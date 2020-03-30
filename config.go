package goback

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
	Enabled    bool     `json:"enabled"`
	Storage    *Storage `json:"-"`
}

//
// func (c *Config) Save() error {
// 	b, err := yaml.Marshal(c)
// 	if err != nil {
// 		return err
// 	}
//
// 	return ioutil.WriteFile("config_.yaml", b, 0644)
// }

type Storage struct {
	Id       int    `json:"id"`
	Protocol int    `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dir      string `json:"dir"`
}

func (c *Config) findJobById(id int) *Job {
	var job *Job
	for i, j := range c.Jobs {
		if j.Id == id {
			job = &c.Jobs[i]
			break
		}
	}

	if job != nil {
		storage := c.findStorageById(job.StorageId)
		if storage == nil {
			log.Error("Storage not found")
			return nil
		}
		job.Storage = storage
	}
	return job
}

func (c *Config) findStorageById(id int) *Storage {
	for i, storage := range c.Storages {
		if storage.Id == id {
			return &c.Storages[i]
		}
	}
	return nil
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
