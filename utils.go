package goback

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devplayg/goutils"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var ErrorBucketNotFound = errors.New("bucket not found")

func IntToByte(v int) []byte {
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
	if !filepath.IsAbs(name) {
		return false
	}
	fi, err := os.Stat(name)
	if !os.IsNotExist(err) {
		return true
	}
	if fi != nil {
		return true
	}
	return false
}

func uniqueStrings(arr []string) []string {
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

func GetHumanizedSize(size uint64) string {
	humanized := humanize.Bytes(size)

	str := fmt.Sprintf("%d B", size)
	if humanized == str {
		return str
	}
	return fmt.Sprintf("%s (%s)", str, humanized)
}

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
		encoded, err = goutils.GobEncode(data)
	} else if encoding == JsonEncoding {
		encoded, err = json.Marshal(data)
	} else {
		return fmt.Errorf("failed to backup data: invalid encoding(%d)", encoding)
	}
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	// Compress
	compressed, err := goutils.Gzip(encoded)
	if err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}
	if err := ioutil.WriteFile(path, compressed, 0644); err != nil {
		return err
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

func newSummaryStats(s *Summary) *Summary {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", s.Date.Format("2006-01")+"-01 00:00:00", time.Local)

	return &Summary{
		//Id:             0,
		//BackupId:       0,
		Date:   t,
		SrcDir: s.SrcDir,
		//BackupDir:      "",
		//BackupType:     0,
		//State:          0,
		//WorkerCount:    0,
		//TotalSize:      0,
		//TotalCount:     0,
		AddedCount:    0,
		AddedSize:     0,
		ModifiedCount: 0,
		ModifiedSize:  0,
		DeletedCount:  0,
		DeletedSize:   0,
		FailedCount:   0,
		FailedSize:    0,
		SuccessCount:  0,
		SuccessSize:   0,
		//ReadingTime:    time.Time{},
		//ComparisonTime: time.Time{},
		//BackupTime:     time.Time{},
		//LoggingTime:    time.Time{},
		//ExecutionTime:  0,
		//Message:        "",
		//Version:        0,
		//Stats:          nil,
		//Keeper:         nil,
		//addedFiles:     nil,
		//modifiedFiles:  nil,
		//deletedFiles:   nil,
		//failedFiles:    nil,
	}
}
