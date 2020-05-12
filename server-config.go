package goback

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

func (s *Server) loadConfig() error {
	err := s.db.View(func(tx *bolt.Tx) error {
		// Config bucket
		b := tx.Bucket(ConfigBucketName)

		data := b.Get(KeyStorage)
		if data != nil {
			var config Config
			if err := json.Unmarshal(data, &config); err != nil {
				return err
			}
			s.config = &config
			return nil
		}

		s.config.Storages = []Storage{
			{Id: 1, Protocol: LocalDisk, Host: "", Port: 0, Username: "", Password: "", Dir: ""},
			{Id: 2, Protocol: Sftp, Host: "", Port: 0, Username: "", Password: "", Dir: ""},
		}

		s.config.Jobs = []Job{
			{Id: 1, BackupType: LocalDisk, SrcDirs: nil, Schedule: "", Ignore: nil, StorageId: 1, Enabled: false, Storage: nil},
			{Id: 2, BackupType: LocalDisk, SrcDirs: nil, Schedule: "", Ignore: nil, StorageId: 1, Enabled: false, Storage: nil},
		}

		// Job setting
		return nil
	})
	if err != nil {
		return err
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
		return b.Put(KeyStorage, data)
	})
}
