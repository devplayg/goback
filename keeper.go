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
	Protocol int    `json:"protocol"` // Local / FTP / SFTP
	Host     string `json:"host"`     // local or remote host
	Dir      string `json:"dir"`
}

func NewKeeper(job *Job) Keeper {
	if job.Storage.Protocol == LocalDisk {
		return NewLocalKeeper(job.Storage.Dir)
	}

	if job.Storage.Protocol == Sftp {
		return NewSftpKeeper(job.Storage)
	}

	return nil
}
