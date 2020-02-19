package goback

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	App  AppConfig
	Jobs []Job `json:"job"`
}

type Job struct {
	Dir      string `json:"dir"`
	Schedule struct {
		Full        []string `json:"full"`
		Incremental []string `json:"incremental"`
	}
	Storage []Storage `json:"storage"`
	Ignore  []string  `json:"ignore"`
}

type Storage struct {
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
