package goback

import (
	log "github.com/sirupsen/logrus"
)

func (b *Backup) generateFirstBackupData(srcDir string) error {
	log.Infof("generating source data from %s", srcDir)

	// 1. Issue summary
	b.issueSummary(srcDir, InitialBackup)

	// 2. Collect files in source directories
	currentFileMaps, err := b.getCurrentFileMaps(srcDir)
	if err != nil {
		return err
	}

	b.summary.ComparisonTime = b.summary.ReadingTime
	b.summary.BackupTime = b.summary.ComparisonTime

	// 3. Write result
	if err := b.writeResult(currentFileMaps, nil); err != nil {
		return err
	}

	return nil
}
