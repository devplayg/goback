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

type FilesTopN struct {
    *File
    Count int64
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
    d := DirInfo{
        Checksum: checksumStr,
        dbPath:   filepath.Join(dstDir, fmt.Sprintf(FilesDbName, checksumStr)),
    }
    return &d
}

type FileGroupReport struct {
    ExtensionMap     map[string]int64
    SizeDistribution map[int64]int64
    TopN             struct {
        BySize      []*FilesTopN
        ByName      []*FilesTopN
        ByExtension []*FilesTopN
    }
}
