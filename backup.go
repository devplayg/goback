package goback

import (
	"fmt"
	"github.com/devplayg/goutils"
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
	Id                int
	srcDirMap         map[string]*dirInfo
	summary           *Summary
	summaries         []*Summary
	workerCount       int
	version           int
	started           time.Time
	rank              int
	minFileSizeToRank int64
	keeper            Keeper
	dbDir             string
	job               *Job
}

func NewBackup(id int, job *Job, dbDir string, keeper Keeper, started time.Time) *Backup {
	return &Backup{
		Id:                id,
		job:               job,
		srcDirMap:         make(map[string]*dirInfo),
		workerCount:       runtime.GOMAXPROCS(0) * 2,
		version:           3,
		started:           started,
		rank:              50,
		minFileSizeToRank: 10 * MB,
		keeper:            keeper,
		dbDir:             dbDir,
		summaries:         make([]*Summary, 0),
	}
}

// Initialize backup
func (b *Backup) init() error {
	if err := b.initKeeper(); err != nil {
		return err
	}

	return nil
}

func (b *Backup) initKeeper() error {
	if err := b.keeper.Init(b.started); err != nil {
		return fmt.Errorf("failed to initialize the keeper; %w", err)
	}

	return nil
}

// Start backup
func (b *Backup) Start() ([]*Summary, error) {
	if err := b.init(); err != nil {
		return nil, err
	}

	// Backup directory sequentially
	for _, dir := range b.job.SrcDirs {
		if err := b.startDirBackup(dir); err != nil {
			log.Error(err)
		}
	}

	// Stop backup
	if err := b.Stop(); err != nil {
		return nil, err
	}

	return b.summaries, nil
}

// Sequentially
func (b *Backup) startDirBackup(dir string) error {
	srcDir, err := IsValidDir(dir)
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
	if err := goutils.GobDecode(data, &files); err != nil {
		return nil, err
	}

	fileMap := sync.Map{}
	for _, f := range files {
		fileMap.Store(f.Path, f.WrapInFileWrapper(true))
	}

	return &fileMap, nil
}

func (b *Backup) issueSummary(dir string, backupType int) {
	summary := NewSummary(backupType, dir, b)
	b.summaries = append(b.summaries, summary)
	b.summary = summary

	log.WithFields(logrus.Fields{
		"protocol": b.keeper.Description().Protocol,
		"host":     b.keeper.Description().Host,
		"dir":      dir,
	}).Infof("backup started %s", strings.Repeat("-", 30))
}

func (b *Backup) Stop() error {
	if b.keeper.Active() {
		if err := b.keeper.Close(); err != nil {
			log.Error(err)
		}
	}
	return nil
}

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
