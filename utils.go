package goback

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/golibs/converter"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrorBucketNotFound = errors.New("bucket not found")

func iToB(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func IsValidDir(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		return "", errors.New("directory not found: " + absDir)
	}

	fi, err := os.Lstat(absDir)
	if err != nil {
		return "", err
	}

	if !fi.Mode().IsDir() {
		return "", errors.New("invalid source directory: " + fi.Name())
	}

	return absDir, nil
}

func DirExists(name string) bool {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return true
	}
	if filepath.IsAbs(name) {
		return true
	}
	return false
}

func UniqueStrings(arr []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0)
	for _, e := range arr {
		e = strings.TrimSpace(e)
		if _, value := keys[e]; !value {
			keys[e] = true
			list = append(list, e)
		}
	}

	return list
}

//
// func GetFileHash(path string) (string, error) {
// 	highwayhash, err := highwayhash.New(HashKey[:])
// 	if err != nil {
// 		return "", err
// 	}
// 	file, err := os.Open(path) // specify your file here
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()
//
// 	if _, err = io.Copy(highwayhash, file); err != nil {
// 		return "", err
// 	}
//
// 	checksum := highwayhash.Sum(nil)
// 	return hex.EncodeToString(checksum), nil
// }

func GetHumanizedSize(size uint64) string {
	humanized := humanize.Bytes(size)

	str := fmt.Sprintf("%d B", size)
	if humanized == str {
		return str
	}
	return fmt.Sprintf("%s (%s)", str, humanized)
}

// func LoadOrCreateDatabase(path string) (*os.File, error) {
//	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
//	if err != nil {
//		return nil, err
//	}
//	return db, nil
// }

func NewSizeDistMap() map[int64]*SizeDistStats {
	m := make(map[int64]*SizeDistStats)
	for _, sizeDist := range fileSizeCategories {
		m[sizeDist] = NewSizeDistStats(sizeDist, 0, 0)
	}
	return m
}

func WriteBackupData(data interface{}, path string, encoding int) error {
	var encoded []byte
	var err error

	// Encode
	if encoding == GobEncoding {
		encoded, err = converter.EncodeToBytes(data)
	} else if encoding == JsonEncoding {
		encoded, err = json.Marshal(data)
	} else {
		return fmt.Errorf("failed to backup data: invalid encoding(%d)", encoding)
	}
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	// Compress
	compressed, err := compress.Compress(encoded, compress.GZIP)
	if err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}
	if err := ioutil.WriteFile(path, compressed, 0644); err != nil {
		return err
	}
	return nil
}

func LoadBackupData(path string, output interface{}, encoding int) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("database not found: %s", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	decompressed, err := compress.Decompress(data, compress.GZIP)
	if err != nil {
		return err
	}

	if encoding == GobEncoding {
		if err := converter.DecodeFromBytes(decompressed, output); err != nil {
			return err
		}
	} else if encoding == JsonEncoding {
		if err := json.Unmarshal(decompressed, output); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("failed to restore  data: invalid encoding(%d)", encoding)
	}

	return nil
}

func GetFileNameKey(name string, size int64) string {
	return fmt.Sprintf("%s-%d", name, size)
}

func FindProperBackupDirName(dir string) string {
	i := 0
	for {
		var d string
		if i < 1 {
			d = dir
		} else {
			d = dir + "-" + strconv.Itoa(i)
		}
		if _, err := os.Stat(d); os.IsNotExist(err) {
			return d
		}
		i++
	}
}

func DecodeSummaries(data []byte) ([]*Summary, int, int, error) {
	if len(data) < 1 {
		return make([]*Summary, 0), 0, 0, nil
	}

	var summaries []*Summary
	err := converter.DecodeFromBytes(data, &summaries)

	lastIdx := len(summaries) - 1
	return summaries, summaries[lastIdx].BackupId, summaries[lastIdx].Id, err
}
