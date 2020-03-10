package main

import (
	"fmt"
	"github.com/devplayg/goback"
	"github.com/devplayg/himma"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"runtime"
	"time"
)

const (
	appName    = "goback"
	appVersion = "1.1.0"
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
		app := himma.Application{
			AppName:     "SecuBACKUP",
			Description: "INCREMENTAL BACKUP ",
			Url:         "https://devplayg.com",
			Phrase1:     "KEEP YOUR DATA SAFE",
			Phrase2:     "Powered by Go",
			Year:        time.Now().Year(),
			Version:     appVersion,
			Company:     "SECUSOLUTION",
		}
		c := goback.NewController(*web, *addr, &app)
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
		fmt.Printf("%s v%s\n", appName, appVersion)
		fs.PrintDefaults()
		// fmt.Println("\n  usage) backup -s [directory to backup] -d [directory where backup files will be stored]")
		fmt.Println("")
		fmt.Printf("  ex) goback -s /dir/to/backup -d /dir/to/be/saved\n")
		fmt.Printf("  ex) goback -s /dir/to/backup1 -s /dir/to/backup2 -d /dir/to/be/saved\n")
		os.Exit(0)
	}
	_ = fs.Parse(os.Args[1:])

	if *version || (len(*web) < 1 && len(*srcDirs) < 1 && len(*dstDir) < 1) {
		fs.Usage()
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
