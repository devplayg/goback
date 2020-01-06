package goback

import (
	"github.com/boltdb/bolt"
)

func GetLastDbData(db *bolt.DB, bucketName []byte) ([]byte, []byte, error) {
	var key []byte
	var val []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrorBucketNotFound
		}
		c := b.Cursor()
		key, val = c.Last()
		return nil
	})
	return key, val, err
}

func IssueDbInt64Id(db *bolt.DB, bucketName []byte) (int64, error) {
	var id int64
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrorBucketNotFound
		}
		newId, _ := b.NextSequence()
		id = int64(newId)
		return nil
	})
	return id, err
}
