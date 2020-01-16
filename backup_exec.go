package goback

import (
	"fmt"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func (b *Backup) startBackup() error {
	// Ready
	summary, err := b.newSummary()
	if err != nil {
		return err
	}
	b.summary = summary
	log.WithFields(log.Fields{
		"summaryId": b.summary.Id,
	}).Debug("New backup has been started")

	if b.fileBackupEnable {
		tempDir, err := ioutil.TempDir(b.dstDir, "backup-")
		if err != nil {
			return err
		}
		b.tempDir = tempDir
	}
	defer func() {
		targetDir := filepath.Join(b.dstDir, b.summary.Date.Format("20060102")+"-"+strconv.FormatInt(b.summary.Id, 10))
		if err := os.Rename(b.tempDir, targetDir); err != nil {
			log.Error(err)
		}
	}()

	// 1. Collect current files
	currentFileMaps, extensions, sizeDistribution, count, size, err := GetCurrentFileMaps(b.srcDirArr, b.workerCount, b.hashComparision)
	if err != nil {
		return nil
	}
	b.summary.ReadingTime = time.Now()
	b.summary.TotalCount = count
	b.summary.TotalSize = size
	b.summary.Extensions = extensions
	b.summary.SizeDistribution = sizeDistribution
	log.WithFields(log.Fields{
		"execTime": b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
		"files":    count,
		"size":     fmt.Sprintf("%d(%s)", size, humanize.Bytes(size)),
	}).Debug("1) collected current files")

	// 2. Compares file maps
	if err := b.CompareFileMaps(currentFileMaps); err != nil {
		return err
	}
	b.summary.ComparisonTime = time.Now()
	log.WithFields(log.Fields{
		"execTime":    b.summary.ComparisonTime.Sub(b.summary.ReadingTime).Seconds(),
		"changeFiles": GetChangeFilesDesc(b.summary.AddedCount, b.summary.ModifiedCount, b.summary.DeletedCount),
		"changeSize":  GetChangeSizeDesc(b.summary.AddedCount, b.summary.ModifiedCount, b.summary.DeletedCount),
	}).Debug("2) comparing files done")

	// 3. Backup added or changed files
	if err := b.BackupFiles(); err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"execTime": b.summary.BackupTime.Sub(b.summary.ComparisonTime).Seconds(),
	}).Debug("3) backup done")

	// 4. Encode changed files
	if err := b.encodedChangedFiles(); err != nil {
		return err
	}

	// 5. Write result
	if err := b.writeResult(currentFileMaps); err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"execTime": b.summary.LoggingTime.Sub(b.summary.ComparisonTime).Seconds(),
	}).Debug("4) logging done")

	if err := b.writeSummary(); err != nil {
		return err
	}

	return nil
}

func (b *Backup) encodedChangedFiles() error {
	var err error

	if b.addedData, err = EncodeFileMap(b.addedFiles); err != nil {
		return err
	}
	if b.modifiedData, err = EncodeFileMap(b.modifiedFiles); err != nil {
		return err
	}
	if b.deletedData, err = EncodeFileMap(b.deletedFiles); err != nil {
		return err
	}
	if b.failedData, err = EncodeFileMap(b.failedFiles); err != nil {
		return err
	}

	return nil
}

func (b *Backup) backupFileGroup() ([][]*File, error) {
	fileGroup := make([][]*File, b.workerCount)
	for i := range fileGroup {
		fileGroup[i] = make([]*File, 0)
	}

	i := 0
	b.addedFiles.Range(func(k, v interface{}) bool {
		file := k.(*File)
		workerId := i % b.workerCount
		fileGroup[workerId] = append(fileGroup[workerId], file)
		i++
		return true
	})
	b.modifiedFiles.Range(func(k, v interface{}) bool {
		file := k.(*File)
		workerId := i % b.workerCount
		fileGroup[workerId] = append(fileGroup[workerId], file)
		i++
		return true
	})

	return fileGroup, nil

}

func (b *Backup) BackupFiles() error {
	fileGroup, err := b.backupFileGroup()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for i := range fileGroup {
		if len(fileGroup[i]) < 1 {
			continue
		}

		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			t := time.Now()
			if err := b.backupFiles(workerId, fileGroup[workerId]); err != nil {
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
	b.summary.BackupTime = time.Now()

	return nil
}

func (b *Backup) backupFiles(workerId int, files []*File) error {
	for _, f := range files {
		err := b.backupFile(f)
		if err != nil {
			log.Errorf("failed to backup: %s; %s", f.Path, err.Error())
			continue
		}
	}
	return nil
}

func (b *Backup) backupFile(file *File) error {
	//return nil // wondory
	path, dur, err := BackupFile(b.tempDir, file.Path)
	if err != nil {
		b.failedFiles.Store(file, nil)
		atomic.AddUint64(&b.summary.BackupFailureCount, uint64(1))
		file.Result = Failure
		file.Message = err.Error()
		return err
	}

	atomic.AddUint64(&b.summary.BackupSuccessCount, uint64(1))
	file.Result = Success
	file.Duration = dur
	if err := os.Chtimes(path, file.ModTime, file.ModTime); err != nil {
		return err
	}

	return nil
}

func (b *Backup) writeWhatHappened(file *File, whatHappened int) {
	file.WhatHappened = whatHappened
	if whatHappened == FileAdded {
		b.addedFiles.Store(file, nil)
		atomic.AddUint64(&b.summary.AddedCount, uint64(1))
		atomic.AddUint64(&b.summary.AddedSize, uint64(file.Size))
		return
	}
	if whatHappened == FileModified {
		b.modifiedFiles.Store(file, nil)
		atomic.AddUint64(&b.summary.ModifiedCount, uint64(1))
		atomic.AddUint64(&b.summary.ModifiedSize, uint64(file.Size))
		return
	}
	if whatHappened == FileDeleted {
		b.deletedFiles.Store(file, nil)
		atomic.AddUint64(&b.summary.DeletedCount, uint64(1))
		atomic.AddUint64(&b.summary.DeletedSize, uint64(file.Size))
		return
	}
}

func (b *Backup) CompareFileMaps(currentFileMaps []*sync.Map) error {
	wg := sync.WaitGroup{}
	for i := range currentFileMaps {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			if err := b.compareFileMap(workerId, b.lastFileMap, currentFileMaps[workerId]); err != nil {
				log.Error(err)
			}
		}(i)
	}
	wg.Wait()

	// for _, v := range lastFileMap {
	b.lastFileMap.Range(func(k, v interface{}) bool {
		file := v.(*File)
		b.writeWhatHappened(file, FileDeleted)
		return true
	})

	return nil
}

func (b *Backup) compareFileMap(workerId int, lastFileMap, myMap *sync.Map) error {
	var count int64
	//t := time.Now()
	myMap.Range(func(k, v interface{}) bool {
		count++
		path := k.(string)
		current := v.(*File)

		if val, have := lastFileMap.Load(path); have {
			last := val.(*File)
			if last.ModTime.Unix() != current.ModTime.Unix() || last.Size != current.Size {
				//log.WithFields(log.Fields{
				//	"workerId": workerId,
				//}).Debugf("modified: %s", path)
				b.writeWhatHappened(current, FileModified)
			}
			lastFileMap.Delete(path)
			return true
		}

		b.writeWhatHappened(current, FileAdded)
		return true
	})
	//log.WithFields(log.Fields{
	//	"workerId": workerId,
	//	"count":    count,
	//	"duration": time.Since(t).Seconds(),
	//}).Debugf("comparison is over")
	return nil
}
