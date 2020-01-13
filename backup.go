package goback

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Backup struct {
	srcDirArr        []string
	dstDir           string
	db               *bolt.DB
	fileDb           *bolt.DB
	summary          *Summary
	hashComparision  bool
	debug            bool
	workerCount      int
	fileBackupEnable bool
	tempDir          string

	addedFiles    *sync.Map
	modifiedFiles *sync.Map
	deletedFiles  *sync.Map
	failedFiles   *sync.Map

	addedData    []byte
	modifiedData []byte
	deletedData  []byte
	failedData   []byte

	lastFileMap *sync.Map
}

func NewBackup(srcDirArr []string, dstDir string, hashComparision, debug bool) *Backup {
	b := Backup{
		srcDirArr:       srcDirArr,
		dstDir:          dstDir,
		hashComparision: hashComparision,
		debug:           debug,
		workerCount:     runtime.NumCPU(),

		addedFiles:       &sync.Map{},
		modifiedFiles:    &sync.Map{},
		deletedFiles:     &sync.Map{},
		failedFiles:      &sync.Map{},
		fileBackupEnable: true,
	}
	return &b
}

func (b *Backup) init() error {
	if err := b.initDirectories(); err != nil {
		return err
	}

	if err := b.initDatabase(); err != nil {
		return err
	}

	return nil
}

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
		}).Infof("directory to backup")
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
	}).Infof("storage")

	return nil
}

// Initialize database
func (b *Backup) initDatabase() error {
	db, fileDb, err := InitDatabase(b.dstDir)
	if err != nil {
		return err
	}
	b.db, b.fileDb = db, fileDb
	return nil
	//db, err := bolt.Open(filepath.Join(b.dstDir, "backup_log.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	//if err != nil {
	//	return err
	//}
	//err = db.Update(func(tx *bolt.Tx) error {
	//	_, err = tx.CreateBucketIfNotExists(BucketSummary)
	//	return err
	//})
	//if err != nil {
	//	return err
	//}
	//b.db = db
	//
	//fileDb, err := bolt.Open(filepath.Join(b.dstDir, "backup_origin.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	//if err != nil {
	//	return err
	//}
	//err = fileDb.Update(func(tx *bolt.Tx) error {
	//	_, err = tx.CreateBucketIfNotExists(BucketFiles)
	//	return err
	//})
	//if err != nil {
	//	return err
	//}
	//b.fileDb = fileDb

	return nil
}

func (b *Backup) getLastFileMap() (*sync.Map, int64, error) {
	fileMap := sync.Map{}
	var count int64
	err := b.fileDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketFiles)
		if b == nil {
			return ErrorBucketNotFound
		}
		return b.ForEach(func(k, v []byte) error {
			var file File
			if err := json.Unmarshal(v, &file); err != nil {
				return err
			}
			fileMap.Store(string(k), &file)
			count++
			return nil
		})
	})
	return &fileMap, count, err
}

