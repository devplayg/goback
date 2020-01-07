package goback

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (b *Backup) startBackup() error {
	// Ready
	summary, err := b.newSummary()
	if err != nil {
		return err
	}
	b.summary = summary
	log.WithFields(log.Fields{
		"summaryId": summary.Id,
	}).Debug("backup is ready")
	defer func() {
		if err := b.writeSummary(); err != nil {
			log.Error(err)
		}
	}()

	//tempDir, err := ioutil.TempDir(b.dstDir, "bak")
	//if err != nil {
	//	return err
	//}
	//b.tempDir = tempDir

	// Reading
	lastFileMap, err := b.getLastFileMap()
	if err != nil {
		return err
	}

	currentFileMap, size, err := GetFileMap(b.srcDirArr, b.hashComparision)
	if err != nil {
		return nil
	}
	b.summary.ReadingTime = time.Now()
	b.summary.TotalCount = len(currentFileMap)
	b.summary.TotalSize = size

	added, modified, deleted, err := CompareFileMaps(lastFileMap, currentFileMap)
	if err != nil {
		return err
	}
	b.summary.ComparisonTime = time.Now()
	b.summary.AddedFiles = len(added)
	b.summary.ModifiedFiles = len(modified)
	b.summary.DeletedFiles = len(deleted)

	// Write
	if err := b.writeFileMap(currentFileMap); err != nil {
		return err
	}

	//spew.Dump(lastFileMap)
	return nil
}

//func (b *Backup) compareFileMaps(lastFileMap, currentFileMap map[string]*File) ([]*File, []*File, []*File, error) {
//	added := make([]*File, 0)
//	deleted := make([]*File, 0)
//	modified := make([]*File, 0)
//	for path, current := range currentFileMap {
//		if last, had := lastFileMap[path]; had {
//			if last.ModTime.Unix() != current.ModTime.Unix() || last.Size != current.Size {
//				log.Debugf("modified: %s", path)
//				current.WhatHappened = FileModified
//				modified = append(modified, current)
//			}
//			delete(lastFileMap, path)
//
//		} else {
//			log.Debugf("added: %s", path)
//			current.WhatHappened = FileAdded
//			added = append(added, current)
//		}
//	}
//	for _, file := range lastFileMap {
//		file.WhatHappened = FileDeleted
//		deleted = append(deleted, file)
//	}
//
//	log.WithFields(log.Fields{
//		"added":    len(added),
//		"modified": len(modified),
//		"deleted":  len(deleted),
//	}).Debugf("total %d files; comparision result", len(currentFileMap))
//
//	return added, modified, deleted, nil
//}
