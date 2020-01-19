package goback

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/devplayg/golibs/converter"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Backup struct {
	debug bool

	srcDirArr []string
	dstDir    string

	summaryDbPath string
	summaryDb     *os.File
	fileMapDbPath string
	fileMapDb     *os.File

	summary          *Summary
	summaries        []*Summary
	hashComparision  bool
	workerCount      int
	fileBackupEnable bool
	tempDir          string

	addedFiles    *sync.Map
	modifiedFiles *sync.Map
	deletedFiles  *sync.Map
	failedFiles   *sync.Map

	lastFileMap   *sync.Map
	version       int
	nextSummaryId int64
	changesLog    map[string]interface{}
}

func NewBackup(srcDirArr []string, dstDir string, hashComparision, debug bool) *Backup {
	b := Backup{
		srcDirArr:        srcDirArr,
		dstDir:           dstDir,
		hashComparision:  hashComparision,
		debug:            debug,
		workerCount:      runtime.NumCPU(),
		addedFiles:       &sync.Map{},
		modifiedFiles:    &sync.Map{},
		deletedFiles:     &sync.Map{},
		failedFiles:      &sync.Map{},
		fileBackupEnable: true,
		version:          1,
	}
	return &b
}

// Initialize backup
func (b *Backup) init() error {
	if err := b.initDirectories(); err != nil {
		return err
	}

	if err := b.initDatabase(); err != nil {
		return err
	}

	return nil
}

// Initialize directories
func (b *Backup) initDirectories() error {
	if len(b.srcDirArr) < 1 {
		return errors.New("empty source directories")
	}
	for i := range b.srcDirArr {
		dir, err := filepath.Abs(b.srcDirArr[i])
		if err != nil {
			return err
		}
		b.srcDirArr[i] = dir
		if err := isValidDir(b.srcDirArr[i]); err != nil {
			return errors.New("source directory error: " + err.Error())
		}

		log.WithFields(log.Fields{
			"dir": b.srcDirArr[i],
		}).Infof("source")
	}

	if len(b.dstDir) < 1 {
		return errors.New("empty source directories")
	}
	b.dstDir = filepath.Clean(b.dstDir)
	if err := isValidDir(b.dstDir); err != nil {
		return errors.New("destination directory error: " + err.Error())
	}

	log.WithFields(log.Fields{
		"dir": b.dstDir,
	}).Infof("backup")

	return nil
}

// Initialize database
func (b *Backup) initDatabase() error {
	b.summaryDbPath = filepath.Join(b.dstDir, SummaryDbName)
	b.fileMapDbPath = filepath.Join(b.dstDir, FileMapDbName)
	summaryDb, fileMapDb, err := InitDatabase(b.summaryDbPath, b.fileMapDbPath)
	if err != nil {
		return err
	}
	b.summaryDb, b.fileMapDb = summaryDb, fileMapDb

	return err
}

func (b *Backup) loadFileMapDb() (*sync.Map, error) {
	data, err := ioutil.ReadAll(b.fileMapDb)
	if err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return &sync.Map{}, nil
	}

	var files []*File
	if err := converter.DecodeFromBytes(data, &files); err != nil {
		return nil, err
	}

	fileMap := sync.Map{}
	for _, f := range files {
		fileMap.Store(f.Path, f)
	}

	return &fileMap, nil
}

func (b *Backup) Start() error {
	if err := b.init(); err != nil {
		return err
	}
	defer b.Stop()

	// Load last backup
	lastSummary, lastFileMap, err := b.loadLastBackup()
	if err != nil {
		return err
	}

	// Issue summary
	if err := b.issueSummary(b.nextSummaryId, b.srcDirArr, b.dstDir, b.workerCount, b.version); err != nil {
		return err
	}

	// Create temp directory
	tempDir, err := ioutil.TempDir(b.dstDir, "backup-")
	if err != nil {
		return err
	}
	b.tempDir = tempDir
	defer func() {
		dstDir := filepath.Join(b.dstDir, b.summary.Date.Format("20060102")+"-"+strconv.FormatInt(b.summary.Id, 10))
		if err := os.Rename(b.tempDir, dstDir); err != nil {
			log.Error(err)
		}
	}()

	if lastSummary == nil || lastSummary.TotalCount < 1 || !IsEqualStringSlices(lastSummary.SrcDirArr, b.srcDirArr) {
		return b.generateFirstBackupData()
	}
	b.lastFileMap = lastFileMap

	if err := b.startBackup(); err != nil {
		return err
	}

	return nil
}

func (b *Backup) issueSummary(id int64, srcDirs []string, dstDir string, workCount, version int) error {
	b.summary = NewSummary(b.nextSummaryId, b.srcDirArr, b.dstDir, b.workerCount, b.version)

	tempSummaries := append(b.summaries, b.summary)
	data, err := converter.EncodeToBytes(tempSummaries)
	if err != nil {
		return err
	}
	if err := b.summaryDb.Truncate(0); err != nil {
		return err
	}
	if _, err := b.summaryDb.WriteAt(data, 0); err != nil {
		return err
	}
	return nil
}

