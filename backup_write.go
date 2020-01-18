package goback

import (
	"encoding/json"
	"github.com/devplayg/golibs/converter"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func (b *Backup) writeResult(currentFileMaps []*sync.Map) error {
	defer func() {
		b.summary.LoggingTime = time.Now()
		log.WithFields(log.Fields{
			"files": b.summary.TotalCount,
			// "changeFiles": GetChangeFilesDesc(b.summary.AddedCount, b.summary.ModifiedCount, b.summary.DeletedCount),
			"execTime": time.Since(b.summary.Date).Seconds(),
		}).Info("current files recorded")
	}()

	if err := b.writeChangesLog(); err != nil {
		return err
	}

	// wondory
	if err := b.writeFileMap(currentFileMaps); err != nil {
		return err
	}

	return nil
}

func (b *Backup) writeFileMap(fileMaps []*sync.Map) error {
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
	if err := b.fileMapDb.Truncate(0); err != nil {
		return err
	}

	if _, err := b.fileMapDb.WriteAt(data, 0); err != nil {
		return err
	}

	return err
}

func (b *Backup) writeSummary() error {
	b.summaries = append(b.summaries, b.summary)
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
	b.summary.ExecutionTime = b.summary.LoggingTime.Sub(b.summary.Date).Seconds()

	// Write changes log
	data, err = json.MarshalIndent(b.changesLog, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(b.backupDir, "changes.json"), data, 0644); err != nil {
		return err
	}

	log.Info(strings.Repeat("=", 50))
	log.WithFields(log.Fields{
		"summaryId": b.summary.Id,
		"files":     b.summary.TotalCount,
		"totalSize": GetHumanizedSize(b.summary.TotalSize),
		"execTime":  b.summary.ExecutionTime,
	}).Info("# summary")

	log.WithFields(log.Fields{
		"backupFailed": b.summary.FailedCount,
		"size":         GetHumanizedSize(b.summary.FailedSize),
	}).Info("# summary")
	log.WithFields(log.Fields{
		"backupSuccess": b.summary.SuccessCount,
		"size":          GetHumanizedSize(b.summary.SuccessSize),
	}).Info("# summary")
	log.Info(strings.Repeat("=", 50))

	return nil
}

func (b *Backup) writeChangesLog() error {
	m := make(map[string]interface{})

	added := make([]*FileWrapper, 0)
	b.addedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		added = append(added, file)
		return true
	})

	modified := make([]*FileWrapper, 0)
	b.modifiedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		modified = append(modified, file)
		return true
	})

	failed := make([]*FileWrapper, 0)
	b.failedFiles.Range(func(k, v interface{}) bool {
		file := k.(*FileWrapper)
		failed = append(failed, file)
		return true
	})

	// The remaining files in LastFileMap are deleted files.
	deleted := make([]*FileWrapper, 0)
	if b.lastFileMap != nil {
		b.lastFileMap.Range(func(k, v interface{}) bool {
			file := v.(*File)

			fileWrapper := FileWrapper{
				File:         file,
				WhatHappened: FileDeleted,
				Result:       0,
				Duration:     0,
				Message:      "",
			}
			deleted = append(deleted, &fileWrapper)
			return true
		})
	}

	m["added"] = added
	m["modified"] = modified
	m["failed"] = failed
	m["deleted"] = deleted
	m["summary"] = b.summary

	b.changesLog = m
	return nil

	// data, err := json.MarshalIndent(m, "", "    ")
	// if err != nil {
	// 	return err
	// }

	// return ioutil.WriteFile(filepath.Join(b.backupDir, "changes.json"), data, 0644)
}
