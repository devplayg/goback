package main

import (
	"fmt"
	"github.com/devplayg/goback"
	"github.com/spf13/pflag"
	"os"
	"strconv"
	"time"
)

const (
	appName        = "SecuBACKUP"
	appDescription = "Smart incremental backup"
	appVersion     = "3.0.1"
)

var (
	fs    = pflag.NewFlagSet(appName+" "+appVersion, pflag.ExitOnError)
	debug = fs.Bool("debug", false, "Debug") // GODEBUG=http2debug=2
	// trace        = fs.Bool("trace", false, "Trace")
	verbose = fs.BoolP("verbose", "v", false, "Verbose")
	version = fs.Bool("version", false, "Version")
	// insecure     = fs.Bool("insecure", false, "Disable TLS")
	// certFile     = fs.String("certFile", "server.crt", "SSL Certificate file")
	// keyFile      = fs.String("keyFile", "server.key", "SSL Certificate key file")
	// batchSize    = fs.Int("batchsize", 10000, "Batch size")
	// batchTimeout = fs.Int("batchtime", 1000, "Batch timeout, in milliseconds")
	// worker       = fs.Int("worker", 0, "Worker count")
	// monitor      = fs.Bool("mon", false, "Monitoring operation on HTTP")
	port    = fs.Uint16P("port", "p", 8000, "Monitoring address")
	devMode = fs.Bool("dev", false, "Developer mode")
)

func main() {
	engine := goback.NewEngine(&goback.AppConfig{
		Name:          appName,
		Description:   appDescription,
		Version:       appVersion,
		Url:           "https://github.com/devplayg/goback",
		Text1:         "MAKE YOUR DATA SAFE",
		Text2:         "powered by Go",
		Year:          time.Now().Year(),
		Company:       "SECUSOLUTION",
		Debug:         *debug,
		Trace:         false,
		Address:       ":" + strconv.Itoa(int(*port)),
		Verbose:       *verbose,
		DeveloperMode: *devMode,
	})

	// Start
	if err := engine.Start(); err != nil {
		panic(err)
	}
}

func init() {
	_ = fs.Parse(os.Args[1:])

	if *version {
		fmt.Printf("%s %s\n", appName, appVersion)
		os.Exit(0)
	}
}
