package goback

import (
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

func (b *Backup) startBackup(srcDir string, lastFileMap *sync.Map) error {
	// 1. Issue summary
	b.issueSummary(srcDir, Incremental)

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

	// b.ftpSite = newFtpSite(Sftp, "127.0.0.1", 22, "/backup/", "devplayg", "devplayg123!@#")
	// b.sendChangedFiles()

	// 5. Write result
	// bb, _ := json.MarshalIndent(b.summary, "", "  ")
	// fmt.Println(string(bb))
	if err := b.writeResult(currentFileMaps, lastFileMap); err != nil {
		return err
	}

	return nil
}

// Added or modified files will be backed up
func (b *Backup) getBackupFileGroup() ([][]*FileWrapper, uint64) {
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

	return fileGroup, i
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
	log.WithFields(logrus.Fields{
		"count": b.summary.AddedCount,
		"size":  GetHumanizedSize(b.summary.AddedSize),
	}).Info("added files")
	log.WithFields(logrus.Fields{
		"count": b.summary.ModifiedCount,
		"size":  GetHumanizedSize(b.summary.ModifiedSize),
	}).Info("modified files")
	log.WithFields(logrus.Fields{
		"count": b.summary.DeletedCount,
		"size":  GetHumanizedSize(b.summary.DeletedSize),
	}).Info("deleted files")

	b.writeBackupState(Compared)
	return nil
}

func (b *Backup) compareFileMap(workerId int, lastFileMap, currentFileMap *sync.Map) error {
	var count int64
	//	// t := time.Now()
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
	//	// log.WithFields(log.Fields{
	//	// 	"workerId": workerId,
	//	// 	"count":    count,
	//	// 	"duration": time.Since(t).Seconds(),
	//	// }).Debugf("  - comparison is over")
	return nil
}
