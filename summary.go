package goback

import (
	"encoding/json"
	"sync"
	"time"
)

type Summary struct {
	Id        int64
	Date      time.Time
	SrcDirArr []string
	DstDir    string
	State     int

	TotalSize  int64
	TotalCount int64

	AddedCount    uint64
	ModifiedCount uint64
	DeletedCount  uint64

	addedFiles    *sync.Map
	modifiedFiles *sync.Map
	deletedFiles  *sync.Map

	Extensions       map[string]int64
	SizeDistribution map[int64]int64

	BackupSuccess int64
	BackupFailure int64
	BackupSize    int64
	Message       string

	ReadingTime    time.Time // Step 1
	ComparisonTime time.Time // Step 2
	BackupTime     time.Time // Step 3
	LoggingTime    time.Time // Step 4
	ExecutionTime  float64
}

func (s *Summary) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
