package goback

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"
)

type Summary struct {
	Id        int64     `json:"id"`
	Date      time.Time `json:"date"`
	SrcDirArr []string  `json:"srcDirs"`
	DstDir    string    `json:"dstDir"`
	// State     int       `json:"state"`

	WorkerCount int    `json:"workerCount"`
	TotalSize   uint64 `json:"totalSize"`
	TotalCount  int64  `json:"totalCount"`

	// Thread-safe
	AddedCount    uint64 `json:"countAdded"`
	AddedSize     uint64 `json:"sizeAdded"`
	ModifiedCount uint64 `json:"countModified"`
	ModifiedSize  uint64 `json:"sizeModified"`
	DeletedCount  uint64 `json:"countDeleted"`
	DeletedSize   uint64 `json:"sizeDeleted"`
	FailedCount   uint64 `json:"countFailed"`
	FailedSize    uint64 `json:"sizeFailed"`
	SuccessCount  uint64 `json:"countSuccess"`
	SuccessSize   uint64 `json:"sizeSuccess"`

	Extensions       map[string]int64 `json:"extensions"`
	SizeDistribution map[int64]int64  `json:"sizeDistribution"`

	Message string `json:"message"`
	Version int    `json:"version"`

	ReadingTime    time.Time `json:"readingTime"`    // Step 1
	ComparisonTime time.Time `json:"comparisonTime"` // Step 2
	BackupTime     time.Time `json:"backupTime"`     // Step 3
	LoggingTime    time.Time `json:"loggingTime"`    // Step 4
	ExecutionTime  float64   `json:"execTime"`
}

func (s *Summary) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func NewSummary(id int64, srcDirs []string, dstDir string, workCount, version int) *Summary {
	return &Summary{
		Id:          id,
		Date:        time.Now(),
		SrcDirArr:   srcDirs,
		DstDir:      dstDir,
		WorkerCount: workCount,
		Version:     version,
	}
}

func (s *Summary) addExtension(name string) {
	if s.Extensions != nil {
		ext := strings.ToLower(filepath.Ext(name))
		if len(ext) > 0 {
			s.Extensions[ext]++
		} else {
			s.Extensions["__OTHERS__"]++
		}
	}
}
