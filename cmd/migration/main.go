package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/devplayg/goback"
	"github.com/devplayg/goutils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const dbDir = "db"

var db *bolt.DB

func init() {
	if err := os.Mkdir(dbDir, 0755); err != nil {
		panic(err)
	}
	boltDb, err := bolt.Open(filepath.Join(dbDir, goback.ProcessName+".db"), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	db = boltDb
	if err := db.Batch(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(goback.SummaryBucket); err != nil {
			return fmt.Errorf("failed to create summary bucket; %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(goback.BackupBucket); err != nil {
			return fmt.Errorf("failed to create backup group bucket; %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(goback.ConfigBucket); err != nil {
			return fmt.Errorf("failed to create backup group bucket; %w", err)
		}
		return nil
	}); err != nil {
		return
	}

}
func main() {

	// Check arguments
	if len(os.Args) < 2 {
		return
	}
	backupDir := os.Args[1]
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		fmt.Printf("backup directory (%s) does not exist\n", backupDir)
		return
	}

	fmt.Printf("Loading old database...")
	oldSummaries, err := readSummaryDb(filepath.Join(backupDir, "summary.db"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("done. Found %d backup logs\n", len(oldSummaries))

	var newSummaries []*goback.Summary
	maxBackupId := 0
	for _, s := range oldSummaries {
		newSummaries = append(newSummaries, s.toNew())
		if s.BackupId > maxBackupId {
			maxBackupId = s.BackupId
		}

		// Copy changes
		h := md5.Sum([]byte(s.SrcDir))
		key := hex.EncodeToString(h[:])
		srcFile := filepath.Join(backupDir, filepath.Base(s.BackupDir), fmt.Sprintf("changes-%s.db", key))
		fmt.Printf("[%d/%d] %v / src=%s\n", s.BackupId, s.Id, s.Date.Format("2006-01-02 15:04:05"), s.BackupDir)
		dstFile := filepath.Join(dbDir, fmt.Sprintf("changes-%d-%s.db", s.BackupId, key))
		if _, err := copy(srcFile, dstFile); err != nil {
			if s.BackupType != goback.Initial {
				fmt.Printf("[error] %s\n", err.Error())
			}
			continue
		}
	}

	for i := 1; i <= maxBackupId; i++ {
		if err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(goback.BackupBucket)
			newId, _ := b.NextSequence()
			id := int(newId)

			return b.Put(goback.IntToByte(id), nil)
		}); err != nil {
			fmt.Printf("[error] %s\n", err.Error())
		}
	}

	if err := save(dbDir, newSummaries); err != nil {
		panic(err)
	}
}

func save(dbDir string, summary []*goback.Summary) error {
	db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(goback.SummaryBucket)
		for i := range summary {
			data, err := summary[i].Marshal()
			if err != nil {
				fmt.Printf("[error] %s\n", err.Error())
				continue
			}
			//fmt.Printf("%d\n", summary[i].Id)
			seq, _ := b.NextSequence()
			if seq != uint64(summary[i].Id) {
				panic("what!!! sum id is diff")
			}
			if err := b.Put(goback.IntToByte(summary[i].Id), data); err != nil {
				fmt.Printf("[error] %s\n", err.Error())
				continue
			}
		}
		return nil
	})

	return nil
}

func readSummaryDb(path string) ([]*OldSummary, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decompressed, err := goutils.Gunzip(data)
	if err != nil {
		return nil, err
	}

	var summaries []*OldSummary
	if err := goutils.GobDecode(decompressed, &summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

type OldSummary struct {
	Id          int       `json:"id"`
	BackupId    int       `json:"backupId"`
	Date        time.Time `json:"date"`
	SrcDir      string    `json:"srcDir"`
	DstDir      string    `json:"-"`
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

	Message string        `json:"message"`
	Version int           `json:"-"`
	Stats   *goback.Stats `json:"stats"`

	addedFiles    *sync.Map
	modifiedFiles *sync.Map
	deletedFiles  *sync.Map
	failedFiles   *sync.Map
}

func (s *OldSummary) toNew() *goback.Summary {
	return &goback.Summary{
		Id:             s.Id,
		BackupId:       s.BackupId,
		Date:           s.Date,
		SrcDir:         s.SrcDir,
		State:          s.State,
		BackupDir:      s.BackupDir,
		BackupType:     s.BackupType,
		WorkerCount:    s.WorkerCount,
		TotalSize:      s.TotalSize,
		TotalCount:     s.TotalCount,
		AddedCount:     s.AddedCount,
		AddedSize:      s.AddedSize,
		ModifiedCount:  s.ModifiedCount,
		ModifiedSize:   s.ModifiedSize,
		DeletedCount:   s.DeletedCount,
		DeletedSize:    s.DeletedSize,
		FailedCount:    s.FailedCount,
		FailedSize:     s.FailedSize,
		SuccessCount:   s.SuccessCount,
		SuccessSize:    s.SuccessSize,
		ReadingTime:    s.ReadingTime,
		ComparisonTime: s.ComparisonTime,
		BackupTime:     s.BackupTime,
		LoggingTime:    s.LoggingTime,
		ExecutionTime:  s.ExecutionTime,
		Message:        s.Message,
		Version:        3,
		Stats:          s.Stats,
	}
}

// --------------------

type Stats struct {
	ExtRanking  []*goback.ExtStats      `json:"extRanking"`
	NameRanking []*goback.NameStats     `json:"nameRanking"`
	SizeRanking []*goback.FileGrid      `json:"sizeRanking"`
	SizeDist    []*goback.SizeDistStats `json:"sizeDist"`
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
