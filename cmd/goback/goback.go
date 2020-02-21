package main

import (
	"github.com/devplayg/goback"
	"time"
)

func main() {
	config := goback.AppConfig{
		Name:        "goback",
		Description: "smart bak",
		Version:     "service beta",
		Url:         "https://",
		Text1:       "ph1",
		Text2:       "ph2",
		Year:        time.Now().Year(),
		Company:     "dev",
		Debug:       true,
		Trace:       false,
	}
	engine := goback.NewEngine(&config)
	if err := engine.Start(); err != nil {
		panic(err)
	}

	//config, _ := goback.LoadConfig("config.yaml")
	//spew.Dump(config)
	//config.Save()

	// full,
	// initial,
	// incremental

	// local
	// remote
	// remote

	//keepers := make([]goback.Keeper, 0)
	//keepers = append(keepers, goback.NewLocalKeeper("d:/backup"))
	//keepers = append(keepers, goback.NewSftpKeeper("127.0.0.1", 22, "devplayg", "devplayg123!@#", "/backup1"))
	//keepers = append(keepers, goback.NewSftpKeeper("127.0.0.1", 22, "devplayg", "devplayg123!@#", "/backup2"))
	//
	//srcDirs := []string{
	//	"c:/temp",
	//	"d:/data",
	//	"d:/temp",
	//	"d:/드라마",
	//	"D:/Dropbox/01 - Music",
	//}
	//log.SetLevel(log.DebugLevel)
	//backup := goback.NewBackup(srcDirs, keepers, goback.Incremental, true)
	//if err := backup.Start(); err != nil {
	//	log.Error(err)
	//}
	//
	//app := himma.Application{
	//	AppName:     "SecuBACKUP",
	//	Description: "INCREMENTAL BACKUP ",
	//	Url:         "https://devplayg.com",
	//	Phrase1:     "KEEP YOUR DATA SAFE",
	//	Phrase2:     "Powered by Go",
	//	Year:        time.Now().Year(),
	//	Version:     appVersion,
	//	Company:     "SECUSOLUTION",
	//}
	//c := goback.NewController(backup.DbDir, "127.0.0.1:8000", &app)
	//if err := c.Start(); err != nil {
	//	log.Error(err)
	//}
}
