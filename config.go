package goback

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type SystemConfig struct {
	App    AppConfig
	Backup []BackupDir `json:"backup"`
}

type BackupDir struct {
	Dir      string `json:"dir"`
	Schedule struct {
		Full        []string `json:"full"`
		Incremental []string `json:"incremental"`
	}
	Storages []BackupStorage `json:"storages"`
	Ignore   []string        `json:"ignore"`
}

type BackupStorage struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dir      string `json:"dir"`
}

func (c *SystemConfig) Save() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("config_.yaml", b, 0644)
}

func loadConfig(path string) (*SystemConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config SystemConfig
	err = yaml.Unmarshal(b, &config)
	return &config, err
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
	Addr        string
}
