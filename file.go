package goback

import (
	"time"
)

type File struct {
	Size    int64     `json:"size"`
	Hash    string    `json:"hash"`
	ModTime time.Time `json:"mtime"`
	Path    string    `json:"path"`
}

type FileWrapper struct {
	*File
	WhatHappened int     `json:"how"`
	Result       int     `json:"result"`
	Duration     float64 `json:"dur"`
	Message      string  `json:"msg"`
}

func NewFileWrapper(path string, size int64, modTime time.Time) *FileWrapper {
	return &FileWrapper{
		File: &File{
			Size:    size,
			ModTime: modTime,
			Path:    path,
			Hash:    "",
		},
		WhatHappened: 0,
		Result:       0,
		Duration:     0,
		Message:      "",
	}
}
