package goback

import (
	"encoding/json"
	"time"
)

type Summary struct {
	Id        int64     `json:"id"`
	Date      time.Time `json:"d"`
	SrcDirArr []string  `json:"src"`
	DstDir    string    `json:"dst"`
	State     int       `json:"s"`

	TotalSize  int64 `json:"sz"`
	TotalCount int64 `json:"c"`

	// Thread-safe
	AddedCount    uint64 `json:"c_a"`
	ModifiedCount uint64 `json:"c_m"`
	DeletedCount  uint64 `json:"c_d"`

	Extensions       map[string]int64 `json:"ext"`
	SizeDistribution map[int64]int64  `json:"szd"`

	BackupSuccessCount int64 `json:"ok"`
	BackupFailureCount int64 `json:"fail"`
	//BackupSize    int64
	Message string `json:"msg"`
	Version int    `json:"v"`

	ReadingTime    time.Time `json:"t_r"` // Step 1
	ComparisonTime time.Time `json:"t_c"` // Step 2
	BackupTime     time.Time `json:"t_b"` // Step 3
	LoggingTime    time.Time `json:"t_l"` // Step 4
	ExecutionTime  float64   `json:"t_e"`
}

func (s *Summary) Marshal() ([]byte, error) {
	return json.Marshal(s)
}
