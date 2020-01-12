package goback

import (
	"crypto/sha256"
	"hash"
)

const (
	FileModified = 1 << iota // 1
	FileAdded    = 1 << iota // 2
	FileDeleted  = 1 << iota // 4

	BackupReady    = 1
	BackupRunning  = 2
	BackupFinished = 3

	Failure = -1
	Success = 1

	BackupPrefixStr = "backup-"
)

var (
	BucketSummary = []byte("summary")
	BucketFiles   = []byte("files")

	BucketAdded    = []byte("added")
	BucketModified = []byte("modified")
	BucketDeleted  = []byte("deleted")
	BucketFailed   = []byte("failed")

	HashKey     = sha256.Sum256([]byte("goback"))
	Highwayhash hash.Hash
)
