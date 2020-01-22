package goback

import (
	"encoding/json"
	"fmt"
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
	b.writeBackupState(Completed)

	data, err := converter.EncodeToBytes(b.summaries)
	if err != nil {
		return err
	}

	if err := b.summaryDb.Truncate(0); err != nil {
		return err
	}
	if _, err := b.summaryDb.WriteAt(data, 0); err != nil {
		return err
	}

	// 	log.Info(strings.Repeat("=", 50))
	// 	log.WithFields(log.Fields{
	// 		"ID": b.summary.Id,
	// 	}).Info("# summary")
	// 	log.WithFields(log.Fields{
	// 		"files":     b.summary.TotalCount,
	// 		"totalSize": GetHumanizedSize(b.summary.TotalSize),
	// 		"execTime":  b.summary.ExecutionTime,
	// 	}).Info("# summary")
	//
	// 	log.WithFields(log.Fields{
	// 		"backupFailed": b.summary.FailedCount,
	// 		"size":         GetHumanizedSize(b.summary.FailedSize),
	// 	}).Info("# summary")
	// 	log.WithFields(log.Fields{
	// 		"backupSuccess": b.summary.SuccessCount,
	// 		"size":          GetHumanizedSize(b.summary.SuccessSize),
	// 	}).Info("# summary")
	// 	log.Info(strings.Repeat("=", 50))
	//

	bdd, _ := json.Marshal(b.summary)
	fmt.Println(string(bdd))
	return nil
}

func (b *Backup) writeChangesLog(lastFileMap *sync.Map) error {
	m := make(map[string]interface{})

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
			failed = append(failed, file)
			return true
		})
	}

	m["added"] = added
	m["modified"] = modified
	m["failed"] = failed
	m["deleted"] = deleted
	m["summary"] = b.summary

	// Write changes log
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	path := filepath.Join(b.tempDir, "changes-"+b.srcDirMap[b.summary.SrcDir].Checksum+".json")
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}
