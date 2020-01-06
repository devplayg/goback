package goback

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"time"
)

func (b *Backup) generateFirstBackupData() error {
	log.Debug("generating first backup data")
	summary, err := b.newSummary()
	if err != nil {
		return err
	}
	b.summary = summary

	fileMap, err := b.collectFilesToBackup()
	if err != nil {
		return err
	}
	b.summary.ReadingTime = time.Now()
	b.summary.ComparisonTime = b.summary.ReadingTime

	if err := b.writeFileMap(fileMap); err != nil {
		return err
	}
	b.summary.LoggingTime = time.Now()

	if err := b.writeSummary(); err != nil {
		return err
	}

	return nil
}

func (b *Backup) writeFileMap(fileMap map[string]*File) error {
	return b.fileDb.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketFiles)
		if b == nil {
			return ErrorBucketNotFound
		}

		for path, f := range fileMap {
			if err := b.Put([]byte(path), f.data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *Backup) collectFilesToBackup() (map[string]*File, error) {
	log.Debug("first backup; generating first backup data")
	//fileMap := make(map[string]*File)
	b.summary.State = BackupRunning
	b.summary.Message = "collecting initialize data"
	fileMap, size, err := GetFileMap(b.srcDir, b.hashComparision)
	if err != nil {
		return nil, err
	}
	b.summary.ReadingTime = time.Now()
	b.summary.TotalCount = len(fileMap)
	b.summary.TotalSize = size

	for path, _ := range fileMap {
		fileMap[path].Marshal()
	}

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

func (b *Backup) writeSummary() error {
	log.WithFields(log.Fields{
		"summaryId": b.summary.Id,
		"count":     b.summary.TotalCount,
		"size":      b.summary.TotalSize,
		"execTime":  b.summary.LoggingTime.Sub(b.summary.ReadingTime).Seconds(),
	}).Debug("files found")

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketSummary)
		if bucket == nil {
			return ErrorBucketNotFound
		}
		//id, _ := b.NextSequence()
		//summaryId = int64(id)
		data, err := json.Marshal(b.summary)
		if err != nil {
			return err
		}
		return bucket.Put(Int64ToBytes(b.summary.Id), data)
	})
}
