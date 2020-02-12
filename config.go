package goback

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	appName string
	appDesc string
	version string
	Backup  []BackupDir `json:"backup"`
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

func (c *Config) Save() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("config_.yaml", b, 0644)
}

func loadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	return &config, err
}
