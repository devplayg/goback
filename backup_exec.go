package goback

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

//
//func (b *Backup) getLastSummary() (*Summary, error) {
//	lastSummary, err := b.getLastSummary()
//	if err != nil {
//		return lastSummary, nil, err
//	}
//
//	return lastSummary, nil, nil
//}
//

func (b *Backup) startBackup() error {
	return nil
}

func (b *Backup) getLastSummary() (*Summary, error) {
	var lastSummaryData []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketSummary)
		c := b.Cursor()
		_, lastSummaryData = c.Last()
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(lastSummaryData) == 0 {
		return nil, nil
	}

	var summary Summary
	err = json.Unmarshal(lastSummaryData, &summary)
	return &summary, err
}
