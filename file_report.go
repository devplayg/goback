package goback

import (
    "path/filepath"
    "strings"
)

type FilesReport struct {
    Extension        map[string]*FileStats `json:"extension"`
    SizeDistribution map[int64]int64       `json:"sizeDistribution"`
    SizeRanking      []*FileStats          `json:"sizeRanking"`
    NameRanking      []*FileStats          `json:"nameRanking"`
    nameMap          map[string]*FileStats
}

func NewFilesReport() *FilesReport {
    report := FilesReport{
        Extension:        make(map[string]*FileStats),
        SizeDistribution: NewSizeDistribution(),
        SizeRanking:        make([]*FileStats, 0), // path: size, path
        NameRanking:        make([]*FileStats, 0), // name: count, size, name
    }

    return &report
}

func NewSizeDistribution() map[int64]int64 {
    m := make(map[int64]int64)
    for _, size := range fileSizeCategories {
        m[size] = 0
    }
    return m
}

func (r *FilesReport) addExtension(name string, size int64) {
    ext := strings.ToLower(filepath.Ext(name))
    if _, have := r.Extension[ext]; !have {
        r.Extension[ext] = NewFileStats(ext, size)
        return
    }

    r.Extension[ext].Count++
    r.Extension[ext].Size += size
}

func (r *FilesReport) addSize(size int64) {
    for i := range fileSizeCategories {
        if size <= fileSizeCategories[i] {
            r.SizeDistribution[fileSizeCategories[i]]++
            return
        }
    }
    r.SizeDistribution[size]++
}
