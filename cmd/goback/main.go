package main

import (
	"fmt"
	"github.com/devplayg/goback"
	"github.com/devplayg/himma/v2"
	"github.com/devplayg/hippo/v2"
	"github.com/spf13/pflag"
	"os"
	"strconv"
	"time"
)

const (
	appName        = "GoBack"
	appDescription = "Smart incremental backup"
	appVersion     = "3.1.0"
)

var (
	fs      = pflag.NewFlagSet(appName+" "+appVersion, pflag.ExitOnError)
	debug   = fs.Bool("debug", false, "Debug") // GODEBUG=http2debug=2
	verbose = fs.BoolP("verbose", "v", false, "Verbose")
	version = fs.Bool("version", false, "Version")
	// insecure     = fs.Bool("insecure", false, "Disable TLS")
	// certFile     = fs.String("certFile", "server.crt", "SSL Certificate file")
	// keyFile      = fs.String("keyFile", "server.key", "SSL Certificate key file")
	port = fs.Uint16("port", 8000, "Monitoring address")
	// securePort   = fs.Uint16("secureport", 8443, "Secured monitoring address")
	resetAccount = fs.Bool("resetkey", false, "Reset access key and secret key")
)

func main() {
	appConfig := &goback.AppConfig{
		HimmaConfig: &himma.Config{
			AppName:     appName,
			Description: appDescription,
			Url:         "https://github.com/devplayg/goback",
			Phrase1:     "MAKE YOUR DATA SAFE",
			Phrase2:     "powered by Go",
			Year:        time.Now().Year(),
			Version:     appVersion,
			Company:     "devplayg",
		},
		HippoConfig: &hippo.Config{
			Name:        appName,
			Description: appDescription,
			Version:     appVersion,
			Debug:       *debug,
			Trace:       false,
			// CertFile:    *certFile,
			// KeyFile:     *keyFile,
			// Insecure:    *insecure,
		},
		Address: ":" + strconv.Itoa(int(*port)), // HTTP
		// SecureAddress: ":" + strconv.Itoa(int(*securePort)), // HTTPS
		Verbose:      *verbose,
		ResetAccount: *resetAccount,
	}

	// Start
	if err := goback.NewEngine(appConfig).Start(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	_ = fs.Parse(os.Args[1:])

	if *version {
		fmt.Printf("%s %s\n", appName, appVersion)
		os.Exit(0)
	}
}
