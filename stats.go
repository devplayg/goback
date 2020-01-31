package goback

import (
	"sort"
)

type ExtStats struct {
	Ext   string `json:"ext"`
	Size  int64  `json:"size"`
	Count int64  `json:"count"`
}

func NewExtensionStats(ext string, size, count int64) *ExtStats {
	return &ExtStats{
		Ext:   ext,
		Size:  size,
		Count: count,
	}
}

type NameStats struct {
	Name  string   `json:"name"`
	Size  int64    `json:"size"`
	Paths []string `json:"paths"`
	Count int64    `json:"count"`
}

func NewNameStats(file *FileGrid) *NameStats {
	stats := NameStats{
		Name:  file.Name,
		Size:  file.Size,
		Paths: make([]string, 0),
		Count: 1,
	}
	stats.Paths = append(stats.Paths, file.Dir)
	return &stats
}

type SizeDistStats struct {
	SizeDist int64 `json:"sizeDist"`
	Size     int64 `json:"size"`
	Count    int64 `json:"count"`
}

func NewSizeDistStats(sizeDist, size, count int64) *SizeDistStats {
	return &SizeDistStats{
		SizeDist: sizeDist,
		Size:     size,
		Count:    count,
	}
}

type Stats struct {
	ExtRanking  []*ExtStats      `json:"extRanking"`
	NameRanking []*NameStats     `json:"nameRanking"`
	SizeRanking []*FileGrid      `json:"sizeRanking"`
	SizeDist    []*SizeDistStats `json:"sizeDist"`

	extMap          map[string]*ExtStats
	nameMap         map[string]*NameStats
	sizeDistMap     map[int64]*SizeDistStats
	sizeRankMinSize int64
}

func NewStatsReport(sizeRankMinSize int64) *Stats {
	return &Stats{
		ExtRanking:  make([]*ExtStats, 0),
		SizeRanking: make([]*FileGrid, 0),  // path: size, path
		NameRanking: make([]*NameStats, 0), // name: count, size, name

		extMap:      make(map[string]*ExtStats),
		sizeDistMap: NewSizeDistMap(),
		nameMap:     make(map[string]*NameStats),

		sizeRankMinSize: sizeRankMinSize,
	}

}

func (s *Stats) addToStats(file *FileGrid) {
	s.addToExtStats(file)
	s.addToSizeStats(file)
	s.addToNameStats(file)
}

func (s *Stats) addToExtStats(file *FileGrid) {
	if _, have := s.extMap[file.Ext]; !have {
		s.extMap[file.Ext] = NewExtensionStats(file.Ext, file.Size, 1)
		return
	}
	s.extMap[file.Ext].Count++
	s.extMap[file.Ext].Size += file.Size
}

func (s *Stats) addToNameStats(file *FileGrid) {
	name := GetFileNameKey(file.Name, file.Size)
	if _, have := s.nameMap[name]; !have {
		s.nameMap[name] = NewNameStats(file)
		return
	}
	s.nameMap[name].Paths = append(s.nameMap[name].Paths, file.Dir)
	s.nameMap[name].Count++
}

func (s *Stats) addToSizeStats(file *FileGrid) {
	// Size ranking
	if file.Size >= s.sizeRankMinSize {
		s.SizeRanking = append(s.SizeRanking, file)
	}

	for i := 1; i < len(fileSizeCategories); i++ {
		if file.Size <= fileSizeCategories[i] {
			s.sizeDistMap[fileSizeCategories[i]].Count++
			s.sizeDistMap[fileSizeCategories[i]].Size += file.Size
			return
		}
	}
}

func (s *Stats) rank(rank int) {
	// Size ranking
	sort.Slice(s.SizeRanking, func(i, j int) bool {
		return s.SizeRanking[i].Size > s.SizeRanking[j].Size
	})
	if len(s.SizeRanking) > rank {
		s.SizeRanking = s.SizeRanking[0:rank]
	}

	// Name ranking
	for _, stats := range s.nameMap {
		if stats.Count <= 1 {
			continue
		}
		s.NameRanking = append(s.NameRanking, stats)
	}
	sort.Slice(s.NameRanking, func(i, j int) bool {
		return (s.NameRanking[i].Size * s.NameRanking[i].Count) > (s.NameRanking[j].Size * s.NameRanking[j].Count)
	})
	if len(s.NameRanking) > rank {
		s.NameRanking = s.NameRanking[0:rank]
	}

	// Extension (keep all)
	for _, stats := range s.extMap {
		s.ExtRanking = append(s.ExtRanking, stats)
	}
	sort.Slice(s.ExtRanking, func(i, j int) bool {
		return s.ExtRanking[i].Size > s.ExtRanking[j].Size
	})
	// if len(s.ExtensionRanking) > rank {
	//    s.ExtensionRanking = s.ExtensionRanking[0:rank]
	//    return
	// }

	for _, stats := range s.sizeDistMap {
		s.SizeDist = append(s.SizeDist, stats)
	}
	sort.Slice(s.SizeDist, func(i, j int) bool {
		return s.SizeDist[i].SizeDist > s.SizeDist[j].SizeDist
	})
}

type StatsReportWithList struct {
	Files  []*FileGrid `json:"files"`
	Size   uint64      `json:"size"`
	Report *Stats      `json:"report"`
}

func CreateFilesReportWithList(files []*FileGrid, size uint64, minSizeToRank int64, rank int) *StatsReportWithList {
	r := StatsReportWithList{
		Files:  files,
		Size:   size,
		Report: NewStatsReport(minSizeToRank),
	}
	for _, f := range files {
		r.Report.addToStats(f)
	}
	r.Report.rank(rank)
	return &r
}
