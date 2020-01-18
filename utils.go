package goback

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/minio/highwayhash"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DefaultDateFormat = "2006-01-02 15:04:05"
)

var ErrorBucketNotFound = errors.New("bucket not found")

const (
	kB = 1000
	MB = 1000000
	GB = 1000000000
)

var fileSizeCategories = []int64{
	1 * kB,
	5 * kB,
	10 * kB,
	50 * kB,
	100 * kB,
	500 * kB,

	1 * MB,
	5 * MB,
	10 * MB,
	50 * MB,
	100 * MB,
	500 * MB,

	1 * GB,
	5 * GB,
	10 * GB,
	50 * GB,
	100 * GB,
	500 * GB,
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

			fi := NewFileWrapper(path, file.Size(), file.ModTime())
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

func GetFileSizeCategory(size int64) int64 {
	for i := range fileSizeCategories {
		if size <= fileSizeCategories[i] {
			return fileSizeCategories[i]
		}
	}
	return -1

}

func IsEqualStringSlices(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
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

// func EncodeToBytes(p interface{}) ([]byte, error) {
// 	buf := bytes.Buffer{}
// 	enc := gob.NewEncoder(&buf)
// 	if err := enc.Encode(p); err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

func BackupFile(srcPath, tempDir string) (string, float64, error) {
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

//
// func GetChangeFilesDesc(added uint64, modified uint64, deleted uint64) string {
// 	return fmt.Sprintf("added=%d, modified=%d, deleted=%d", added, modified, deleted)
// }
//
// func GetChangeSizeDesc(added uint64, modified uint64, deleted uint64) string {
// 	return fmt.Sprintf("added=%d(%s), modified=%d(%s), deleted=%d(%s)", added, humanize.Bytes(added), modified, humanize.Bytes(modified), deleted, humanize.Bytes(deleted))
// }

func GetHumanizedSize(size uint64) string {
	humanized := humanize.Bytes(size)

	str := fmt.Sprintf("%d B", size)
	if humanized == str {
		return str
	}
	return fmt.Sprintf("%s (%s)", str, humanized)
}

// func CreateFileIfNotExists(path string) error {
// 	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	if err := f.Close(); err != nil {
// 		return nil
// 	}
// 	return nil
// }