func (b *Backup) loadLastBackup() (*Summary, *sync.Map, error) {
	t := time.Now()
	lastSummary, backupCount, err := b.getLastSummary()
	if err != nil {
		return nil, nil, err
	}
	b.nextSummaryId = int64(backupCount) + 1

	fileMap, err := b.loadFileMapDb()
	if err != nil {
		return nil, nil, err
	}

	if lastSummary != nil {
		log.WithFields(log.Fields{
			"loadingTime": time.Since(t).Seconds(),
			"files":       lastSummary.TotalCount,
			"summaryId":   lastSummary.Id,
			"size":        fmt.Sprintf("%d(%s)", lastSummary.TotalSize, humanize.Bytes(lastSummary.TotalSize)),
			"date":        lastSummary.Date.Format(DefaultDateFormat),
		}).Info("last backup")
	}
	return lastSummary, fileMap, err
}

func (b *Backup) getLastSummary() (*Summary, int, error) {
	summaries, err := b.loadSummaryDb()
	if err != nil {
		return nil, 0, err
	}
	if summaries == nil {
		return nil, 0, nil
	}
	b.summaries = summaries

	return summaries[len(summaries)-1], len(summaries), nil
}

func (b *Backup) Stop() error {
	if err := b.summaryDb.Close(); err != nil {
		log.Error(err)
	}

	if err := b.fileMapDb.Close(); err != nil {
		log.Error(err)
	}

	return nil
}

// func (b *Backup) writeToDatabase(newMap sync.Map, originMap sync.Map) error {
// log.Info("writing to database")
//
// rs, err := b.dbLogTx.Exec("insert into bak_summary(date,src_dir,dst_dir,state,total_size,total_count,backup_modified,backup_added,backup_deleted,backup_success,backup_failure,backup_size,execution_time,message) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
//	b.summary.Date.Format(time.RFC3339),
//	b.summary.SrcDir,
//	b.summary.DstDir,
//	b.summary.State,
//	b.summary.TotalSize,
//	b.summary.TotalCount,
//	b.summary.BackupModified,
//	b.summary.BackupAdded,
//	b.summary.BackupDeleted,
//	b.summary.BackupSuccess,
//	b.summary.BackupFailure,
//	b.summary.BackupSize,
//	b.summary.ExecutionTime,
//	b.summary.Message,
// )
// if err != nil {
//	return err
// }
//
// b.summary.Id, _ = rs.LastInsertId()
// log.Infof("backup_id=%d", b.summary.ID)
//
// var maxInsertSize uint32 = 500
// var lines []string
// var eventLines []string
// var i uint32 = 0
// var j uint32 = 0
//
// // Delete original data
// b.dbOriginTx.Exec("delete from bak_origin")
//
// // Modified or added files
// newMap.Range(func(key, value interface{}) bool {
//	f := value.(*File)
//	path := strings.Replace(f.Path, "'", "''", -1)
//	lines = append(lines, fmt.Sprintf("select '%s', %d, '%s'", path, f.Size, f.ModTime.Format(time.RFC3339)))
//
//	i += 1
//
//	if i%maxInsertSize == 0 || i == b.summary.TotalCount {
//		err := b.insertIntoOrigin(lines)
//		checkErr(err)
//		lines = nil
//	}
//
//	if f.State != 0 {
//		eventLines = append(eventLines, fmt.Sprintf("select %d, '%s', %d, '%s', %d, '%s'", b.summary.ID, path, f.Size, f.ModTime.Format(time.RFC3339), f.State, f.Message))
//		j += 1
//
//		if j%maxInsertSize == 0 {
//			err := b.insertIntoLog(eventLines)
//			checkErr(err)
//			eventLines = nil
//		}
//	}
//	return true
// })
// if len(eventLines) > 0 {
//	err := b.insertIntoLog(eventLines)
//	checkErr(err)
//	eventLines = nil
// }
//
// // Deleted files
// eventLines = make([]string, 0)
// j = 0
// originMap.Range(func(key, value interface{}) bool {
//	atomic.AddUint32(&b.summary.BackupSuccess, 1)
//	f := value.(*File)
//	log.Debugf("deleted: %s", f.Path)
//	f.State = FileDeleted
//	path := strings.Replace(f.Path, "'", "''", -1)
//	eventLines = append(eventLines, fmt.Sprintf("select %d, '%s', %d, '%s', %d, '%s'", b.summary.ID, path, f.Size, f.ModTime.Format(time.RFC3339), f.State, f.Message))
//	j += 1
//
//	if j%maxInsertSize == 0 {
//		err := b.insertIntoLog(eventLines)
//		checkErr(err)
//		eventLines = nil
//	}
//	return true
// })
// if len(eventLines) > 0 {
//	err := b.insertIntoLog(eventLines)
//	checkErr(err)
//	eventLines = nil
// }
// atomic.AddUint32(&b.summary.BackupDeleted, j)

// return nil
// }

