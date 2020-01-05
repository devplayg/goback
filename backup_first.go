package goback

import (
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (b *Backup) generateFirstBackupData() error {
	log.Debug("generating first backup data")
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
		"count": b.summary.TotalCount,
		"execTime": b.summary.LoggingTime.Sub(b.summary.ReadingTime).Seconds(),
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
		//if b.fileHashComparison {
		//	b, err := ioutil.ReadFile(path)
		//	if err != nil {
		//		log.Error(err)
		//		return nil
		//	}
		//	highwayhash.Sum()
		//}

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
