package goback

import (
	"crypto/sha256"
	"hash"
)

const (
	FileModified = 1 << iota // 1
	FileAdded    = 1 << iota // 2
	FileDeleted  = 1 << iota // 4

	FileBackupFailed    = -1
	FileBackupSucceeded = 1

	FilesDbName   = "files-%s.db"
	SummaryDbName = "summary.db"

	InitialBackup = 1
	NormalBackup  = 2

	Failed = -1

	GobEncoding = 1
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
)

var fileSizeCategories = []int64{
	1 * kB,
	5 * kB,
	10 * kB,
	50 * kB,
	100 * kB,
	500 * kB,

	1 * MB,
	5 * MB,
	10 * MB,
	50 * MB,
	100 * MB,
	500 * MB,

	1 * GB,
	5 * GB,
	10 * GB,
	50 * GB,
	100 * GB,
	500 * GB,
}

var (
	HashKey     = sha256.Sum256([]byte("goback"))
	Highwayhash hash.Hash
)
