package goback

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
)

type Config struct {
	Storages []*Storage `json:"storages"`
	Jobs     []*Job     `json:"jobs"`
}

func NewConfig() *Config {
	return &Config{
		Storages: nil,
		Jobs:     nil,
	}
}

type Job struct {
	Id         int      `json:"id" schema:"id"`
	BackupType int      `json:"backup-type"`
	SrcDirs    []string `json:"dirs" schema:"srcDirs"`
	Schedule   string   `json:"schedule" schema:"schedule"`
	Ignore     []string `json:"ignore"`
	StorageId  int      `json:"storage-id"`
	Enabled    bool     `json:"enabled"`
	Storage    *Storage `json:"-"`
	running    bool

	Checksum    string `json:"-"`
	cronEntryId *cron.EntryID
}

func (j *Job) IsValid() error {
	if len(j.SrcDirs) < 1 {
		return fmt.Errorf("no directories")
	}
	return nil
}

func (j *Job) Tune() {
	j.SrcDirs = uniqueStrings(j.SrcDirs)
	j.Schedule = strings.TrimSpace(j.Schedule)
	if j.Storage != nil {
		j.Storage.Tune()
	}
}

type Storage struct {
	Id       int    `json:"id"`
	Protocol int    `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dir      string `json:"dir"`

	Checksum string `json:"-"`
}

func (s *Storage) Tune() {
	if s == nil {
		return
	}
	s.Host = strings.TrimSpace(s.Host)
	s.Username = strings.TrimSpace(s.Username)
	s.Password = strings.TrimSpace(s.Password)
	s.Dir = strings.TrimSpace(s.Dir)
}

type AppConfig struct {
	Address     string
	Name        string
	Description string
	Version     string
	Url         string
	Text1       string
	Text2       string
	Year        int
	Company     string
	Debug       bool
	Trace       bool
	Verbose     bool
	// LogImgPath  string
}
