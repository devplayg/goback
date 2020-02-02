package goback

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"path/filepath"
)

const (
	FileModified = 1
	FileAdded    = 2
	FileDeleted  = 4

	FileBackupFailed    = -1
	FileBackupSucceeded = 1

	FilesDbName   = "files-%s.db"
	SummaryDbName = "summary.db"

	Initial     = 1
	Incremental = 2

	LocalDisk = 1
	Ftp       = 3
	Sftp      = 4

	// Failed = -1

	GobEncoding  = 1
	JsonEncoding = 2
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

var (
	HashKey     = sha256.Sum256([]byte("goback"))
	Highwayhash hash.Hash
)

type dirInfo struct {
	checksum string
	dbPath   string
}

func newDirInfo(srcDir, dstDir string) *dirInfo {
	b := md5.Sum([]byte(srcDir))
	checksum := hex.EncodeToString(b[:])
	return &dirInfo{
		checksum: checksum,
		dbPath:   filepath.Join(dstDir, fmt.Sprintf(FilesDbName, checksum)),
	}
}
