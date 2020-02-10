package goback

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

func (b *Backup) backupFiles() error {
	fileGroup, count := b.getBackupFileGroup()
	if count < 1 {
		return nil
	}

	// Backup
	for i := range b.keepers {
		if !b.keepers[i].Active() {
			continue
		}
		if err := b.backupFileGroup(fileGroup, i); err != nil {
			log.Error(fmt.Errorf("failed to do backup: %w", err))
			continue
		}
	}

	// Check result
	successVal := 1<<3 - 1
	log.Debug(successVal)
	for i := range fileGroup {
		for j := range fileGroup[i] {
			if fileGroup[i][j].State == successVal {
				atomic.AddUint64(&b.summary.SuccessCount, 1)
				atomic.AddUint64(&b.summary.SuccessSize, uint64(fileGroup[i][j].Size))
				continue
			}

			b.summary.failedFiles.Store(fileGroup[i][j], nil)
			atomic.AddUint64(&b.summary.FailedCount, 1)
			atomic.AddUint64(&b.summary.FailedSize, uint64(fileGroup[i][j].Size))
		}
	}

	return nil
}

func (b *Backup) backupFileGroup(fileGroup [][]*FileWrapper, keeperIdx int) error {
	log.WithFields(log.Fields{
		"protocol": b.keepers[keeperIdx].Description().Protocol,
		"host":     b.keepers[keeperIdx].Description().Host,
		"keepers":  fmt.Sprintf("%d/%d", keeperIdx+1, len(b.keepers)),
	}).Debug("backup")
	defer func() {
		b.writeBackupState(Copied)
		log.WithFields(log.Fields{
			"execTime": b.summary.BackupTime.Sub(b.summary.ComparisonTime).Seconds(),
			"success":  b.summary.SuccessCount,
			"failed":   b.summary.FailedCount,
			"dir":      b.summary.SrcDir,
		}).Info("directory backup done")
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
			// t := time.Now()
			if err := b.backupFilesInGroup(workerId, fileGroup[workerId], keeperIdx); err != nil {
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
func (b *Backup) backupFilesInGroup(workerId int, files []*FileWrapper, keeperIdx int) error {
	for _, file := range files {
		path, dur, err := b.keepers[keeperIdx].keep(file.Path)
		if err != nil { // failed to backup
			file.ExecTime = append(file.ExecTime, 0)
			// 		b.summary.failedFiles.Store(file, nil)
			// 		atomic.AddUint64(&b.summary.FailedCount, 1)
			// 		atomic.AddUint64(&b.summary.FailedSize, uint64(file.Size))
			// 		file.State = FileBackupFailed
			file.Message = append(file.Message, err.Error())
			log.WithFields(log.Fields{
				"workerId": workerId,
			}).Error(fmt.Errorf("failed to backup: %s; %w", file.Path, err))
			continue
		}

		// Success
		// 	atomic.AddUint64(&b.summary.SuccessCount, 1)
		// 	atomic.AddUint64(&b.summary.SuccessSize, uint64(file.Size))
		file.State |= 1 << keeperIdx
		file.ExecTime = append(file.ExecTime, dur)
		file.Message = append(file.Message, "")
		if err := b.keepers[keeperIdx].Chtimes(path, file.ModTime, file.ModTime); err != nil {
			log.WithFields(log.Fields{
				"workerId": workerId,
			}).Error(fmt.Errorf("failed to change file modification time: %s; %w", file.Path, err))
			continue
		}
	}
	return nil
}
