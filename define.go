package goback

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"hash"
	"path/filepath"
)

const (
	FileModified = 1
	FileAdded    = 2
	FileDeleted  = 4

	FileBackupFailed    = -1
	FileBackupSucceeded = 1

	//DbDir = "db"
	FilesDbName       = "files-%s.db"
	SummaryDbName     = "summary.db"
	SummaryTempDbName = "summary.db.lock"
	ChangesDbName     = "changes-%s-%d.db"
	ConfigFileName    = "config.yaml"

	Initial     = 1
	Incremental = 2
	Full        = 4

	LocalDisk = 1
	Ftp       = 2
	Sftp      = 4

	// Failed = -1

	GobEncoding  = 1
	JsonEncoding = 2
)

var (
	SummaryBucketName []byte = []byte("summary")
	BackupBucketName  []byte = []byte("backup")
	ConfigBucketName  []byte = []byte("config")
)

const (
	Started = iota + 1
	Read
	Compared
	Copied
	Logged
)

const (
	kB = 1000
	MB = 1000000
	GB = 1000000000
	TB = 1000000000000
)

var fileSizeCategories = []int64{
	0,

	5 * kB,
	50 * kB,
	500 * kB,

	5 * MB,
	50 * MB,
	500 * MB,

	5 * GB,
	50 * GB,

	5 * TB,
}

var log *logrus.Logger

type Fields *logrus.Fields

var (
	HashKey     = sha256.Sum256([]byte("goback"))
	Highwayhash hash.Hash
)

type dirInfo struct {
	checksum string
	dbPath   string
}

func newDirInfo(srcDir, dbDir string) *dirInfo {
	b := md5.Sum([]byte(srcDir))
	checksum := hex.EncodeToString(b[:])
	return &dirInfo{
		checksum: checksum,
		dbPath:   filepath.Join(dbDir, fmt.Sprintf(FilesDbName, checksum)),
	}
}
