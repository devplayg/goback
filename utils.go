package goback

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/minio/highwayhash"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var ErrorBucketNotFound = errors.New("bucket not found")
var fileSizeCategories = []int64{
	1, 10, 50, 100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000, 5000000, 10000000,
}

func isValidDir(dir string) error {
	if len(dir) < 1 {
		return errors.New("empty directory")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return errors.New("directory not found: " + dir)
	}

	fi, err := os.Lstat(dir)
	if err != nil {
		return err
	}

	if !fi.Mode().IsDir() {
		return errors.New("invalid source directory: " + fi.Name())
	}

	return nil
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func UnmarshalSummary(data []byte) (*Summary, error) {
	var summary Summary
	err := json.Unmarshal(data, &summary)
	return &summary, err
}

func GetFileHash(path string) (string, error) {

	highwayhash, err := highwayhash.New(HashKey[:])
	if err != nil {
		return "", err
	}
	file, err := os.Open(path) // specify your file here
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err = io.Copy(highwayhash, file); err != nil {
		return "", err
	}

	checksum := highwayhash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}

func GetFileMap(dirs []string, hashComparision bool) (*sync.Map, map[string]int64, map[int64]int64, int64, int64, error) {
	fileMap := sync.Map{}
	extensions := make(map[string]int64)
	sizeDistribution := make(map[int64]int64)

	var size int64
	var count int64

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
			if file.IsDir() {
				return nil
			}

			if !file.Mode().IsRegular() {
				return nil
			}

			fi := newFile(path, file.Size(), file.ModTime())
			if hashComparision {
				h, err := GetFileHash(path)
				if err != nil {
					return err
				}
				fi.Hash = h
			}

			// Statistics
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if len(ext) > 0 {
				extensions[ext]++
			} else {
				extensions["__OTHERS__"]++
			}
			sizeDistribution[GetFileSizeCategory(file.Size())]++
			size += fi.Size
			count++

			return nil
		})
		if err != nil {
			return nil, nil, nil, 0, 0, err
		}
	}

	return &fileMap, extensions, sizeDistribution, count, size, nil
}

func GetCurrentFileMaps(dirs []string, workerCount int, hashComparision bool) ([]*sync.Map, map[string]int64, map[int64]int64, int64, int64, error) {
	fileMaps := make([]*sync.Map, workerCount)
	extensions := make(map[string]int64)
	sizeDistribution := make(map[int64]int64)

	for i := range fileMaps {
		fileMaps[i] = &sync.Map{}
	}

	var size int64
	var count int64

	for _, dir := range dirs {
		i := 0
		err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
			if file.IsDir() {
				return nil
			}

			if !file.Mode().IsRegular() {
				return nil
			}

			fi := newFile(path, file.Size(), file.ModTime())
			if hashComparision {
				h, err := GetFileHash(path)
				if err != nil {
					return err
				}
				fi.Hash = h
			}

			// Statistics
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if len(ext) > 0 {
				extensions[ext]++
			} else {
				extensions["__OTHERS__"]++
			}
			sizeDistribution[GetFileSizeCategory(file.Size())]++
			size += fi.Size
			count++

			// Distribute works
			workerId := i % workerCount
			fileMaps[workerId].Store(path, fi)
			i++

			return nil
		})
		if err != nil {
			return nil, nil, nil, 0, 0, err
		}
	}

	return fileMaps, extensions, sizeDistribution, count, size, nil
}

func GetFileSizeCategory(size int64) int64 {
	for i := range fileSizeCategories {
		if size <= fileSizeCategories[i] {
			return fileSizeCategories[i]
		}
	}
	return -1

}

func CompareFileMaps(lastFileMap, currentFileMap map[string]*File) ([]*File, []*File, []*File, error) {
	added := make([]*File, 0)
	deleted := make([]*File, 0)
	modified := make([]*File, 0)
	for path, current := range currentFileMap {
		if last, had := lastFileMap[path]; had {
			if last.ModTime.Unix() != current.ModTime.Unix() || last.Size != current.Size {
				log.Debugf("modified: %s", path)
				current.WhatHappened = FileModified
				modified = append(modified, current)
			}
			delete(lastFileMap, path)

		} else {
			log.Debugf("added: %s", path)
			current.WhatHappened = FileAdded
			added = append(added, current)
		}

		current.Marshal()
	}
	for _, file := range lastFileMap {
		file.WhatHappened = FileDeleted
		deleted = append(deleted, file)
	}

	log.WithFields(log.Fields{
		"added":    len(added),
		"modified": len(modified),
		"deleted":  len(deleted),
	}).Debugf("total %d files; comparision result", len(currentFileMap))

	return added, modified, deleted, nil
}

func IsEqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
