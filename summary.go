package goback

import "time"

type Summary struct {
	Id         int64
	Date       time.Time
	SrcDir     string
	DstDir     string
	State      int
	TotalSize  uint64
	TotalCount uint32

	BackupAdded    uint32
	BackupModified uint32
	BackupDeleted  uint32

	BackupSuccess uint32
	BackupFailure uint32

	BackupSize uint64
	Message    string

	ReadingTime    time.Time
	ComparisonTime time.Time
	LoggingTime    time.Time
	ExecutionTime  float64
}

func newSummary(lastId int64, srcDir string) *Summary {
	return &Summary{
		Id:     lastId,
		Date:   time.Now(),
		SrcDir: srcDir,
		State:  BackupReady,
	}
}
