package goback

// import (
// 	"fmt"
// 	"github.com/dustin/go-humanize"
// 	log "github.com/sirupsen/logrus"
// 	"sync"
// 	"time"
// )
//
// func (b *Backup) generateFirstBackupData() error {
// 	log.Info("generating source data..")
//
// 	// 1. Collect files in source directories
// 	fileMap, err := b.collectFilesToBackup()
// 	if err != nil {
// 		return err
// 	}
//
// 	// No comparison
// 	b.summary.ComparisonTime = b.summary.ReadingTime
//
// 	// No backup
// 	b.summary.BackupTime = b.summary.ComparisonTime
//
// 	// 2. Write result
// 	if err := b.writeResult([]*sync.Map{fileMap}); err != nil {
// 		return err
// 	}
//
// 	// 3. Write summary
// 	if err := b.writeSummary(); err != nil {
// 		log.Error(err)
// 	}
//
// 	return nil
// }
//
// func (b *Backup) collectFilesToBackup() (*sync.Map, error) {
// 	defer func() {
// 		b.summary.ReadingTime = time.Now()
//
// 		log.WithFields(log.Fields{
// 			"execTime": b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
// 			"files":    b.summary.TotalCount,
// 			"size":     fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
// 		}).Info("files loaded")
//
// 	}()
// 	fileMap, extensions, sizeDistribution, count, size, err := GetFileMap(b.srcDirs, b.hashComparision)
// 	if err != nil {
// 		return fileMap, err
// 	}
// 	b.summary.TotalCount = count
// 	b.summary.TotalSize = size
// 	b.summary.Extensions = extensions
// 	b.summary.SizeDistribution = sizeDistribution
//
// 	return fileMap, nil
// }
