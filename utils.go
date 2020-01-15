package goback

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/minio/highwayhash"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
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

func GetFileMap(dirs []string, hashComparision bool) (*sync.Map, map[string]int64, map[int64]int64, int64, uint64, error) {
	fileMap := sync.Map{}
	extensions := make(map[string]int64)
	sizeDistribution := make(map[int64]int64)

	var size uint64
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
			size += uint64(fi.Size)
			count++

			fileMap.Store(path, fi)

			return nil
		})
		if err != nil {
			return nil, nil, nil, 0, 0, err
		}
	}

	return &fileMap, extensions, sizeDistribution, count, size, nil
}

func GetCurrentFileMaps(dirs []string, workerCount int, hashComparision bool) ([]*sync.Map, map[string]int64, map[int64]int64, int64, uint64, error) {
	fileMaps := make([]*sync.Map, workerCount)
	extensions := make(map[string]int64)
	sizeDistribution := make(map[int64]int64)

	for i := range fileMaps {
		fileMaps[i] = &sync.Map{}
	}

	var size uint64
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
			size += uint64(fi.Size)
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

func EncodeFileMap(fileMap *sync.Map) ([]byte, error) {
	files := make([]*File, 0)
	fileMap.Range(func(k, v interface{}) bool {
		files = append(files, k.(*File))
		return true
	})
	return EncodeFiles(files)
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

func EncodeFiles2(files []*File) ([]byte, error) {
	b, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}

	return Compress(b)
}

func EncodeFiles(files []*File) ([]byte, error) {
	//return json.Marshal(files)
	if len(files) < 1 {
		return nil, nil
	}

	b, err := EncodeToBytes(files)
	if err != nil {
		return nil, err
	}
	compressed, err := Compress(b)
	if err != nil {
		return nil, err
	}
	return compressed, nil
}

func DecodeToFiles(s []byte) ([]*File, error) {
	var files []*File
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	if _, err := zw.Write(data); err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decompress(s []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(s))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if err := reader.Close(); err != nil {
		return nil, err
	}
	return data, nil
}

func EncodeToBytes(p interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BackupFile(tempDir, srcPath string) (string, float64, error) {
	// Set source
	t := time.Now()
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", 0.0, err

	}
	defer srcFile.Close()

	// Set destination
	if runtime.GOOS == "windows" {
		srcPath = strings.ReplaceAll(srcPath, ":", "")
	}
	dstPath := filepath.Join(tempDir, srcPath)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0644); err != nil {
		return "", 0.0, err
	}
	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return "", 0.0, err
	}
	defer dstFile.Close()

	// Copy
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", 0.0, err
	}

	return dstPath, time.Since(t).Seconds(), err
}

func GetChangeFilesDesc(added uint64, modified uint64, deleted uint64) string {
	return fmt.Sprintf("added=%d, modified=%d, deleted=%d", added, modified, deleted)
}

func GetChangeSizeDesc(added uint64, modified uint64, deleted uint64) string {
	return fmt.Sprintf("added=%d(%s), modified=%d(%s), deleted=%d(%s)", added, humanize.Bytes(added), modified, humanize.Bytes(modified), deleted, humanize.Bytes(deleted))
}
