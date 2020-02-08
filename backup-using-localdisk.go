package goback

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"sync/atomic"
)

func (b *Backup) backupFilesToLocal() error {
	fileGroup, count := b.createBackupFileGroup()
	defer func() {
		b.writeBackupState(Copied)
		if count > 0 {
			log.WithFields(log.Fields{
				"execTime": b.summary.BackupTime.Sub(b.summary.ComparisonTime).Seconds(),
				"success":  b.summary.SuccessCount,
				"failed":   b.summary.FailedCount,
				"dir":      b.summary.SrcDir,
			}).Info("directory backup done")
		}
	}()
	if count < 1 {
		return nil
	}

	log.WithFields(log.Fields{
		"files":   count,
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
			// t := time.Now()
			if err := b.backupFilesInGroup(workerId, fileGroup[workerId]); err != nil {
				log.Error(err)
			}
			// log.WithFields(log.Fields{
			//     "workerId":  workerId,
			//     "processed": len(fileGroup[workerId]),
			//     "duration":  time.Since(t).Seconds(),
			// }).Debug("worker backup done")
		}(i)
	}
	wg.Wait()

	return nil
}

// Thread-safe
func (b *Backup) backupFilesInGroup(workerId int, files []*FileWrapper) error {
	for _, file := range files {
		path, dur, err := b.keeper.Keep(file.Path, b.tempDir)
		if err != nil { // failed to backup
			b.summary.failedFiles.Store(file, nil)
			atomic.AddUint64(&b.summary.FailedCount, 1)
			atomic.AddUint64(&b.summary.FailedSize, uint64(file.Size))
			file.State = FileBackupFailed
			file.Message = err.Error()
			log.WithFields(log.Fields{
				"workerId": workerId,
			}).Error(fmt.Errorf("failed to backup: %s; %w", file.Path, err))
			continue
		}

		// Success
		atomic.AddUint64(&b.summary.SuccessCount, 1)
		atomic.AddUint64(&b.summary.SuccessSize, uint64(file.Size))
		file.State = FileBackupSucceeded
		file.ExecTime = dur
		if err := os.Chtimes(path, file.ModTime, file.ModTime); err != nil {
			log.WithFields(log.Fields{
				"workerId": workerId,
			}).Error(fmt.Errorf("failed to change file modification time: %s; %w", file.Path, err))
			continue
		}
	}
	return nil
}
