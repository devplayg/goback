package goback

import (
	"fmt"
	"github.com/devplayg/goutils"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func (b *Backup) writeResult(currentFileMaps []*sync.Map, lastFileMap *sync.Map) error {
	if lastFileMap != nil {
		if err := b.writeChangesLog(lastFileMap); err != nil {
			return err
		}
	}

	if err := b.writeFileMaps(currentFileMaps); err != nil {
		return err
	}

	b.writeBackupState(Logged)

	b.summary.Message = fmt.Sprintf("%3.1fs / %3.1fs / %3.1fs / %3.1fs",
		b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
		b.summary.ComparisonTime.Sub(b.summary.ReadingTime).Seconds(),
		b.summary.BackupTime.Sub(b.summary.ComparisonTime).Seconds(),
		b.summary.LoggingTime.Sub(b.summary.BackupTime).Seconds(),
	)
	return nil
}

func (b *Backup) writeFileMaps(fileMaps []*sync.Map) error {
	files := make([]*File, 0) // test
	for _, m := range fileMaps {
		m.Range(func(k, v interface{}) bool {
			fileWrapper := v.(*FileWrapper)
			files = append(files, fileWrapper.File) // test
			return true
		})
	}

	data, err := goutils.GobEncode(files)
	if err != nil {
		return fmt.Errorf("failed to encode file map: %w", err)
	}

	return ioutil.WriteFile(b.srcDirMap[b.summary.SrcDir].dbPath, data, 0644)
}

func (b *Backup) writeChangesLog(lastFileMap *sync.Map) error {
	m := make(map[string]*StatsReportWithList)

	added := make([]*FileGrid, 0)
	b.summary.addedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		added = append(added, file.WrapInFileGrid())
		return true
	})

	modified := make([]*FileGrid, 0)
	b.summary.modifiedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		modified = append(modified, file.WrapInFileGrid())
		return true
	})

	failed := make([]*FileGrid, 0)
	b.summary.failedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		failed = append(failed, file.WrapInFileGrid())
		return true
	})

	// The remaining files in LastFileMap are deleted files.
	deleted := make([]*FileGrid, 0)
	if lastFileMap != nil {
		b.summary.deletedFiles.Range(func(k, v interface{}) bool {
			file := k.(*FileWrapper)
			deleted = append(deleted, file.WrapInFileGrid())
			return true
		})
	}

	m["added"] = CreateFilesReportWithList(added, b.summary.AddedSize, 0, b.rank)
	m["modified"] = CreateFilesReportWithList(modified, b.summary.ModifiedSize, 0, b.rank)
	m["failed"] = CreateFilesReportWithList(failed, b.summary.FailedSize, 0, b.rank)
	m["deleted"] = CreateFilesReportWithList(deleted, b.summary.DeletedSize, 0, b.rank)

	path := filepath.Join(b.dbDir, fmt.Sprintf(ChangesDbName, b.Id, b.srcDirMap[b.summary.SrcDir].checksum))
	if err := WriteBackupData(m, path, JsonEncoding); err != nil {
		return err
	}

	return nil
}
