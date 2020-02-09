package goback

import "time"

type Keeper interface {
	Init(t time.Time) error
	keep(string) (string, float64, error)
	Close() error
	Description() *KeeperDesc
	Chtimes(name string, atime time.Time, mtime time.Time) error
	Active() bool
}

type KeeperDesc struct {
	Protocol int    // Local / FTP / SFTP
	Host     string // local or remote host
	Dir      string
}