// func (b *Backup) insertIntoLog(rows []string) error {
//	query := fmt.Sprintf("insert into bak_log(id, path, size, mtime, state, message) %s", strings.Join(rows, " union all "))
//	_, err := b.dbLogTx.Exec(query)
//	return err
// }
//
// func (b *Backup) insertIntoOrigin(rows []string) error {
//	query := fmt.Sprintf("insert into bak_origin(path, size, mtime) %s", strings.Join(rows, " union all "))
//	_, err := b.dbOriginTx.Exec(query)
//	defer func() {
//		if r := recover(); r != nil {
//			if err != nil {
//				log.Println(query)
//			}
//		}
//	}()
//	checkErr(err)
//	return err
// }
//
// func (b *Backup) Close() error {
//	b.summary.ExecutionTime = b.summary.LoggingTime.Sub(b.summary.Date).Seconds()
//	b.summary.Message += fmt.Sprintf("reading: %3.1fs, comparing: %3.1fs, writing: %3.1fs",
//		b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
//		b.summary.ComparisonTime.Sub(b.summary.ReadingTime).Seconds(),
//		b.summary.LoggingTime.Sub(b.summary.ComparisonTime).Seconds(),
//	)
//	b.dbLogTx.Exec("update bak_summary set backup_deleted = ?, execution_time = ?, message = ? where id = ?",
//		b.summary.BackupDeleted,
//		b.summary.ExecutionTime,
//		b.summary.Message,
//		b.summary.ID,
//	)
//
//	b.dbLogTx.Commit()
//	b.dbOriginTx.Commit()
//	b.dbOrigin.Close()
//	b.dbLog.Close()
//
//	if b.summary.ID > 1 { // ID 1 is about initializing data
//		log.WithFields(log.Fields{
//			"modified": b.summary.BackupModified,
//			"added":    b.summary.BackupAdded,
//			"deleted":  b.summary.BackupDeleted,
//		}).Infof("files: %d", b.summary.BackupModified+b.summary.BackupAdded+b.summary.BackupDeleted)
//		log.WithFields(log.Fields{
//			"success": b.summary.BackupSuccess,
//			"failure": b.summary.BackupFailure,
//		}).Infof("backup result")
//		log.Infof("backup size: %d(%s)", b.summary.BackupSize, humanize.Bytes(b.summary.BackupSize))
//	}
//	log.WithFields(log.Fields{
//		"files": b.summary.TotalCount,
//		"size":  fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
//	}).Info("source directory")
//
//	log.WithFields(log.Fields{
//		"reading":    fmt.Sprintf("%3.1fs", b.summary.ReadingTime.Sub(b.summary.Date).Seconds()),
//		"comparison": fmt.Sprintf("%3.1fs", b.summary.ComparisonTime.Sub(b.summary.ReadingTime).Seconds()),
//		"writing":    fmt.Sprintf("%3.1fs", b.summary.LoggingTime.Sub(b.summary.ComparisonTime).Seconds()),
//	}).Infof("execution time: %3.1fs", b.summary.ExecutionTime)
//
//	return nil
// }
//
// func (b *Backup) BackupFile(path string) (string, float64, error) {
// Set source
// t := time.Now()
// from, err := os.Open(path)
// if err != nil {
//	return "", time.Since(t).Seconds(), err
//
// }
// defer from.Close()

// Set destination
// /data/a/a.txt
// /backup
// /backup/temp/
// dst := filepath.Join(b.tempDir, path)
// log.Debug(dst)
// return "", 0.0, nil
//	err = os.MkdirAll(filepath.Dir(dst), 0644)
//	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
//	if err != nil {
//		return "", time.Since(t).Seconds(), err
//	}
//	defer to.Close()
//
//	// Copy
//	_, err = io.Copy(to, from)
//	if err != nil {
//		return "", time.Since(t).Seconds(), err
//	}
//
//	return dst, time.Since(t).Seconds(), err
// }
//
// func checkErr(err error) {
//	if err != nil {
//		log.Errorf("[Error] %s", err.Error())
//	}
// }
//
// func (b *Backup) newSummary() *Summary {
// 	return &Summary{
// 		Date:        time.Now(),
// 		SrcDirArr:   b.srcDirArr,
// 		DstDir:      b.dstDir,
// 		WorkerCount: b.workerCount,
// 		Version: b.version,
// 	}

// 	// id, err := IssueDbInt64Id(b.db, BucketSummary)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	id, err := b.issueSummaryId()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &Summary{
// 		Id:          id,
// 		Date:        time.Now(),
// 		SrcDirArr:   b.srcDirArr,
// 		DstDir:      b.dstDir,
// 		WorkerCount: b.workerCount,
// 		// State:     BackupReady,
// 		Version: 1,
// 	}, nil
// }
//
// func (b *Backup) issueSummaryId() (int64, error) {
// 	summaries, err := b.readSummary()
// 	if summaries == nil {
// 		b.newSummary()
// 	}
// 	// log.Error(len(summaries))
// 	return 0, err
// }

func (b *Backup) loadSummaryDb() ([]*Summary, error) {
	data, err := ioutil.ReadAll(b.summaryDb)
	if err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return nil, nil
	}

	var summaries []*Summary
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&summaries); err != nil {
		return nil, err
	}
	return summaries, nil
}

func Decode(data []byte, to interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(to)
}
