package main

import (
	"fmt"
	"github.com/devplayg/goback"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"runtime"
)

const (
	appName    = "goback"
	appVersion = "1.0.3"
)

var (
	fs      = pflag.NewFlagSet(appName, pflag.ContinueOnError)
	srcDirs = fs.StringArrayP("src", "s", []string{}, "Source directories")
	dstDir  = fs.StringP("dst", "d", "", "Destination directory")
	debug   = fs.Bool("debug", false, "Debug")
	verbose = fs.BoolP("verbose", "v", false, "Verbose")
	version = fs.Bool("version", false, "Version")
	// hashComparision = fs.Bool("hash", false, "Hash comparison")
	web  = fs.StringP("web", "w", "", "Database directory")
	addr = fs.String("addr", "0.0.0.0:8000", "Listen address and port")
)

func main() {
	if len(*web) > 0 { // Web UI
		c := goback.NewController(*web, *addr, appVersion)
		if err := c.Start(); err != nil {
			log.Error(err)
		}
		return
	}

	backup := goback.NewBackup(*srcDirs, *dstDir, false, *debug)
	if err := backup.Start(); err != nil {
		log.Error(err)
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fs.Usage = func() {
		fmt.Printf("backup v%s\n", appVersion)
		fmt.Println("Description: Incremental backup")
		fmt.Println("Usage: backup -s [directory to backup] -d [directory where backup files will be stored]")
		fmt.Println("Usage: backup -s /data -d /backup")
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[1:])
	if *version || (len(*web) < 1 && len(*srcDirs) < 1 && len(*dstDir) < 1) {
		//fmt.Printf("%s %s\n", appName, appVersion)
		fmt.Printf("ex) goback -s /dir/to/backup -d /dir/to/be/saved\n")
		fmt.Printf("ex) goback -s /dir/to/backup1 -s /dir/to/backup2 -d /dir/to/be/saved\n")
		os.Exit(0)
	}

	initLogger()
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		DisableColors: true,
	})

	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}
