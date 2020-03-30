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
	BackupDir   string    `json:"-"`
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

	Message string      `json:"message"`
	Version int         `json:"-"`
	Stats   *Stats      `json:"stats"`
	Keeper  *KeeperDesc `json:"keeper"`

	addedFiles    *sync.Map
	modifiedFiles *sync.Map
	deletedFiles  *sync.Map
	failedFiles   *sync.Map
}

func (s *Summary) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func NewSummary(backupType int, srcDir string, b *Backup) *Summary {
	return &Summary{
		BackupId:      b.Id,
		Date:          time.Now(),
		SrcDir:        srcDir,
		WorkerCount:   b.workerCount,
		Version:       b.version,
		BackupType:    backupType,
		State:         Started,
		Stats:         NewStatsReport(b.minFileSizeToRank),
		addedFiles:    &sync.Map{},
		modifiedFiles: &sync.Map{},
		deletedFiles:  &sync.Map{},
		failedFiles:   &sync.Map{},
		Keeper:        b.keeper.Description(),
	}
}
