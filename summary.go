package goback

import (
	"encoding/json"
	"time"
)

type Summary struct {
	Id         int64
	Date       time.Time
	SrcDirArr  []string
	DstDir     string
	State      int
	TotalSize  int64
	TotalCount int

	AddedFiles    int // target to backup
	ModifiedFiles int // target to backup
	DeletedFiles  int

	BackupSuccess int
	BackupFailure int
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
