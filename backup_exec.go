package goback

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

func (b *Backup) startBackup() error {

	// Ready
	summary, err := b.newSummary()
	if err != nil {
		return err
	}
	b.summary = summary
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
	currentFileMaps, extensions, sizeDistribution, count, size, err := GetCurrentFileMaps(b.srcDirArr, b.workerCount, b.hashComparision)
	if err != nil {
		return nil
	}
	b.summary.ReadingTime = time.Now()
	b.summary.TotalCount = count
	b.summary.TotalSize = size
	b.summary.Extensions = extensions
	b.summary.SizeDistribution = sizeDistribution
	log.Trace(lastFileMap)

	if err := b.CompareFileMaps(lastFileMap, currentFileMaps); err != nil {
		return err
	}

	// Write
	if err := b.writeFileMap(currentFileMaps); err != nil {
		return err
	}

	//spew.Dump(lastFileMap)
	return nil
}

func (b *Backup) writeWhatHappened(file *File, whatHappened int) {
	file.WhatHappened = whatHappened
	if whatHappened == FileAdded {
		b.summary.addedFiles.Store(file, nil)
		atomic.AddUint64(&b.summary.AddedCount, uint64(1))
		return
	}
	if whatHappened == FileModified {
		b.summary.modifiedFiles.Store(file, nil)
		atomic.AddUint64(&b.summary.ModifiedCount, uint64(1))
		return
	}
	if whatHappened == FileDeleted {
		atomic.AddUint64(&b.summary.DeletedCount, uint64(1))
		b.summary.deletedFiles.Store(file, nil)
		return
	}
}

func (b *Backup) CompareFileMaps(lastFileMap *sync.Map, currentFileMaps []*sync.Map) error {
	wg := sync.WaitGroup{}
	for i := range currentFileMaps {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			//log.WithFields(log.Fields{
			//	"files": len(m),
			//}).Debugf("worker-%d has been started", workerId)

			if err := b.compareFileMap(workerId, lastFileMap, currentFileMaps[workerId]); err != nil {
				log.Error(err)
			}
		}(i)
	}
	wg.Wait()

	// for _, v := range lastFileMap {
	lastFileMap.Range(func(k, v interface{}) bool {
		file := v.(*File)
		b.writeWhatHappened(file, FileDeleted)
		return true
	})

	return nil
}

func (b *Backup) compareFileMap(workerId int, lastFileMap, myMap *sync.Map) error {
	myMap.Range(func(k, v interface{}) bool {
		path := k.(string)
		current := v.(*File)

		if val, have := lastFileMap.Load(path); have {
			last := val.(*File)
			if last.ModTime.Unix() != current.ModTime.Unix() || last.Size != current.Size {
				//log.WithFields(log.Fields{
				//	"workerId": workerId,
				//}).Debugf("modified: %s", path)
				b.writeWhatHappened(current, FileModified)
			}
			lastFileMap.Delete(path)
			return true
		}
		//log.WithFields(log.Fields{
		//	"workerId": workerId,
		//}).Debugf("added: %s", path)

		b.writeWhatHappened(current, FileAdded)
		return true
	})
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
