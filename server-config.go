package goback

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

func (s *Server) loadConfig() error {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ConfigBucketName)
		if b == nil {
			return ErrorBucketNotFound
		}
		data = b.Get(ConfigBucketName)
		return nil
	})
	if err == nil {
		var config Config
		err := json.Unmarshal(data, &config)
		if err != nil {
			return err
		}
		s.config = &config
	}

	return nil
}

func (s *Server) saveConfig() error {
	data, err := json.Marshal(s.config)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ConfigBucketName)
		if b == nil {
			return ErrorBucketNotFound
		}
		return b.Put(ConfigBucketName, data)
	})
}
