package goback

import (
	"path/filepath"
	"strings"
	"time"
)

// Default file structure
type File struct {
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mtime"`

	name string
	dir  string
	ext  string
}

func (f *File) fill() {
	dir, name := filepath.Split(f.Path)
	f.dir = dir
	f.name = name
	f.ext = strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
}

func (f *File) WrapInFileWrapper(fill bool) *FileWrapper {
	if fill {
		f.fill()
	}
	return &FileWrapper{
		File:         f,
		WhatHappened: 0,
		State:        0,
		ExecTime:     make([]float64, 0),
		Message:      make([]string, 0),
	}
}

func NewFile(path string, size int64, modTime time.Time) *File {
	file := File{
		Path:    path,
		Size:    size,
		ModTime: modTime,
	}
	file.fill()
	return &file
}

// File wrapper structure
type FileWrapper struct {
	*File
	WhatHappened int       `json:"how"`
	State        int       `json:"state"`
	ExecTime     []float64 `json:"execTimes"`
	Message      []string  `json:"msg"`
}

func NewFileWrapper(path string, size int64, modTime time.Time) *FileWrapper {
	return &FileWrapper{
		File:         NewFile(path, size, modTime),
		WhatHappened: 0,
		State:        0,
		ExecTime:     make([]float64, 0),
		Message:      make([]string, 0),
	}
}

func (f *FileWrapper) WrapInFileGrid() *FileGrid {
	return &FileGrid{
		Dir:      f.dir,
		Name:     f.name,
		Ext:      f.ext,
		Size:     f.Size,
		State:    f.State,
		ModTime:  f.ModTime,
		ExecTime: f.ExecTime,
		Message:  f.Message,
	}
}

type FileGrid struct {
	Name     string    `json:"name"`
	Dir      string    `json:"dir"`
	Ext      string    `json:"ext"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"mtime"`
	State    int       `json:"state"`
	ExecTime []float64 `json:"execTimes"`
	Message  []string  `json:"msg"`
}
