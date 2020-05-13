package main

import (
	"github.com/devplayg/goback"
	"time"
)

func main() {
	engine := goback.NewEngine(&goback.AppConfig{
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
		Address:     ":8000",
	})

	// Start
	if err := engine.Start(); err != nil {
		panic(err)
	}
}
