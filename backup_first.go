package goback

import (
	"fmt"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (b *Backup) generateFirstBackupData() error {
	defer func() {
		if err := b.writeSummary(); err != nil {
			log.Error(err)
		}
	}()

	// Ready
	log.Info("generating source data..")
	b.summary = NewSummary(b.nextSummaryId, b.srcDirArr, b.dstDir, b.workerCount, b.version)

	if err := b.createBackupDir(); err != nil {
		return err
	}

	// Collect files in source directories
	fileMap, err := b.collectFilesToBackup()
	if err != nil {
		return err
	}

	// No comparison
	b.summary.ComparisonTime = b.summary.ReadingTime

	// No backup
	b.summary.BackupTime = b.summary.ComparisonTime

	// Write result
	if err := b.writeResult([]*sync.Map{fileMap}); err != nil {
		return err
	}

	return nil
}

func (b *Backup) collectFilesToBackup() (*sync.Map, error) {
	defer func() {
		b.summary.ReadingTime = time.Now()

		log.WithFields(log.Fields{
			"execTime": b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
			"files":    b.summary.TotalCount,
			"size":     fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
		}).Info("files loaded")

	}()
	fileMap, extensions, sizeDistribution, count, size, err := GetFileMap(b.srcDirArr, b.hashComparision)
	if err != nil {
		return fileMap, err
	}
	b.summary.TotalCount = count
	b.summary.TotalSize = size
	b.summary.Extensions = extensions
	b.summary.SizeDistribution = sizeDistribution

	return fileMap, nil
}
