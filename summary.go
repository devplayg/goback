package goback

import (
	"encoding/json"
	"sync"
	"time"
)

type Summary struct {
	Id          int       `json:"id"`
	BackupId    int       `json:"backupId"`
	Date        time.Time `json:"date"`
	SrcDir      string    `json:"srcDir"`
	DstDir      string    `json:"dstDir"`
	BackupType  int       `json:"backupType"`
	State       int       `json:"state"`
	WorkerCount int       `json:"workerCount"`

	// Thread-safe
	TotalSize     uint64 `json:"totalSize"`
	TotalCount    int64  `json:"totalCount"`
	AddedCount    uint64 `json:"countAdded"`
	AddedSize     uint64 `json:"sizeAdded"`
	ModifiedCount uint64 `json:"countModified"`
	ModifiedSize  uint64 `json:"sizeModified"`
	DeletedCount  uint64 `json:"countDeleted"`
	DeletedSize   uint64 `json:"sizeDeleted"`

	// Backup
	FailedCount  uint64 `json:"countFailed"`
	FailedSize   uint64 `json:"sizeFailed"`
	SuccessCount uint64 `json:"countSuccess"`
	SuccessSize  uint64 `json:"sizeSuccess"`

	ReadingTime    time.Time `json:"readingTime"`    // Step 1
	ComparisonTime time.Time `json:"comparisonTime"` // Step 2
	BackupTime     time.Time `json:"backupTime"`     // Step 3
	LoggingTime    time.Time `json:"loggingTime"`    // Step 4
	ExecutionTime  float64   `json:"execTime"`       // Result

	Message string `json:"message"`
	Version int    `json:"version"`
	Stats   *Stats `json:"stats"`

	addedFiles    *sync.Map
	modifiedFiles *sync.Map
	deletedFiles  *sync.Map
	failedFiles   *sync.Map
}

func (s *Summary) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func NewSummary(summaryId, backupId int, srcDir, dstDir string, backupType, workCount, version int, sizeRankMinSize int64) *Summary {
	return &Summary{
		Id:            summaryId,
		BackupId:      backupId,
		Date:          time.Now(),
		SrcDir:        srcDir,
		DstDir:        dstDir,
		WorkerCount:   workCount,
		Version:       version,
		BackupType:    backupType,
		State:         Started,
		Stats:         NewStatsReport(sizeRankMinSize),
		addedFiles:    &sync.Map{},
		modifiedFiles: &sync.Map{},
		deletedFiles:  &sync.Map{},
		failedFiles:   &sync.Map{},
	}
}
