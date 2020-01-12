package goback

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type File struct {
	Size         int64     `json:"s"`
	Hash         string    `json:"h"`
	ModTime      time.Time `json:"t"`
	Path         string    `json:"p"`
	WhatHappened int       `json:"w"`
	Result       int       `json:"r"`
	Duration     float64   `json:"d"`
	Message      string    `json:"m"`
	data         []byte
}

func (f *File) Marshal() {
	b, err := json.Marshal(f)
	if err != nil {
		log.Error(err)
	}
	f.data = b
}

func newFile(path string, size int64, modTime time.Time) *File {
	return &File{
		Path:    path,
		Size:    size,
		ModTime: modTime,
	}
}

type FileExtended struct {
	*File
	Result   int     `json:"r"`
	Duration float64 `json:"d"`
}

//
//func NewFileExtended(file *File) *FileExtended{
//	return &FileExtended{
//		File:     file,
//		Result:   0,
//		Duration: 0,
//	}
//}

type Directory struct {
	//Files
}

//type FileMap struct {
//	sync.Map
//}
