package goback

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Backup struct {
	debug bool

	srcDirs   []string
	srcDirMap map[string]string
	dstDir    string

	// summaryDbPath string
	summaryDb *os.File
	// fileMapDbPath string
	fileMapDb *os.File

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

	// lastFileMap   *sync.Map
	version       int
	nextSummaryId int64
	changesLog    map[string]interface{}
}

func NewBackup(srcDirs []string, dstDir string, hashComparision, debug bool) *Backup {
	return &Backup{
		srcDirs:          srcDirs,
		srcDirMap:        make(map[string]string),
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
}

// Initialize backup
func (b *Backup) init() error {
	if err := b.initDirectories(); err != nil {
		return err
	}

	if err := b.initDatabase(); err != nil {
		return err
	}
	if err := b.loadSummaryDb(); err != nil {
		return err
	}

	return nil
}

// Initialize directories
func (b *Backup) initDirectories() error {
	if len(b.srcDirs) < 1 {
		return errors.New("empty source directories")
	}

	// Check if source directories is valid
	// 	for src, h := range b.srcDirMap {
	// 		dir, err := filepath.Abs(src)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		b.srcDirMap[dir] = h
	// 		if err := isValidDir(dir); err != nil {
	// 			return errors.New("source directory error: " + err.Error())
	// 		}
	//
	// 		log.WithFields(log.Fields{
	// 			"dir": dir,
	// 		}).Infof("source")
	// 	}
	//
	// Check if destination is valid
	if len(b.dstDir) < 1 {
		return errors.New("empty source directories")
	}
	dstDir, err := IsValidDir(b.dstDir)
	if err != nil {
		return errors.New("destination directory error: " + err.Error())
	}
	b.dstDir = dstDir

	log.WithFields(log.Fields{
		"directory": b.dstDir,
	}).Infof("backup")

	return nil
}

// Initialize database
func (b *Backup) initDatabase() error {
	summaryDb, err := LoadOrCreateDatabase(filepath.Join(b.dstDir, SummaryDbName))
	if err != nil {
		return fmt.Errorf("failed to load summary database: %w", err)
	}
	b.summaryDb = summaryDb
	// b.summaryDbPath = filepath.Join(b.dstDir, SummaryDbName)
	// 	b.fileMapDbPath = filepath.Join(b.dstDir, FileMapDbName)
	// 	summaryDb, fileMapDb, err := InitDatabase(b.summaryDbPath, b.fileMapDbPath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	b.summaryDb, b.fileMapDb = summaryDb, fileMapDb
	//
	// 	return err
	// }
	//
	// func (b *Backup) loadFileMapDb() (*sync.Map, error) {
	// 	data, err := ioutil.ReadAll(b.fileMapDb)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	if len(data) < 1 {
	// 		return &sync.Map{}, nil
	// 	}
	//
	// 	var files []*File
	// 	if err := converter.DecodeFromBytes(data, &files); err != nil {
	// 		return nil, err
	// 	}
	//
	// 	fileMap := sync.Map{}
	// 	for _, f := range files {
	// 		fileMap.Store(f.Path, f)
	// 	}
	//
	// 	return &fileMap, nil
	return nil
}

func (b *Backup) Start() error {
	if err := b.init(); err != nil {
		return err
	}
	defer b.Stop()

	for _, dir := range b.srcDirs {
		if err := b.doBackup(dir); err != nil {
			log.Error(err)
		}
	}

	// 	// Load last backup
	// 	lastSummary, lastFileMap, err := b.loadLastBackup()
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	// Issue summary
	// 	if err := b.issueSummary(b.nextSummaryId, b.srcDirs, b.dstDir, b.workerCount, b.version); err != nil {
	// 		return err
	// 	}
	//
	// 	// Create temp directory
	// 	tempDir, err := ioutil.TempDir(b.dstDir, "backup-")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	b.tempDir = tempDir
	// 	defer func() {
	// 		dstDir := filepath.Join(b.dstDir, b.summary.Date.Format("20060102")+"-"+strconv.FormatInt(b.summary.Id, 10))
	// 		if err := os.Rename(b.tempDir, dstDir); err != nil {
	// 			log.Error(err)
	// 		}
	// 	}()
	//
	// 	if lastSummary == nil || lastSummary.TotalCount < 1 || !IsEqualStringSlices(lastSummary.SrcDirArr, b.srcDirs) {
	// 		return b.generateFirstBackupData()
	// 	}
	// 	b.lastFileMap = lastFileMap
	//
	// 	if err := b.startBackup(); err != nil {
	// 		return err
	// 	}
	//
	return nil
}

func (b *Backup) doBackup(target string) error {
	dir, err := IsValidDir(target)
	if err != nil {
		return errors.New("source directory error: " + err.Error())
	}

	if _, have := b.srcDirMap[dir]; have {
		return errors.New("duplicate source directory: " + dir)
	}

	dirSum := md5.Sum([]byte(dir))
	b.srcDirMap[dir] = hex.EncodeToString(dirSum[:])
	log.Debugf("%s => %s", dir, b.srcDirMap[dir])

	// lastFileMap, err := b.loadLastFileMap(dir)
	// if err != nil {
	// 	return fmt.Errorf("failed to load last file map: %w", err)
	// }

	// spew.Dump(b.summaries)
	// log.Debug(absDir)
	// Load last backup
	// lastSummary, lastFileMap, err := b.loadLastBackup()
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	// Issue summary
	// 	if err := b.issueSummary(b.nextSummaryId, b.srcDirs, b.dstDir, b.workerCount, b.version); err != nil {
	// 		return err
	// 	}
	//
	// 	// Create temp directory
	// 	tempDir, err := ioutil.TempDir(b.dstDir, "backup-")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	b.tempDir = tempDir
	// 	defer func() {
	// 		dstDir := filepath.Join(b.dstDir, b.summary.Date.Format("20060102")+"-"+strconv.FormatInt(b.summary.Id, 10))
	// 		if err := os.Rename(b.tempDir, dstDir); err != nil {
	// 			log.Error(err)
	// 		}
	// 	}()
	//
	// 	if lastSummary == nil || lastSummary.TotalCount < 1 || !IsEqualStringSlices(lastSummary.SrcDirArr, b.srcDirs) {
	// 		return b.generateFirstBackupData()
	// 	}
	// 	b.lastFileMap = lastFileMap
	//
	// 	if err := b.startBackup(); err != nil {
	// 		return err
	// 	}
	//

	return nil
}

func (b *Backup) loadLastFileMap(dir string) (*sync.Map, error) {
	// 	data, err := ioutil.ReadAll(b.fileMapDb)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	if len(data) < 1 {
	// 		return &sync.Map{}, nil
	// 	}
	//
	// 	var files []*File
	// 	if err := converter.DecodeFromBytes(data, &files); err != nil {
	// 		return nil, err
	// 	}
	//
	// 	fileMap := sync.Map{}
	// 	for _, f := range files {
	// 		fileMap.Store(f.Path, f)
	// 	}
	//
	// 	return &fileMap, nil

	return nil, nil
}

//
// func (b *Backup) issueSummary(id int64, srcDirs []string, dstDir string, workCount, version int) error {
// 	b.summary = NewSummary(id, srcDirs, dstDir, workCount, version)
//
// 	tempSummaries := append(b.summaries, b.summary)
// 	data, err := converter.EncodeToBytes(tempSummaries)
// 	if err != nil {
// 		return err
// 	}
// 	if err := b.summaryDb.Truncate(0); err != nil {
// 		return err
// 	}
// 	if _, err := b.summaryDb.WriteAt(data, 0); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (b *Backup) loadLastBackup() (*Summary, *sync.Map, error) {
// 	t := time.Now()
// 	lastSummary, backupCount, err := b.getLastSummary()
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	b.nextSummaryId = int64(backupCount) + 1
//
// 	fileMap, err := b.loadFileMapDb()
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	if lastSummary != nil {
// 		log.WithFields(log.Fields{
// 			"loadingTime": time.Since(t).Seconds(),
// 			"files":       lastSummary.TotalCount,
// 			"summaryId":   lastSummary.Id,
// 			"size":        fmt.Sprintf("%d(%s)", lastSummary.TotalSize, humanize.Bytes(lastSummary.TotalSize)),
// 			"date":        lastSummary.Date.Format(DefaultDateFormat),
// 		}).Info("last backup")
// 	}
// 	return lastSummary, fileMap, err
// }
//
// func (b *Backup) getLastSummary() (*Summary, int, error) {
// 	summaries, err := b.loadSummaryDb()
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	if summaries == nil {
// 		return nil, 0, nil
// 	}
// 	b.summaries = summaries
//
// 	return summaries[len(summaries)-1], len(summaries), nil
// }
//
func (b *Backup) Stop() error {
	if err := b.summaryDb.Close(); err != nil {
		log.Error(err)
	}

	// 	if err := b.fileMapDb.Close(); err != nil {
	// 		log.Error(err)
	// 	}
	//
	return nil
}

func (b *Backup) loadSummaryDb() error {
	data, err := ioutil.ReadAll(b.summaryDb)
	if err != nil {
		return err
	}

	if len(data) < 1 {
		b.summaries = make([]*Summary, 0)
		return nil
	}

	var summaries []*Summary
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&summaries); err != nil {
		return fmt.Errorf("failed to load summary database: %w", err)
	}
	b.summaries = summaries
	return nil
}
