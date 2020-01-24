package goback

import (
	"fmt"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/golibs/converter"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func (b *Backup) writeResult(currentFileMaps []*sync.Map, lastFileMap *sync.Map) error {
	// 	defer func() {
	// 		b.summary.LoggingTime = time.Now()
	// 		log.WithFields(log.Fields{
	// 			"files": b.summary.TotalCount,
	// 			// "changeFiles": GetChangeFilesDesc(b.summary.AddedCount, b.summary.ModifiedCount, b.summary.DeletedCount),
	// 			"execTime": time.Since(b.summary.Date).Seconds(),
	// 		}).Info("current files recorded")
	// 	}()
	//
	if lastFileMap != nil {
		if err := b.writeChangesLog(lastFileMap); err != nil {
			return err
		}
	}

	if err := b.writeFileMaps(currentFileMaps); err != nil {
		return err
	}

	b.writeBackupState(Logged)
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

	data, err := converter.EncodeToBytes(files)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(b.srcDirMap[b.summary.SrcDir].dbPath, data, 0644)
}

func (b *Backup) writeSummary() error {
	encoded, err := converter.EncodeToBytes(b.summaries)
	if err != nil {
		return fmt.Errorf("failed to encode summary data: %w", err)
	}
	compressed, err := compress.Compress(encoded, compress.GZIP)
	if err != nil {
		return fmt.Errorf("failed to compress summarydata: %w", err)
	}

	if err := b.summaryDb.Truncate(0); err != nil {
		return err
	}
	if _, err := b.summaryDb.WriteAt(compressed, 0); err != nil {
		return err
	}

	return nil
}

func (b *Backup) writeChangesLog(lastFileMap *sync.Map) error {
	m := make(map[string]*StatsReportWithList)

	added := make([]*FileWrapper, 0)
	b.summary.addedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		added = append(added, file)
		return true
	})

	modified := make([]*FileWrapper, 0)
	b.summary.modifiedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		modified = append(modified, file)
		return true
	})

	failed := make([]*FileWrapper, 0)
	b.summary.failedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		failed = append(failed, file)
		return true
	})

	// The remaining files in LastFileMap are deleted files.
	deleted := make([]*FileWrapper, 0)
	if lastFileMap != nil {
		b.summary.deletedFiles.Range(func(k, v interface{}) bool {
			file := k.(*FileWrapper)
			deleted = append(deleted, file)
			return true
		})
	}

	m["added"] = CreateFilesReportWithList(added, 0, b.rank)
	m["modified"] = CreateFilesReportWithList(modified, 0, b.rank)
	m["failed"] = CreateFilesReportWithList(failed, 0, b.rank)
	m["deleted"] = CreateFilesReportWithList(deleted, 0, b.rank)

	if err := WriteBackupData(m, filepath.Join(b.tempDir, "changes-"+b.srcDirMap[b.summary.SrcDir].Checksum+".db"), JsonEncoding); err != nil {
		return err
	}
	// data, err := converter.EncodeToBytes(m)
	// if err != nil {
	//     return fmt.Errorf("failed to encode changes log: %w", err)
	// }
	// compressed, err := compress.Compress(data, compress.GZIP)
	// if err != nil {
	//     return fmt.Errorf("failed to compress encoded changes log: %w", err)
	// }
	// if err := ioutil.WriteFile(filepath.Join(b.tempDir, "changes-"+b.srcDirMap[b.summary.SrcDir].Checksum+".data"), compressed, 0644); err != nil {
	//     return err
	// }

	return nil
}
