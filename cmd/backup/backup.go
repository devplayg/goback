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
	appName    = "backup"
	appVersion = "2.0.0"
)

var (
	fs              = pflag.NewFlagSet(appName, pflag.ContinueOnError)
	srcDirArr       = fs.StringArrayP("src", "s", []string{}, "Source directories")
	dstDir          = fs.StringP("dst", "d", "", "Destination directory")
	debug           = fs.Bool("debug", false, "Debug")
	verbose         = fs.BoolP("verbose", "v", false, "Verbose")
	version         = fs.Bool("version", false, "Version")
	hashComparision = fs.Bool("hash", false, "Hash comparison")
)

func main() {
	backup := goback.NewBackup(*srcDirArr, *dstDir, *hashComparision, *debug)
	if err := backup.Start(); err != nil {
		log.Error(err)
	}

	////	Start backup files
	//b := goback.NewBackup(*srcDir, *dstDir, *debug)
	//defer b.Close()
	//
	//// Initialize backup
	//if err := b.Initialize(); err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//// Start backup
	//if err = b.Start(); err != nil {
	//	log.Error(err)
	//}
}

func printHelp() {
	fmt.Printf("backup v%s\n", appVersion)
	fmt.Println("Description: Incremental backup")
	fmt.Println("Usage: backup -s [directory to backup] -d [directory where backup files will be stored]")
	fmt.Println("Usage: backup -s /data -d /backup")
	fs.PrintDefaults()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fs.Usage = printHelp
	_ = fs.Parse(os.Args[1:])
	if *version {
		fmt.Printf("%s %s\n", appName, appVersion)
		return
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
