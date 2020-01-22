package goback

import (
    "math"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

type FilesReport struct {
    ExtensionRanking []*FileStats          `json:"extRanking"`
    SizeDistribution map[int64]int64       `json:"sizeDistribution"`
    SizeRanking      []*File               `json:"sizeRanking"`
    NameRanking      []*FileNameStats      `json:"nameRanking"`

    extension        map[string]*FileStats
    nameMap          map[string]*FileNameStats
}

func NewFilesReport() *FilesReport {
    return &FilesReport{
        ExtensionRanking: make([]*FileStats, 0),
        extension:        make(map[string]*FileStats),
        SizeDistribution: NewSizeDistribution(),
        SizeRanking:      make([]*File, 0),          // path: size, path
        NameRanking:      make([]*FileNameStats, 0), // name: count, size, name
        nameMap:          make(map[string]*FileNameStats),
    }
}

func (r *FilesReport) addToReport(file *File) {
    r.addExtension(file.Path, file.Size)
    r.addSize(file)
    r.addName(file)
}

func (r *FilesReport) addExtension(name string, size int64) {
    ext := strings.ToLower(filepath.Ext(name))
    if _, have := r.extension[ext]; !have {
        r.extension[ext] = NewFileStats(ext, size)
        return
    }

    r.extension[ext].Count++
    r.extension[ext].Size += size
}

func (r *FilesReport) addName(file *File) {
    name := GetFileNameKey(file)
    if _, have := r.nameMap[name]; !have {
        r.nameMap[name] = NewFileNameStats(file)
        return
    }
    r.nameMap[name].Paths = append(r.nameMap[name].Paths, file.Path)
    r.nameMap[name].Count++
}

func (r *FilesReport) addSize(file *File) {
    // file size is larger thant 10 MB
    if file.Size > 10*MB {
        r.SizeRanking = append(r.SizeRanking, file)
    }

    // Size-zero file
    if file.Size == 0 {
        r.SizeDistribution[0]++
        return
    }

    // Small file
    if file.Size < fileSizeCategories[0] {
        r.SizeDistribution[fileSizeCategories[0]]++
        return
    }

    // Big file
    if file.Size > fileSizeCategories[len(fileSizeCategories)-1] {
        r.SizeDistribution[file.Size]++
        return
    }

    n := math.Pow10(len(strconv.FormatInt(file.Size, 10)))
    if n/2 > float64(file.Size) {
        r.SizeDistribution[int64(n/2)]++
        return
    }

    r.SizeDistribution[int64(n)]++
}

func (r *FilesReport) tune(rank int) {
    // Size ranking
    sort.Slice(r.SizeRanking, func(i, j int) bool {
        return r.SizeRanking[i].Size > r.SizeRanking[j].Size
    })

    if len(r.SizeRanking) > rank {
        r.SizeRanking = r.SizeRanking[0:rank]
        return
    }

    // Name ranking
    for _, stats := range r.nameMap {
        if stats.Count <= 1 {
            continue
        }
        r.NameRanking = append(r.NameRanking, stats)
    }
    sort.Slice(r.NameRanking, func(i, j int) bool {
        return (r.NameRanking[i].Size * r.NameRanking[i].Count) > (r.NameRanking[j].Size * r.NameRanking[j].Count)
    })
    if len(r.NameRanking) > rank {
        r.NameRanking = r.NameRanking[0:rank]
        return
    }

    // Extension
    for _, stats := range r.extension {
        r.ExtensionRanking = append(r.ExtensionRanking, stats)
    }
    sort.Slice(r.ExtensionRanking, func(i, j int) bool {
        return r.ExtensionRanking[i].Size > r.ExtensionRanking[j].Size
    })
    if len(r.ExtensionRanking) > rank {
        r.ExtensionRanking = r.ExtensionRanking[0:rank]
        return
    }
}
