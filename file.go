package goback

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "path/filepath"
    "time"
)

type File struct {
    Size    int64     `json:"size"`
    ModTime time.Time `json:"mtime"`
    Path    string    `json:"path"`
    Hash    string    `json:"hash"`
}

type FileWrapper struct {
    *File
    WhatHappened int     `json:"how"`
    Result       int     `json:"result"`
    Duration     float64 `json:"dur"`
    Message      string  `json:"msg"`
}

type ExtensionStats struct {
    Ext   string `json:"ext"`
    Size  int64  `json:"size"`
    Count int64  `json:"count"`
}

func NewFileStats(ext string, size int64) *ExtensionStats {
    return &ExtensionStats{
        Ext:   ext,
        Size:  size,
        Count: 1,
    }
}

type FileNameStats struct {
    Name  string   `json:"name"`
    Size  int64    `json:"size"`
    Paths []string `json:"paths"`
    Count int64    `json:"count"`
}

func NewFileNameStats(file *File) *FileNameStats {
    stats := FileNameStats{
        Name:  filepath.Base(file.Path),
        Size:  file.Size,
        Paths: make([]string, 0),
        Count: 1,
    }
    stats.Paths = append(stats.Paths, file.Path)
    return &stats
}

func GetFileNameKey(file *File) string {
    return fmt.Sprintf("%s-%d", filepath.Base(file.Path), file.Size)
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

func WrapFile(file *File) *FileWrapper {
    return &FileWrapper{
        File:         file,
        WhatHappened: 0,
        Result:       0,
        Duration:     0,
        Message:      "",
    }
}

type DirInfo struct {
    Checksum string
    dbPath   string
}

func NewDirInfo(srcDir, dstDir string) *DirInfo {
    checksum := md5.Sum([]byte(srcDir))
    checksumStr := hex.EncodeToString(checksum[:])
    return &DirInfo{
        Checksum: checksumStr,
        dbPath:   filepath.Join(dstDir, fmt.Sprintf(FilesDbName, checksumStr)),
    }
}