func (b *Backup) Start() error {
	if err := b.init(); err != nil {
		return err
	}
	defer b.Stop()

	t := time.Now()
	lastSummary, err := b.getLastSummary()
	if err != nil {
		return err
	}

	lastFileMap, lastBackupFileCount, err := b.getLastFileMap()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"duration": time.Since(t).Seconds(),
	}).Debug("load last backup")

	if lastSummary == nil || lastSummary.TotalCount < 1 || !IsEqualStringSlices(lastSummary.SrcDirArr, b.srcDirArr) || lastBackupFileCount == 0 {
		return b.generateFirstBackupData()
	}
	b.lastFileMap = lastFileMap

	if err := b.startBackup(); err != nil {
		return err
	}

	return nil

	// Load last backup data
	//lastSummary, err := b.getLastSummary()
	//if err != nil {
	//	return err
	//}
	//spew.Dump(lastSummary)

	//originMap, originCount := b.getOriginMap(lastSummary)
	//spew.Dump(originMap)
	//spew.Dump(originCount)

	// Write initial data to database
	//newMap := sync.Map{}
	//if originCount < 1 || b.srcDir != lastSummary.SrcDir {
	//	b.summary.State = Running
	//	b.summary.Message = "collecting initialize data"
	//	log.Info(b.summary.Message)
	//
	//	err := filepath.Walk(b.srcDir, func(path string, f os.FileInfo, err error) error {
	//		if f.IsDir() {
	//			return nil
	//		}
	//
	//		if !f.Mode().IsRegular() {
	//			return nil
	//		}
	//
	//		fi := newFile(path, f.Size(), f.ModTime())
	//		newMap.Store(path, fi)
	//		b.summary.TotalCount += 1
	//		b.summary.TotalSize += uint64(f.Size())
	//		return nil
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	os.RemoveAll(b.tempDir)
	//	b.summary.ReadingTime = time.Now()
	//	b.summary.ComparisonTime = b.summary.ReadingTime
	//
	//	log.Infof("writing initial data")
	//	b.writeToDatabase(newMap, sync.Map{})
	//	b.summary.LoggingTime = time.Now()
	//	return nil
	//}
	//	b.summary.ReadingTime = time.Now()
	//
	//	// Search files and compare with previous data
	//	log.Infof("comparing old and new")
	//	b.summary.State = 3
	//	i := 1
	//	err := filepath.Walk(b.srcDir, func(path string, f os.FileInfo, err error) error {
	//		if !f.IsDir() && f.Mode().IsRegular() {
	//
	//			log.Debugf("Start checking: [%d] %s (%d)", i, path, f.Size())
	//			atomic.AddUint32(&b.summary.TotalCount, 1)
	//			atomic.AddUint64(&b.summary.TotalSize, uint64(f.Size()))
	//			fi := newFile(path, f.Size(), f.ModTime())
	//
	//			if inf, ok := originMap.Load(path); ok {
	//				last := inf.(*File)
	//
	//				if last.ModTime.Unix() != f.ModTime().Unix() || last.Size != f.Size() {
	//					log.Debugf("modified: %s", path)
	//					fi.State = FileModified
	//					atomic.AddUint32(&b.summary.BackupModified, 1)
	//					backupPath, dur, err := b.BackupFile(path)
	//					if err != nil {
	//						atomic.AddUint32(&b.summary.BackupFailure, 1)
	//						log.Error(err)
	//						fi.Message = err.Error()
	//						fi.State = fi.State * -1
	//						//spew.Dump(fi)
	//					} else {
	//						fi.Message = fmt.Sprintf("copy_time=%4.1f", dur)
	//						atomic.AddUint32(&b.summary.BackupSuccess, 1)
	//						atomic.AddUint64(&b.summary.BackupSize, uint64(f.Size()))
	//						os.Chtimes(backupPath, f.ModTime(), f.ModTime())
	//						originMap.Delete(path)
	//					}
	//				}
	//				originMap.Delete(path)
	//			} else {
	//				log.Debugf("added: %s", path)
	//				fi.State = FileAdded
	//				atomic.AddUint32(&b.summary.BackupAdded, 1)
	//				backupPath, dur, err := b.BackupFile(path)
	//				if err != nil {
	//					atomic.AddUint32(&b.summary.BackupFailure, 1)
	//					log.Error(err)
	//					fi.Message = err.Error()
	//					fi.State = fi.State * -1
	//					//spew.Dump(fi)
	//				} else {
	//					fi.Message = fmt.Sprintf("copy_time=%4.1f", dur)
	//					atomic.AddUint32(&b.summary.BackupSuccess, 1)
	//					atomic.AddUint64(&b.summary.BackupSize, uint64(f.Size()))
	//					os.Chtimes(backupPath, f.ModTime(), f.ModTime())
	//				}
	//			}
	//			//if fi.State < 0 {
	//			//	log.Debugf("[%d] %s", fi.State, fi.Path)
	//			//}
	//			newMap.Store(path, fi)
	//			i++
	//		}
	//		return nil
	//	})
	//	checkErr(err)
	//
	//	// Rename directory
	//	lastDir := filepath.Join(b.dstDir, b.summary.Date.Format("20060102"))
	//	err = os.Rename(b.tempDir, lastDir)
	//	if err == nil {
	//		b.summary.DstDir = lastDir
	//	} else {
	//
	//		i := 1
	//		for err != nil && i <= 10 {
	//			altDir := lastDir + "_" + strconv.Itoa(i)
	//			err = os.Rename(b.tempDir, altDir)
	//			if err == nil {
	//				b.summary.DstDir = altDir
	//			}
	//			i += 1
	//		}
	//		if err != nil {
	//			b.summary.Message = err.Error()
	//			b.summary.State = -1
	//			b.summary.DstDir = b.tempDir
	//			os.RemoveAll(b.tempDir)
	//			return err
	//		}
	//	}
	//	b.summary.ComparisonTime = time.Now()
	//
	//	// Write data to database
	//	err = b.writeToDatabase(newMap, originMap)
	//	b.summary.LoggingTime = time.Now()
	//	return err
	return nil
}

