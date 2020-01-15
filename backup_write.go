package goback

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

func (b *Backup) writeResult(currentFileMaps []*sync.Map) error {
	if err := b.writeBackupResult(); err != nil {
		return err
	}

	if err := b.writeFileMap(currentFileMaps); err != nil {
		return err
	}

	if err := b.writeSummary(); err != nil {
		return err
	}

	return nil
}

func (b *Backup) writeSummary() error {
	b.summary.ExecutionTime = time.Since(b.summary.Date).Seconds()
	b.logSummary()

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketSummary)
		if bucket == nil {
			return ErrorBucketNotFound
		}
		data, err := json.Marshal(b.summary)
		if err != nil {
			return err
		}
		return bucket.Put(Int64ToBytes(b.summary.Id), data)
	})
}

func (b *Backup) logSummary() {
	log.WithFields(log.Fields{
		"summaryId":   b.summary.Id,
		"files":       b.summary.TotalCount,
		"totalSize":   fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
		"changeFiles": GetChangeFilesDesc(b.summary.AddedCount, b.summary.ModifiedCount, b.summary.DeletedCount),
		"changeSize":  GetChangeSizeDesc(b.summary.AddedCount, b.summary.ModifiedCount, b.summary.DeletedCount),
		"execTime":    time.Since(b.summary.Date).Seconds(),
	}).Info("summary")

}

func (b *Backup) writeBackupResult() error {
	return b.db.Batch(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BackupPrefixStr + strconv.FormatInt(b.summary.Id, 10)))
		if err != nil {
			return err
		}

		// Write added files
		if err := bucket.Put(BucketAdded, b.addedData); err != nil {
			return err
		}

		// Write modified files
		if err := bucket.Put(BucketModified, b.modifiedData); err != nil {
			return err
		}

		// Write deleted files
		if err := bucket.Put(BucketDeleted, b.deletedData); err != nil {
			return err
		}

		// Write files that failed to back up
		if err := bucket.Put(BucketFailed, b.failedData); err != nil {
			return err
		}

		return nil
	})
}
