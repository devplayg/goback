package main

import (
	"fmt"
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
		fmt.Println(err.Error())
	}
}
