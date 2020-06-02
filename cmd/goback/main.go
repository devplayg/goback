package main

import (
	"github.com/devplayg/goback"
	"time"
)

func main() {
	engine := goback.NewEngine(&goback.AppConfig{
		Name:        "GoBack",
		Description: "Smart incremental backup",
		Version:     "2.0",
		Url:         "https://github.com/devplayg/goback",
		Text1:       "MAKE YOUR DATA SAFE",
		Text2:       "powered by Go",
		Year:        time.Now().Year(),
		Company:     "DevPlayG",
		Debug:       false,
		Trace:       false,
		Address:     ":8000",
	})

	// Start
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
