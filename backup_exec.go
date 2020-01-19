package goback

// import (
// 	"fmt"
// 	"github.com/dustin/go-humanize"
// 	log "github.com/sirupsen/logrus"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )
//
// func (b *Backup) startBackup() error {
// 	// 1. Collect files in source directories
// 	currentFileMaps, err := b.getCurrentFileMaps()
// 	if err != nil {
// 		return nil
// 	}
//
// 	// 2. Compares file maps
// 	if err := b.CompareFileMaps(currentFileMaps); err != nil {
// 		return err
// 	}
//
// 	// 3. Backup added or changed files
// 	if err := b.BackupFiles(); err != nil {
// 		return err
// 	}
//
// 	// 4. Write result
// 	if err := b.writeResult(currentFileMaps); err != nil {
// 		return err
// 	}
//
// 	// 5. Write summary
// 	if err := b.writeSummary(); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (b *Backup) getCurrentFileMaps() ([]*sync.Map, error) {
// 	fileMaps := make([]*sync.Map, b.workerCount)
// 	b.summary.Extensions = make(map[string]int64)
// 	b.summary.SizeDistribution = make(map[int64]int64)
//
// 	for i := range fileMaps {
// 		fileMaps[i] = &sync.Map{}
// 	}
//
// 	for _, dir := range b.srcDirs {
// 		i := 0
// 		err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
// 			if file.IsDir() {
// 				return nil
// 			}
//
// 			if !file.Mode().IsRegular() {
// 				return nil
// 			}
//
// 			fi := NewFileWrapper(path, file.Size(), file.ModTime())
// 			if b.hashComparision {
// 				h, err := GetFileHash(path)
// 				if err != nil {
// 					return err
// 				}
// 				fi.Hash = h
// 			}
//
// 			// Statistics
// 			b.summary.addExtension(file.Name())
// 			b.summary.SizeDistribution[GetFileSizeCategory(file.Size())]++
// 			b.summary.TotalSize += uint64(fi.Size)
// 			b.summary.TotalCount++
//
// 			// Distribute works
// 			workerId := i % b.workerCount
// 			fileMaps[workerId].Store(path, fi)
// 			i++
//
// 			return nil
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	b.summary.ReadingTime = time.Now()
// 	log.WithFields(log.Fields{
// 		"execTime": b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
// 		"files":    b.summary.TotalCount,
// 		"size":     fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
// 	}).Info("current files loaded")
//
// 	return fileMaps, nil
// }
//
// func (b *Backup) BackupFiles() error {
// 	fileGroup, count, err := b.backupFileGroup()
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		b.summary.BackupTime = time.Now()
// 		if count > 0 {
// 			log.WithFields(log.Fields{
// 				"execTim" +
// 					"" +
// 					"e": b.summary.BackupTime.Sub(b.summary.ComparisonTime).Seconds(),
// 				"success": b.summary.SuccessCount,
// 				"failed":  b.summary.FailedCount,
// 			}).Info("backup report")
// 		}
// 	}()
//
// 	log.WithFields(log.Fields{
// 		"workers": b.workerCount,
// 	}).Info("running backup..")
// 	wg := sync.WaitGroup{}
// 	for i := range fileGroup {
// 		if len(fileGroup[i]) < 1 {
// 			continue
// 		}
//
// 		wg.Add(1)
// 		go func(workerId int) {
// 			defer wg.Done()
// 			t := time.Now()
// 			if err := b.backupFiles(workerId, fileGroup[workerId]); err != nil {
// 				log.Error(err)
// 			}
// 			log.WithFields(log.Fields{
// 				"workerId":  workerId,
// 				"processed": len(fileGroup[workerId]),
// 				"duration":  time.Since(t).Seconds(),
// 			}).Debug("backup done")
// 		}(i)
// 	}
// 	wg.Wait()
//
// 	// b.failedFiles.Range(func(key, value interface{}) bool {
// 	// 	file := key.(*FileWrapper)
// 	// 	return true
// 	// })
//
// 	return nil
// }
//
// func (b *Backup) backupFileGroup() ([][]*FileWrapper, uint64, error) {
// 	fileGroup := make([][]*FileWrapper, b.workerCount)
// 	for i := range fileGroup {
// 		fileGroup[i] = make([]*FileWrapper, 0)
// 	}
//
// 	var i uint64 = 0
// 	b.addedFiles.Range(func(k, v interface{}) bool {
// 		file := k.(*FileWrapper)
// 		workerId := i % uint64(b.workerCount)
// 		fileGroup[workerId] = append(fileGroup[workerId], file)
// 		i++
//
// 		return true
// 	})
// 	b.modifiedFiles.Range(func(k, v interface{}) bool {
// 		file := k.(*FileWrapper)
// 		workerId := i % uint64(b.workerCount)
// 		fileGroup[workerId] = append(fileGroup[workerId], file)
// 		i++
//
// 		return true
// 	})
//
// 	return fileGroup, i, nil
// }
//
// // Thread-safe
// func (b *Backup) backupFiles(workerId int, files []*FileWrapper) error {
// 	for _, file := range files {
// 		path, dur, err := BackupFile(file.Path, b.tempDir)
// 		if err != nil { // failed to backup
// 			b.failedFiles.Store(file, nil)
// 			atomic.AddUint64(&b.summary.FailedCount, 1)
// 			atomic.AddUint64(&b.summary.FailedSize, uint64(file.Size))
// 			file.Result = FileBackupFailed
// 			file.Message = err.Error()
// 			log.WithFields(log.Fields{
// 				"workerId": workerId,
// 			}).Error(fmt.Errorf("failed to backup: %s; %w", file.Path, err))
// 			continue
// 		}
//
// 		// Success
// 		atomic.AddUint64(&b.summary.SuccessCount, 1)
// 		atomic.AddUint64(&b.summary.SuccessSize, uint64(file.Size))
// 		file.Result = FileBackupSucceeded
// 		file.Duration = dur
// 		if err := os.Chtimes(path, file.ModTime, file.ModTime); err != nil {
// 			log.WithFields(log.Fields{
// 				"workerId": workerId,
// 			}).Error(fmt.Errorf("failed to change file modification time: %s; %w", file.Path, err))
// 			continue
// 		}
// 	}
// 	return nil
// }
//
// func (b *Backup) writeWhatHappened(file *FileWrapper, whatHappened int) {
// 	file.WhatHappened = whatHappened
// 	if whatHappened == FileAdded {
// 		b.addedFiles.Store(file, nil)
// 		atomic.AddUint64(&b.summary.AddedCount, uint64(1))
// 		atomic.AddUint64(&b.summary.AddedSize, uint64(file.Size))
// 		return
// 	}
// 	if whatHappened == FileModified {
// 		b.modifiedFiles.Store(file, nil)
// 		atomic.AddUint64(&b.summary.ModifiedCount, uint64(1))
// 		atomic.AddUint64(&b.summary.ModifiedSize, uint64(file.Size))
// 		return
// 	}
// }
//
// func (b *Backup) CompareFileMaps(currentFileMaps []*sync.Map) error {
// 	wg := sync.WaitGroup{}
// 	for i := range currentFileMaps {
// 		wg.Add(1)
// 		go func(workerId int) {
// 			defer wg.Done()
// 			if err := b.compareFileMap(workerId, b.lastFileMap, currentFileMaps[workerId]); err != nil {
// 				log.Error(err)
// 			}
// 		}(i)
// 	}
// 	wg.Wait()
//
// 	// The remaining files in LastFileMap are deleted files.
// 	b.lastFileMap.Range(func(k, v interface{}) bool {
// 		file := v.(*File)
//
// 		fileWrapper := FileWrapper{
// 			File:         file,
// 			WhatHappened: FileDeleted,
// 			Result:       0,
// 			Duration:     0,
// 			Message:      "",
// 		}
// 		b.deletedFiles.Store(&fileWrapper, nil)
// 		atomic.AddUint64(&b.summary.DeletedCount, 1)
// 		atomic.AddUint64(&b.summary.DeletedSize, uint64(file.Size))
// 		return true
// 	})
//
// 	// Logging
// 	log.Info(strings.Repeat("=", 50))
// 	log.WithFields(log.Fields{
// 		"added": b.summary.AddedCount,
// 		"size":  GetHumanizedSize(b.summary.AddedSize),
// 	}).Info("# files compared")
// 	log.WithFields(log.Fields{
// 		"modified": b.summary.ModifiedCount,
// 		"size":     GetHumanizedSize(b.summary.ModifiedSize),
// 	}).Info("# files compared")
// 	log.WithFields(log.Fields{
// 		"deleted": b.summary.DeletedCount,
// 		"size":    GetHumanizedSize(b.summary.DeletedSize),
// 	}).Info("# files compared")
// 	log.Info(strings.Repeat("=", 50))
// 	b.summary.ComparisonTime = time.Now()
//
// 	return nil
// }
//
// func (b *Backup) compareFileMap(workerId int, lastFileMap, currentFileMap *sync.Map) error {
// 	var count int64
// 	// t := time.Now()
// 	currentFileMap.Range(func(k, v interface{}) bool {
// 		count++
// 		path := k.(string)
// 		current := v.(*FileWrapper)
//
// 		if val, have := lastFileMap.Load(path); have {
// 			last := val.(*File)
// 			if last.ModTime.Unix() != current.ModTime.Unix() || last.Size != current.Size {
// 				// log.WithFields(log.Fields{
// 				//	"workerId": workerId,
// 				// }).Debugf("modified: %s", path)
// 				b.writeWhatHappened(current, FileModified)
// 			}
// 			lastFileMap.Delete(path)
// 			return true
// 		}
//
// 		b.writeWhatHappened(current, FileAdded)
// 		return true
// 	})
// 	// log.WithFields(log.Fields{
// 	// 	"workerId": workerId,
// 	// 	"count":    count,
// 	// 	"duration": time.Since(t).Seconds(),
// 	// }).Debugf("  - comparison is over")
// 	return nil
// }
