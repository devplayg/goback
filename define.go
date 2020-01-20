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
)
const (
	Started = iota + 1
	Read
	Compared
	Copied
	Logged
	Completed
)

var (
	HashKey     = sha256.Sum256([]byte("goback"))
	Highwayhash hash.Hash
)
