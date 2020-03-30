package goback

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

// Thread-safe
func (s *Server) findSummaries() ([]*Summary, error) {
	summaries := make([]*Summary, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(SummaryBucketName)
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
		b := tx.Bucket(SummaryBucketName)
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
