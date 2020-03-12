package goback

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type Config struct {
	Server struct {
		Address string `json:"address"`
	} `json:"server"`
	Storages []Storage `json:"storages"`
	Jobs     []Job     `json:"jobs"`
}

type Job struct {
	Id         int      `json:"id" schema:"id"`
	BackupType int      `json:"backup-type"`
	SrcDirs    []string `json:"dirs" schema:"srcDirs"`
	Schedule   string   `json:"schedule" schema:"schedule"`
	Ignore     []string `json:"ignore"`
	StorageId  int      `json:"storage-id"`
}

//
// type Directory string
//
// func (d Directory) String() string {
// 	return string(d)
// }
//
// func (d Directory) IsValid() bool {
// 	return true
// }
//
// func DirsToStrSlice(dirs []Directory) []string {
// 	arr := make([]string, 0)
// 	for _, d := range dirs {
// 		arr = append(arr, d.String())
// 	}
// 	return arr
// }

//
// type Dir struct {
// 	Name    string
// 	IsValid bool
// }

type Storage struct {
	Id       int    `json:"id"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dir      string `json:"dir"`
}

func (c *Config) Save() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("config_.yaml", b, 0644)
}

type AppConfig struct {
	Name        string
	Description string
	Version     string
	Url         string
	Text1       string
	Text2       string
	Year        int
	Company     string
	Debug       bool
	Trace       bool
	LogImgPath  string
}