func (b *Backup) Stop() error {
	if err := b.db.Close(); err != nil {
		log.Error(err)
	}

	if err := b.fileDb.Close(); err != nil {
		log.Error(err)
	}

	return nil
}

func (b *Backup) writeToDatabase(newMap sync.Map, originMap sync.Map) error {
	//log.Info("writing to database")
	//
	//rs, err := b.dbLogTx.Exec("insert into bak_summary(date,src_dir,dst_dir,state,total_size,total_count,backup_modified,backup_added,backup_deleted,backup_success,backup_failure,backup_size,execution_time,message) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
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
	//)
	//if err != nil {
	//	return err
	//}
	//
	//b.summary.Id, _ = rs.LastInsertId()
	//log.Infof("backup_id=%d", b.summary.ID)
	//
	//var maxInsertSize uint32 = 500
	//var lines []string
	//var eventLines []string
	//var i uint32 = 0
	//var j uint32 = 0
	//
	//// Delete original data
	//b.dbOriginTx.Exec("delete from bak_origin")
	//
	//// Modified or added files
	//newMap.Range(func(key, value interface{}) bool {
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
	//})
	//if len(eventLines) > 0 {
	//	err := b.insertIntoLog(eventLines)
	//	checkErr(err)
	//	eventLines = nil
	//}
	//
	//// Deleted files
	//eventLines = make([]string, 0)
	//j = 0
	//originMap.Range(func(key, value interface{}) bool {
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
	//})
	//if len(eventLines) > 0 {
	//	err := b.insertIntoLog(eventLines)
	//	checkErr(err)
	//	eventLines = nil
	//}
	//atomic.AddUint32(&b.summary.BackupDeleted, j)

	return nil
}

//func (b *Backup) insertIntoLog(rows []string) error {
//	query := fmt.Sprintf("insert into bak_log(id, path, size, mtime, state, message) %s", strings.Join(rows, " union all "))
//	_, err := b.dbLogTx.Exec(query)
//	return err
//}
//
//func (b *Backup) insertIntoOrigin(rows []string) error {
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
//}
//
//func (b *Backup) Close() error {
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
//}
//
//func (b *Backup) BackupFile(path string) (string, float64, error) {
// Set source
//t := time.Now()
//from, err := os.Open(path)
//if err != nil {
//	return "", time.Since(t).Seconds(), err
//
//}
//defer from.Close()

// Set destination
///data/a/a.txt
///backup
///backup/temp/
//dst := filepath.Join(b.tempDir, path)
//log.Debug(dst)
//return "", 0.0, nil
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
//}
//
//func checkErr(err error) {
//	if err != nil {
//		log.Errorf("[Error] %s", err.Error())
//	}
//}

func (b *Backup) newSummary() (*Summary, error) {
	id, err := IssueDbInt64Id(b.db, BucketSummary)
	if err != nil {
		return nil, err
	}

	return &Summary{
		Id:          id,
		Date:        time.Now(),
		SrcDirArr:   b.srcDirArr,
		DstDir:      b.dstDir,
		WorkerCount: b.workerCount,
		//State:     BackupReady,
		Version: 1,
	}, nil
}

func (b *Backup) getLastSummary() (*Summary, error) {
	_, val, err := GetLastDbData(b.db, BucketSummary)
	if err != nil {
		return nil, err
	}
	//log.WithFields(log.Fields{
	//	"key":   key,
	//	"value": string(val),
	//}).Debug("last backup")

	if len(val) == 0 {
		return nil, nil
	}

	return UnmarshalSummary(val)
}
