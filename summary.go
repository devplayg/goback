package goback

import (
    "encoding/json"
    "path/filepath"
    "strings"
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
    TotalSize   uint64    `json:"totalSize"`
    TotalCount  int64     `json:"totalCount"`

    addedFiles    *sync.Map
    modifiedFiles *sync.Map
    deletedFiles  *sync.Map
    failedFiles   *sync.Map

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

    ExtensionMap     map[string]int64 `json:"extensions"`
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

func NewSummary(summaryId, backupId int, srcDir, dstDir string, backupType, workCount, version int) *Summary {
    return &Summary{
        Id:               summaryId,
        BackupId:         backupId,
        Date:             time.Now(),
        SrcDir:           srcDir,
        DstDir:           dstDir,
        WorkerCount:      workCount,
        Version:          version,
        BackupType:       backupType,
        State:            Started,
        ExtensionMap:     make(map[string]int64),
        SizeDistribution: make(map[int64]int64),

        addedFiles:    &sync.Map{},
        modifiedFiles: &sync.Map{},
        deletedFiles:  &sync.Map{},
        failedFiles:   &sync.Map{},
    }
}

func (s *Summary) addExtension(name string) {
    if s.ExtensionMap != nil {
        ext := strings.ToLower(filepath.Ext(name))
        if len(ext) > 0 {
            s.ExtensionMap[ext]++
        } else {
            s.ExtensionMap["__NO_EXT__"]++
        }
    }
}
