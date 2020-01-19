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

	FileMapDbName = "files.db"
	SummaryDbName = "summary.db"
)

var (
	HashKey     = sha256.Sum256([]byte("goback"))
	Highwayhash hash.Hash
)
