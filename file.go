package goback

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type File struct {
	Size int64 `json:"s"`
	//Result  int    `json:"r"`
	//Message string `json:"m"`
	Hash    string    `json:"h"`
	ModTime time.Time `json:"t"`
	path    string
	data    []byte
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
		path: path,
		Size: size,
		//modTime: modTime,
		//ModTime: modTime.Format(time.RFC3339),
		ModTime: modTime,
	}
}
