package goback

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/minio/highwayhash"
	"io"
	"os"
	"path/filepath"
)

var ErrorBucketNotFound = errors.New("bucket not found")

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

func GetFileMap(dir string, hashComparision bool) (map[string]*File, int64, error) {
	fileMap := make(map[string]*File)
	var size int64

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
		fileMap[path] = fi
		size += fi.Size
		return nil
	})

	return fileMap, size, err
}
