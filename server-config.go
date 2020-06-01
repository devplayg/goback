package goback

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
)

func (s *Server) loadConfig() error {
	err := s.db.View(func(tx *bolt.Tx) error {
		// Config bucket
		b := tx.Bucket(ConfigBucket)

		data := b.Get(KeyConfig)
		if data != nil {
			var config Config
			if err := json.Unmarshal(data, &config); err != nil {
				return err
			}
			s.config = &config
			return nil
		}

		s.config.Storages = []*Storage{
			{Id: 1, Protocol: LocalDisk, Host: "", Port: 0, Username: "", Password: "", Dir: ""},
			//{Id: 2, Protocol: Sftp, Host: "", Port: 0, Username: "", Password: "", Dir: ""},
		}

		s.config.Jobs = []*Job{
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

func (s *Server) findJobById(jobId int) *Job {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	for i, j := range s.config.Jobs {
		if j.Id == jobId {
			return s.config.Jobs[i]
		}
	}

	return nil
}

func (s *Server) findStorageById(id int) *Storage {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()
	for i, storage := range s.config.Storages {
		if storage.Id == id {
			return s.config.Storages[i]
		}
	}
	return nil
}

func (s *Server) saveConfig(inputChecksum string) error {
	data, err := json.Marshal(s.config)
	checksum := sha256.Sum256(data)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ConfigBucket)

		oldChecksum := b.Get(KeyConfigChecksum)
		_inputChecksum, err := hex.DecodeString(inputChecksum)
		if err != nil {
			return err
		}
		if bytes.Compare(oldChecksum, _inputChecksum) != 0 {
			return fmt.Errorf("checkem error; refresh page")
		}

		if err := b.Put(KeyConfig, data); err != nil {
			return err
		}
		if err := b.Put(KeyConfigChecksum, checksum[:]); err != nil {
			return err
		}

		return nil

	})
}
