package goback

import (
	"encoding/json"
	"fmt"
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
			summary.Stats = nil
			summaries = append(summaries, &summary)
			return nil
		})
		return nil
	})

	return summaries, err
}

func (s *Server) findStats() ([]*Summary, error) {
	summaries, err := s.findSummaries()
	if err != nil {
		return nil, err
	}

	statsMap := make(map[string]*Summary)
	for _, s := range summaries {
		month := s.Date.Format("2006-01")
		dir := s.SrcDir
		key := month + dir
		if _, have := statsMap[key]; !have {
			statsMap[key] = newSummaryStats(s)
		}
		statsMap[key].AddedCount += s.AddedCount
		statsMap[key].AddedSize += s.AddedSize
		statsMap[key].ModifiedCount += s.ModifiedCount
		statsMap[key].ModifiedSize += s.ModifiedSize
		statsMap[key].DeletedCount += s.DeletedCount
		statsMap[key].DeletedSize += s.DeletedSize
		statsMap[key].SuccessCount += s.SuccessCount
		statsMap[key].SuccessSize += s.SuccessSize
		statsMap[key].FailedCount += s.FailedCount
		statsMap[key].FailedSize += s.FailedSize
	}

	var stats []*Summary
	for _, s := range statsMap {
		stats = append(stats, s)
	}
	return stats, err
}

func (s *Server) findSummaryById(id int) (*Summary, error) {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(SummaryBucket)
		if b == nil {
			return ErrorBucketNotFound
		}
		data = b.Get(IntToByte(id))
		return nil
	})
	var summary Summary
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("summary-%d not found", id)
	}

	err = json.Unmarshal(data, &summary)
	if err != nil {
		return nil, err
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

		return b.Put(IntToByte(id), nil)
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
			if err := b.Put(IntToByte(id), data); err != nil {
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
