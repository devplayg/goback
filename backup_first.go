package goback

import (
	"fmt"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"sync"
)

func (b *Backup) generateFirstBackupData(srcDir string) error {
	log.Infof("generating source data from %s", srcDir)

	// 1. Issue summary
	b.issueSummary(srcDir, InitialBackup)

	// 2. Collect files in source directories
	currentFileMap, err := b.collectFilesToBackup(srcDir)
	if err != nil {
		return err
	}

	// 3. Write result
	if err := b.writeResult([]*sync.Map{currentFileMap}, nil); err != nil {
		return err
	}

	// 4. Write summary
	if err := b.writeSummary(); err != nil {
		log.Error(err)
	}

	return nil
}

func (b *Backup) collectFilesToBackup(srcDir string) (*sync.Map, error) {
	fileMap, extensions, sizeDistribution, count, size, err := GetFileMap(srcDir, b.hashComparision)
	if err != nil {
		return fileMap, err
	}
	b.summary.TotalCount = count
	b.summary.TotalSize = size
	b.summary.AddedCount = uint64(count)
	b.summary.AddedSize = size
	b.summary.Extensions = extensions
	b.summary.SizeDistribution = sizeDistribution
	b.writeBackupState(Read)
	log.WithFields(log.Fields{
		"execTime": b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
		"files":    b.summary.TotalCount,
		"dir":      srcDir,
		"size":     fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
	}).Info("  - files loaded")
	return fileMap, nil
}
