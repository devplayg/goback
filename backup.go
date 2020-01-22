package goback

import (
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
    "strings"
    "sync"
    "time"
)

type Backup struct {
    Id    int
    debug bool

    srcDirs   []string
    srcDirMap map[string]*DirInfo
    dstDir    string

    summaryDb        *os.File
    summary          *Summary
    summaries        []*Summary
    hashComparision  bool
    workerCount      int
    fileBackupEnable bool
    tempDir          string

    version int
    started time.Time
    rank    int
}

func NewBackup(srcDirs []string, dstDir string, hashComparision, debug bool) *Backup {
    return &Backup{
        srcDirs:          srcDirs,
        srcDirMap:        make(map[string]*DirInfo),
        dstDir:           dstDir,
        hashComparision:  hashComparision,
        debug:            debug,
        workerCount:      runtime.NumCPU(),
        fileBackupEnable: true,
        version:          1,
        started:          time.Now(),
        rank:             50,
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
        return errors.New("source directories not found")
    }

    // Check if destination is valid
    if len(b.dstDir) < 1 {
        return errors.New("backup directory not found")
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
    return nil
}

// Start backup
func (b *Backup) Start() error {
    if err := b.init(); err != nil {
        return err
    }
    defer func() {
        log.WithFields(log.Fields{
            "execTime": time.Since(b.started).Seconds(),
            "dirCount": len(b.srcDirMap),
        }).Info("backup process done")
    }()

    // Create temp directory
    tempDir, err := ioutil.TempDir(b.dstDir, "backup-")
    if err != nil {
        return err
    }
    b.tempDir = tempDir

    // Backup directory sequentially
    for _, dir := range b.srcDirs {
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
        return errors.New("source directory error: " + err.Error())
    }

    if _, have := b.srcDirMap[srcDir]; have {
        return errors.New("duplicate source directory: " + srcDir)
    }

    b.srcDirMap[srcDir] = NewDirInfo(srcDir, b.dstDir)
    lastFileMap, err := b.loadLastFileMap(srcDir)
    if err != nil {
        return fmt.Errorf("failed to load last backup data for %s: %w", srcDir, err)
    }

    log.Debugf("%s backup started", strings.Repeat("=", 50))
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
        fileMap.Store(f.Path, WrapFile(f))
    }

    return &fileMap, nil
}

func (b *Backup) issueSummary(dir string, backupType int) {
    summaryId := len(b.summaries) + 1
    summary := NewSummary(summaryId, b.Id, dir, b.dstDir, backupType, b.workerCount, b.version)
    b.summaries = append(b.summaries, summary)
    b.summary = summary

    log.WithFields(log.Fields{
        "id":  summaryId,
        "dir": dir,
    }).Debug("issued new summary")
}

func (b *Backup) Stop() error {
    if err := b.summaryDb.Close(); err != nil {
        log.Error(err)
    }

    dstDir := filepath.Join(b.dstDir, b.summary.Date.Format("20060102")+"-"+strconv.Itoa(b.summary.BackupId))
    if err := os.Rename(b.tempDir, dstDir); err != nil {
        log.Error(err)
    }

    return nil
}

func (b *Backup) loadSummaryDb() error {
    data, err := ioutil.ReadAll(b.summaryDb)
    if err != nil {
        return err
    }

    if len(data) < 1 {
        b.Id = 1
        b.summaries = make([]*Summary, 0)
        return nil
    }

    var summaries []*Summary
    if err := converter.DecodeFromBytes(data, &summaries); err != nil {
        return fmt.Errorf("failed to load summary database: %w", err)
    }
    if len(summaries) < 1 {
        b.Id = 1
    } else {
        b.Id = summaries[len(summaries)-1].BackupId + 1
    }

    b.summaries = summaries
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
        return

    }
    if state == Completed {
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
    err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
        if file.IsDir() {
            return nil
        }

        if !file.Mode().IsRegular() {
            return nil
        }

        fileWrapper := NewFileWrapper(path, file.Size(), file.ModTime())
        if b.hashComparision {
            h, err := GetFileHash(path)
            if err != nil {
                return err
            }
            fileWrapper.Hash = h
        }

        // Statistics
        b.summary.TotalSize += uint64(fileWrapper.Size)
        b.summary.TotalCount++
        b.summary.Report.addToReport(fileWrapper.File)

        // Distribute works
        workerId := i % b.workerCount
        fileMaps[workerId].Store(path, fileWrapper)
        i++

        return nil
    })
    if err != nil {
        return nil, err
    }

    b.summary.Report.tune(b.rank)
    b.writeBackupState(Read)
    log.WithFields(log.Fields{
        "execTime": b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
        "files":    b.summary.TotalCount,
        "dir":      dir,
        "size":     fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
    }).Info("read files")

    return fileMaps, nil
}
