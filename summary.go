package goback

import (
	"encoding/json"
	"time"
)

type Summary struct {
	Id        int64     `json:"id"`
	Date      time.Time `json:"date"`
	SrcDirArr []string  `json:"srcDirs"`
	DstDir    string    `json:"dstDir"`
	//State     int       `json:"state"`

	TotalSize  int64 `json:"size"`
	TotalCount int64 `json:"count"`

	// Thread-safe
	AddedCount    uint64 `json:"countAdded"`
	ModifiedCount uint64 `json:"countModified"`
	DeletedCount  uint64 `json:"countDeleted"`

	Extensions       map[string]int64 `json:"ext"`
	SizeDistribution map[int64]int64  `json:"sizeDist"`

	BackupSuccessCount uint64 `json:"successCount"`
	BackupFailureCount uint64 `json:"failureCount"`
	//BackupSize    int64
	Message string `json:"msg"`
	Version int    `json:"v"`

	ReadingTime    time.Time `json:"timeReading"` // Step 1
	ComparisonTime time.Time `json:"timeComp"`    // Step 2
	BackupTime     time.Time `json:"timeBak"`     // Step 3
	LoggingTime    time.Time `json:"timeLog"`     // Step 4
	ExecutionTime  float64   `json:"timeExec"`
}

func (s *Summary) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
