package goback

import (
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (b *Backup) generateFirstBackupData() error {

	// Ready
	log.Info("generating first backup data")
	summary, err := b.newSummary()
	if err != nil {
		return err
	}
	b.summary = summary
	defer func() {
		if err := b.writeSummary(); err != nil {
			log.Error(err)
		}
	}()

	// Reading
	fileMap, err := b.collectFilesToBackup()
	if err != nil {
		return err
	}

	// Write
	if err := b.writeFileMap([]*sync.Map{fileMap}); err != nil {
		return err
	}

	return nil
}

func (b *Backup) writeFileMap(fileMaps []*sync.Map) error {
	defer func() {
		b.summary.LoggingTime = time.Now()
	}()
	// Marshal files
	for _, m := range fileMaps {
		m.Range(func(k, v interface{}) bool {
			file := v.(*File)
			file.Marshal()
			return true
		})
	}

	return b.fileDb.Batch(func(tx *bolt.Tx) error {
		tx.DeleteBucket(BucketFiles)
		b, err := tx.CreateBucket(BucketFiles)
		if err != nil {
			return err
		}

		// Save files
		for _, m := range fileMaps {
			m.Range(func(k, v interface{}) bool {
				path := k.(string)
				file := v.(*File)
				if err := b.Put([]byte(path), file.data); err != nil {
					log.Error(err)
					return false
				}
				return true
			})
		}
		return nil
	})
}

func (b *Backup) collectFilesToBackup() (*sync.Map, error) {
	defer func() {
		b.summary.ReadingTime = time.Now()
		b.summary.ComparisonTime = b.summary.ReadingTime
	}()
	log.Debug("first backup; generating first backup data")
	//fileMap := make(map[string]*File)
	//b.summary.State = BackupRunning
	b.summary.Message = "collecting initialize data"
	fileMap, extensions, sizeDistribution, count, size, err := GetFileMap(b.srcDirArr, b.hashComparision)
	if err != nil {
		return fileMap, err
	}
	b.summary.TotalCount = count
	b.summary.TotalSize = size
	b.summary.ReadingTime = time.Now()
	b.summary.Extensions = extensions
	b.summary.SizeDistribution = sizeDistribution

	return fileMap, nil

	//err := filepath.Walk(b.srcDir, func(path string, file os.FileInfo, err error) error {
	//    if file.IsDir() {
	//        return nil
	//    }
	//
	//    if !file.Mode().IsRegular() {
	//        return nil
	//    }
	//
	//    fi := newFile(path, file.Size(), file.ModTime())
	//    //fi.Result = BackupSuccess
	//    if b.hashComparision {
	//        h, err := GetFileHash(path)
	//        if err != nil {
	//            log.Error(err)
	//            return nil
	//        }
	//        fi.Hash = h
	//    }
	//    if marshal {
	//        fi.Marshal()
	//    }
	//    fileMap.Store(path, fi)
	//    b.summary.TotalCount += 1
	//    b.summary.TotalSize += uint64(file.Size())
	//    return nil
	//})
	//return &fileMap, err
}
