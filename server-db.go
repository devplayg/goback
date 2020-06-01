package goback

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

// Thread-safe
func (s *Server) findSummaries() ([]*Summary, error) {
	summaries := make([]*Summary, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(SummaryBucket)
		b.ForEach(func(id, data []byte) error {
			var summary Summary
			if err := json.Unmarshal(data, &summary); err != nil {
				log.Error(err)
				return nil
			}
			summaries = append(summaries, &summary)
			return nil
		})
		return nil
	})

	return summaries, err
}

func (s *Server) findSummaryById(id int) (*Summary, error) {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(SummaryBucket)
		if b == nil {
			return ErrorBucketNotFound
		}
		data = b.Get(iToB(id))
		return nil
	})
	var summary Summary
	if err == nil {
		err := json.Unmarshal(data, &summary)
		if err != nil {
			return nil, err
		}
	}
	return &summary, err
}

func (s *Server) issueDbId(bucketName []byte) (int, error) {
	var id int
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrorBucketNotFound
		}
		newId, _ := b.NextSequence()
		id = int(newId)

		return b.Put(iToB(id), nil)
	})
	return id, err
}

func (s *Server) writeSummaries(results []*Summary) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(SummaryBucket)
		for i := range results {

			newSummaryId, _ := b.NextSequence()
			id := int(newSummaryId)
			results[i].Id = id
			data, err := results[i].Marshal()
			if err != nil {
				log.Error(err)
				continue
			}
			if err := b.Put(iToB(id), data); err != nil {
				log.Error(err)
				continue
			}
		}
		return nil
	})
}

func (s *Server) getDbValue(bucket, key []byte) ([]byte, error) {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return ErrorBucketNotFound
		}
		data = b.Get(key)
		return nil
	})
	return data, err
}
