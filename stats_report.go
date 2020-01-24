package goback

import (
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type StatsReport struct {
	ExtensionRanking []*ExtensionStats `json:"extRanking"`
	SizeDistribution map[int64]int64   `json:"sizeDistribution"`
	SizeRanking      []*File           `json:"sizeRanking"`
	NameRanking      []*FileNameStats  `json:"nameRanking"`

	extension       map[string]*ExtensionStats
	nameMap         map[string]*FileNameStats
	sizeRankMinSize int64
}

func NewStatsReport(sizeRankMinSize int64) *StatsReport {
	return &StatsReport{
		ExtensionRanking: make([]*ExtensionStats, 0),
		extension:        make(map[string]*ExtensionStats),
		SizeDistribution: NewSizeDistribution(),
		SizeRanking:      make([]*File, 0),          // path: size, path
		NameRanking:      make([]*FileNameStats, 0), // name: count, size, name
		nameMap:          make(map[string]*FileNameStats),
		sizeRankMinSize:  sizeRankMinSize,
	}
}

func (r *StatsReport) addToReport(file *FileReport) {
	r.addExtension(file.path, file.Size)
	r.addSize(file)
	r.addName(file)
}

func (r *StatsReport) addExtension(name string, size int64) {
	ext := strings.ToLower(filepath.Ext(name))
	if _, have := r.extension[ext]; !have {
		r.extension[ext] = NewFileStats(ext, size)
		return
	}

	r.extension[ext].Count++
	r.extension[ext].Size += size
}

func (r *StatsReport) addName(file *File) {
	name := GetFileNameKey(file)
	if _, have := r.nameMap[name]; !have {
		r.nameMap[name] = NewFileNameStats(file)
		return
	}
	r.nameMap[name].Paths = append(r.nameMap[name].Paths, file.Path)
	r.nameMap[name].Count++
}

func (r *StatsReport) addSize(file *File) {
	// Size ranking
	if file.Size >= r.sizeRankMinSize {
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

func (r *StatsReport) tune(rank int) {
	// Size ranking
	sort.Slice(r.SizeRanking, func(i, j int) bool {
		return r.SizeRanking[i].Size > r.SizeRanking[j].Size
	})
	if len(r.SizeRanking) > rank {
		r.SizeRanking = r.SizeRanking[0:rank]
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
	}

	// Extension (keep all)
	for _, stats := range r.extension {
		r.ExtensionRanking = append(r.ExtensionRanking, stats)
	}
	sort.Slice(r.ExtensionRanking, func(i, j int) bool {
		return r.ExtensionRanking[i].Size > r.ExtensionRanking[j].Size
	})
	//if len(r.ExtensionRanking) > rank {
	//    r.ExtensionRanking = r.ExtensionRanking[0:rank]
	//    return
	//}
}

type StatsReportWithList struct {
	Files  []*FileReport `json:"files"`
	Report *StatsReport  `json:"report"`
}

func CreateFilesReportWithList(files []*FileReport, sizeRankMinSize int64, rank int) *StatsReportWithList {
	r := StatsReportWithList{
		Files:  files,
		Report: NewStatsReport(sizeRankMinSize),
	}
	for _, f := range files {
		r.Report.addToReport(f)
	}
	r.Report.tune(rank)
	return &r
}
