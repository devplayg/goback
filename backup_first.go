package goback

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (b *Backup) generateFirstBackupData() error {
	log.Debug("generating first backup data")
	id, err := b.issueSummaryId()
	if err != nil {
		return err
	}
	b.summary = newSummary(id, b.srcDir)
	b.summary.ReadingTime = time.Now()
	fileMap, err := b.collectFilesToBackup(true)
	if err != nil {
		return err
	}

	if err := b.writeFileMap(fileMap); err != nil {
		return err
	}
	b.summary.LoggingTime = time.Now()

	log.WithFields(log.Fields{
		"summaryId": b.summary.Id,
		"count":     b.summary.TotalCount,
		"size":      b.summary.TotalSize,
		"execTime":  b.summary.LoggingTime.Sub(b.summary.ReadingTime).Seconds(),
	}).Debug("files found")

	return nil
}

func (b *Backup) writeFileMap(fileMap *sync.Map) error {
	return b.fileDb.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketFiles)
		fileMap.Range(func(key, value interface{}) bool {
			f := value.(*File)
			if err := b.Put([]byte(f.path), f.data); err != nil {
				log.Error(err)
			}
			return true
		})

		return nil
	})
}

func (b *Backup) collectFilesToBackup(marshal bool) (*sync.Map, error) {
	log.Debug("first backup; generating first backup data")
	fileMap := sync.Map{}
	b.summary.State = BackupRunning
	b.summary.Message = "collecting initialize data"
	err := filepath.Walk(b.srcDir, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}

		if !file.Mode().IsRegular() {
			return nil
		}

		fi := newFile(path, file.Size(), file.ModTime())
		fi.Result = BackupSuccess
		if b.hashComparision {
			h, err := b.getFileHash(path)
			if err != nil {
				log.Error(err)
				return nil
			}
			fi.Hash = h
		}
		if marshal {
			fi.Marshal()
		}
		fileMap.Store(path, fi)
		b.summary.TotalCount += 1
		b.summary.TotalSize += uint64(file.Size())
		return nil
	})
	return &fileMap, err
}

func (b *Backup) getFileHash(path string) (string, error) {
	file, err := os.Open(path) // specify your file here
	if err != nil {
		return "", err
	}
	defer file.Close()

	//hash, err := New(key)
	//if err != nil {
	//	fmt.Printf("Failed to create HighwayHash instance: %v", err) // add error handling
	//	return
	//}

	if _, err = io.Copy(b.highwayhash, file); err != nil {
		return "", err
	}

	checksum := b.highwayhash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}
