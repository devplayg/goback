package goback

import (
	"encoding/json"
	"time"
)

type Summary struct {
	Id         int64
	Date       time.Time
	SrcDir     string
	DstDir     string
	State      int
	TotalSize  int64
	TotalCount int

	BackupAdded    int32
	BackupModified int32
	BackupDeleted  int32

	BackupSuccess uint32
	BackupFailure uint32

	BackupSize uint64
	Message    string

	ReadingTime    time.Time
	ComparisonTime time.Time
	LoggingTime    time.Time
	ExecutionTime  float64
}

func (s *Summary) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
