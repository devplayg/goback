package main

import (
	"flag"
	"fmt"
	"github.com/devplayg/goback"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

const (
	ProductName = "backup"
	Version     = "2.0.0"
)

var (
	fs              = flag.NewFlagSet("", flag.ExitOnError)
	srcDir          = fs.String("s", "", "Source directory")
	dstDir          = fs.String("d", "", "Destination directory")
	version         = fs.Bool("v", false, "Version")
	debug           = fs.Bool("debug", false, "Debug")
	hashComparision = fs.Bool("hash", false, "Hash comparison")
)

func main() {
	backup := goback.NewBackup(*srcDir, *dstDir, *hashComparision, *debug)
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
	fmt.Printf("backup v%s\n", Version)
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
		fmt.Printf("%s %s\n", ProductName, Version)
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
