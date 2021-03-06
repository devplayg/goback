package goback

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

func (b *Backup) backupFiles() error {
	fileGroup, count := b.getBackupFileGroup()
	if count < 1 {
		b.summary.BackupTime = b.summary.ComparisonTime
		return nil
	}

	// Backup
	if err := b.backupFileGroup(fileGroup); err != nil {
		log.Error(fmt.Errorf("failed to do backup: %w", err))
	}

	return nil
}

func (b *Backup) backupFileGroup(fileGroup [][]*FileWrapper) error {
	defer func() {
		b.writeBackupState(Copied)
	}()

	// 	log.WithFields(log.Fields{
	// 		"files":   count,
	// 		"workers": b.workerCount,
	// 	}).Info("running backup..")
	wg := sync.WaitGroup{}
	for i := range fileGroup {
		if len(fileGroup[i]) < 1 {
			continue
		}

		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

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
		path, dur, err := b.keeper.keep(file.Path)
		if err != nil { // failed to backup
			b.summary.failedFiles.Store(file, nil)
			atomic.AddUint64(&b.summary.FailedCount, 1)
			atomic.AddUint64(&b.summary.FailedSize, uint64(file.Size))
			file.State = FileBackupFailed
			file.Message = err.Error()
			log.WithFields(logrus.Fields{
				"workerId": workerId,
			}).Error(fmt.Errorf("failed to backup: %s; %w", file.Path, err))
			continue
		}

		// Success
		atomic.AddUint64(&b.summary.SuccessCount, 1)
		atomic.AddUint64(&b.summary.SuccessSize, uint64(file.Size))
		file.State = FileBackupSucceeded
		file.ExecTime = dur
		if err := b.keeper.Chtimes(path, file.ModTime, file.ModTime); err != nil {
			log.WithFields(logrus.Fields{
				"workerId": workerId,
			}).Error(fmt.Errorf("failed to change file modification time: %s; %w", file.Path, err))
			continue
		}
	}
	return nil
}
