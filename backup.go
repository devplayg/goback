package goback

import (
	"fmt"
	"github.com/devplayg/golibs/converter"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Backup struct {
	Id    int
	debug bool
	//srcDirs     []string
	srcDirMap   map[string]*dirInfo
	backupDir   string
	summaryDb   *os.File
	summary     *Summary
	summaries   []*Summary
	workerCount int
	//fileBackupEnable  bool
	version           int
	started           time.Time
	rank              int
	minFileSizeToRank int64
	keeper            Keeper
	//workingDir        string
	dbDir string
	//backupType int
	job *Job
}

func NewBackup(job *Job, dbDir string, keeper Keeper, debug bool) *Backup {
	return &Backup{
		job: job,
		//srcDirs:     job.SrcDirs,
		srcDirMap:   make(map[string]*dirInfo),
		debug:       debug,
		workerCount: runtime.GOMAXPROCS(0) * 2,
		//fileBackupEnable:  true,
		version:           2,
		started:           time.Now(),
		rank:              50,
		minFileSizeToRank: 10 * MB,
		keeper:            keeper,
		backupDir:         "",
		//backupType:        job.BackupType,
		dbDir:     dbDir,
		summaries: make([]*Summary, 0),
	}
}

// Initialize backup
func (b *Backup) init() error {
	//if err := b.initDirectories(); err != nil {
	//	return err
	//}
	//if err := b.initDatabase(); err != nil {
	//	return err
	//}
	//if err := b.loadSummaryDb(); err != nil {
	//	return err
	//}
	if err := b.initKeeper(); err != nil {
		return err
	}

	return nil
}

// Initialize directories
//func (b *Backup) initDirectories() error {
//	if len(b.srcDirs) < 1 {
//		return errors.New("source directories not found")
//	}
//
//	// Working directory
//	workingDir, err := filepath.Abs(os.Args[0])
//	if err != nil {
//		return err
//	}
//	b.workingDir = filepath.Dir(workingDir)
//
//	// Database directory
//	dbDir := filepath.Join(b.workingDir, "db")
//	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
//		if err := os.Mkdir(dbDir, 0600); err != nil {
//			return fmt.Errorf("unable to create database directory: %w", err)
//		}
//	}
//	b.DbDir = dbDir
//
//	return nil
//}

// Initialize database
//func (b *Backup) initDatabase() error {
//	path := filepath.Join(b.DbDir, SummaryDbName)
//	summaryDb, err := LoadOrCreateDatabase(path)
//	if err != nil {
//		return fmt.Errorf("failed to load summary database: %w", err)
//	}
//	b.summaryDb = summaryDb
//	return nil
//}

func (b *Backup) initKeeper() error {
	if err := b.keeper.Init(b.started); err != nil {
		log.WithFields(logrus.Fields{
			"protocol": b.keeper.Description().Protocol,
			"host":     b.keeper.Description().Host,
		}).Errorf("failed to initialize the keeper: %s", err.Error())
		return err
	}

	return nil
}

// Start backup
func (b *Backup) Start() error {
	if err := b.init(); err != nil {
		return err
	}
	// Log.WithFields(Log.Fields{
	// 	"directory": b.dstDir,
	// 	"backupId":  b.Id,
	// }).Infof("backup started")

	defer func() {
		log.WithFields(logrus.Fields{
			"execTime": time.Since(b.started).Seconds(),
			"dirCount": len(b.srcDirMap),
			//"backupId": b.Id,
		}).Info("backup process done")
	}()

	// Backup directory sequentially
	for _, dir := range b.job.SrcDirs {
		if err := b.doBackup(dir); err != nil {
			log.Error(err)
		}
	}

	// Stop backup
	if err := b.Stop(); err != nil {
		return err
	}

	return nil
}

func (b *Backup) doBackup(dirToBackup string) error {
	srcDir, err := IsValidDir(dirToBackup)
	if err != nil {
		return fmt.Errorf("invalid source directory: %w", err)
	}

	if _, ok := b.srcDirMap[srcDir]; ok {
		return nil
	}
	b.srcDirMap[srcDir] = newDirInfo(srcDir, b.dbDir)

	lastFileMap, err := b.loadLastFileMap(srcDir)
	if err != nil {
		return fmt.Errorf("failed to load last backup data for %s: %w", srcDir, err)
	}

	// Full backup
	if b.job.BackupType == Full {
		return nil
	}

	// Incremental backup
	if lastFileMap == nil { // First backup
		return b.generateFirstBackupData(srcDir)
	}

	if err := b.startBackup(srcDir, lastFileMap); err != nil {
		return err
	}

	return nil
}

func (b *Backup) loadLastFileMap(dir string) (*sync.Map, error) {
	if _, err := os.Stat(b.srcDirMap[dir].dbPath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := ioutil.ReadFile(b.srcDirMap[dir].dbPath)
	if err != nil || len(data) < 1 {
		return nil, err
	}

	var files []*File
	if err := converter.DecodeFromBytes(data, &files); err != nil {
		return nil, err
	}

	fileMap := sync.Map{}
	for _, f := range files {
		fileMap.Store(f.Path, f.WrapInFileWrapper(true))
	}

	return &fileMap, nil
}

func (b *Backup) issueSummary(dir string, backupType int) {
	//summaryId := len(b.summaries) + 1
	summary := NewSummary(backupType, dir, b)
	b.summaries = append(b.summaries, summary)
	b.summary = summary

	log.WithFields(logrus.Fields{
		//"summaryId": summaryId,
		//"backupId":  b.Id,
		"dir": dir,
	}).Infof("backup started %s", strings.Repeat("-", 30))
}

func (b *Backup) Stop() error {
	if b.keeper.Active() {
		if err := b.keeper.Close(); err != nil {
			log.Error(err)
		}
	}

	if err := b.writeSummary(); err != nil {
		return err
	}
	if err := b.summaryDb.Close(); err != nil {
		return err
	}

	return nil
}

//
//func (b *Backup) loadSummaryDb() error {
//	data, err := ioutil.ReadAll(b.summaryDb)
//	if err != nil {
//		return err
//	}
//	if len(data) < 1 {
//		b.Id = 1
//		b.summaries = make([]*Summary, 0)
//		return nil
//	}
//
//	decompressed, err := compress.Decompress(data, compress.GZIP)
//	if err != nil {
//		return err
//	}
//
//	var summaries []*Summary
//	if err := converter.DecodeFromBytess(decompressed, &summaries); err != nil {
//		return fmt.Errorf("failed to load summary database: %w", err)
//	}
//	if len(summaries) < 1 {
//		b.Id = 1
//	} else {
//		b.Id = summaries[len(summaries)-1].BackupId + 1
//	}
//	// spew.Dump(summaries)
//
//	b.summaries = summaries
//	return nil
//}

func (b *Backup) writeBackupState(state int) {
	t := time.Now()
	b.summary.State = state

	if state == Read {
		b.summary.ReadingTime = t
		return
	}

	if state == Compared {
		b.summary.ComparisonTime = t
		return

	}
	if state == Copied {
		b.summary.BackupTime = t
		return

	}
	if state == Logged {
		b.summary.LoggingTime = t
		b.summary.ExecutionTime = b.summary.LoggingTime.Sub(b.summary.Date).Seconds()
		return

	}
}

// Using go-routine
func (b *Backup) getCurrentFileMaps(dir string) ([]*sync.Map, error) {
	fileMaps := make([]*sync.Map, b.workerCount)
	for i := range fileMaps {
		fileMaps[i] = &sync.Map{}
	}

	i := 0
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		if !f.Mode().IsRegular() {
			return nil
		}

		fileWrapper := NewFileWrapper(path, f.Size(), f.ModTime())
		// if b.hashComparision {
		// 	h, err := GetFileHash(path)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	fileWrapper.Hash = h
		// }

		// Statistics
		b.summary.TotalSize += uint64(fileWrapper.Size)
		b.summary.TotalCount++
		b.summary.Stats.addToStats(fileWrapper.WrapInFileGrid())

		// Distribute works
		workerId := i % b.workerCount
		fileMaps[workerId].Store(path, fileWrapper)
		i++

		return nil
	})
	if err != nil {
		return nil, err
	}

	b.summary.Stats.rank(b.rank)
	b.writeBackupState(Read)
	log.WithFields(logrus.Fields{
		"execTime": b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
		"files":    b.summary.TotalCount,
		"dir":      dir,
		"size":     fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
	}).Info("read files")

	return fileMaps, nil
}
