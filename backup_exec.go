package goback

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "os"
    "sync"
    "sync/atomic"
    "time"
)

func (b *Backup) startBackup(srcDir string, lastFileMap *sync.Map) error {
    // 1. Issue summary
    b.issueSummary(srcDir, NormalBackup)

    // 2. Collect files in source directories
    currentFileMaps, err := b.getCurrentFileMaps(srcDir)
    if err != nil {
        return nil
    }

    // 3. Compares file maps
    if err := b.compareFileMaps(currentFileMaps, lastFileMap); err != nil {
        return err
    }

    // 4. Backup added or changed files
    if err := b.backupFiles(); err != nil {
        return err
    }

    // 4. Write result
    if err := b.writeResult(currentFileMaps, lastFileMap); err != nil {
        return err
    }

    // 5. Write summary
    if err := b.writeSummary(); err != nil {
        return err
    }

    return nil
}


func (b *Backup) backupFiles() error {
    fileGroup, count, err := b.createBackupFileGroup()
    if err != nil {
        return err
    }

    log.WithFields(log.Fields{
        "workers": b.workerCount,
    }).Info("running backup..")
    wg := sync.WaitGroup{}
    for i := range fileGroup {
        if len(fileGroup[i]) < 1 {
            continue
        }

        wg.Add(1)
        go func(workerId int) {
            defer wg.Done()
            t := time.Now()
            if err := b.backupFileGroup(workerId, fileGroup[workerId]); err != nil {
                log.Error(err)
            }
            log.WithFields(log.Fields{
                "workerId":  workerId,
                "processed": len(fileGroup[workerId]),
                "duration":  time.Since(t).Seconds(),
            }).Debug("backup done")
        }(i)
    }
    wg.Wait()

    // b.failedFiles.Range(func(key, value interface{}) bool {
    // 	file := key.(*FileWrapper)
    // 	return true
    // })

    b.writeBackupState(Copied)
    if count > 0 {
        log.WithFields(log.Fields{
            "execTime": b.summary.BackupTime.Sub(b.summary.ComparisonTime).Seconds(),
            "success":  b.summary.SuccessCount,
            "failed":   b.summary.FailedCount,
        }).Info("backup report")
    }

    return nil
}

// Added or modified files will be backed up
func (b *Backup) createBackupFileGroup() ([][]*FileWrapper, uint64, error) {
    fileGroup := make([][]*FileWrapper, b.workerCount)
    for i := range fileGroup {
        fileGroup[i] = make([]*FileWrapper, 0)
    }

    var i uint64 = 0
    b.summary.addedFiles.Range(func(k, v interface{}) bool {
        file := k.(*FileWrapper)
        workerId := i % uint64(b.workerCount)
        fileGroup[workerId] = append(fileGroup[workerId], file)
        i++

        return true
    })
    b.summary.modifiedFiles.Range(func(k, v interface{}) bool {
        file := k.(*FileWrapper)
        workerId := i % uint64(b.workerCount)
        fileGroup[workerId] = append(fileGroup[workerId], file)
        i++

        return true
    })

    return fileGroup, i, nil
}

// Thread-safe
func (b *Backup) backupFileGroup(workerId int, files []*FileWrapper) error {
    for _, file := range files {
        path, dur, err := BackupFile(file.Path, b.tempDir)
        if err != nil { // failed to backup
            b.summary.failedFiles.Store(file, nil)
            atomic.AddUint64(&b.summary.FailedCount, 1)
            atomic.AddUint64(&b.summary.FailedSize, uint64(file.Size))
            file.Result = FileBackupFailed
            file.Message = err.Error()
            log.WithFields(log.Fields{
                "workerId": workerId,
            }).Error(fmt.Errorf("failed to backup: %s; %w", file.Path, err))
            continue
        }

        // Success
        atomic.AddUint64(&b.summary.SuccessCount, 1)
        atomic.AddUint64(&b.summary.SuccessSize, uint64(file.Size))
        file.Result = FileBackupSucceeded
        file.Duration = dur
        if err := os.Chtimes(path, file.ModTime, file.ModTime); err != nil {
            log.WithFields(log.Fields{
                "workerId": workerId,
            }).Error(fmt.Errorf("failed to change file modification time: %s; %w", file.Path, err))
            continue
        }
    }
    return nil
}

func (b *Backup) writeWhatHappened(file *FileWrapper, whatHappened int) {
    file.WhatHappened = whatHappened
    if whatHappened == FileAdded {
        b.summary.addedFiles.Store(file, nil)
        atomic.AddUint64(&b.summary.AddedCount, uint64(1))
        atomic.AddUint64(&b.summary.AddedSize, uint64(file.Size))
        return
    }
    if whatHappened == FileModified {
        b.summary.modifiedFiles.Store(file, nil)
        atomic.AddUint64(&b.summary.ModifiedCount, uint64(1))
        atomic.AddUint64(&b.summary.ModifiedSize, uint64(file.Size))
        return
    }
}

func (b *Backup) compareFileMaps(currentFileMaps []*sync.Map, lastFileMap *sync.Map) error {
    wg := sync.WaitGroup{}
    for i := range currentFileMaps {
        wg.Add(1)
        go func(workerId int) {
            defer wg.Done()
            if err := b.compareFileMap(workerId, lastFileMap, currentFileMaps[workerId]); err != nil {
                log.Error(err)
            }
        }(i)
    }
    wg.Wait()

    // The remaining files in LastFileMap are deleted files.
    lastFileMap.Range(func(k, v interface{}) bool {
        fileWrapper := v.(*FileWrapper)
        fileWrapper.WhatHappened = FileDeleted
        b.summary.deletedFiles.Store(fileWrapper, nil)
        atomic.AddUint64(&b.summary.DeletedCount, 1)
        atomic.AddUint64(&b.summary.DeletedSize, uint64(fileWrapper.Size))
        return true
    })

    // Logging
    log.WithFields(log.Fields{
        "added": b.summary.AddedCount,
        "size":  GetHumanizedSize(b.summary.AddedSize),
    }).Info("# files compared")
    log.WithFields(log.Fields{
        "modified": b.summary.ModifiedCount,
        "size":     GetHumanizedSize(b.summary.ModifiedSize),
    }).Info("# files compared")
    log.WithFields(log.Fields{
        "deleted": b.summary.DeletedCount,
        "size":    GetHumanizedSize(b.summary.DeletedSize),
    }).Info("# files compared")

    b.writeBackupState(Compared)
    return nil
}

func (b *Backup) compareFileMap(workerId int, lastFileMap, currentFileMap *sync.Map) error {
    var count int64
    // t := time.Now()
    currentFileMap.Range(func(k, v interface{}) bool {
        count++
        path := k.(string)
        current := v.(*FileWrapper)

        if val, have := lastFileMap.Load(path); have {
            last := val.(*FileWrapper)
            if last.ModTime.Unix() != current.ModTime.Unix() || last.Size != current.Size {
                // log.WithFields(log.Fields{
                //	"workerId": workerId,
                // }).Debugf("modified: %s", path)
                b.writeWhatHappened(current, FileModified)
            }
            lastFileMap.Delete(path)
            return true
        }

        b.writeWhatHappened(current, FileAdded)
        return true
    })
    // log.WithFields(log.Fields{
    // 	"workerId": workerId,
    // 	"count":    count,
    // 	"duration": time.Since(t).Seconds(),
    // }).Debugf("  - comparison is over")
    return nil
}
