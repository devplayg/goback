package goback

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
	}).Debug("new backup is ready")
	defer func() {
		if err := b.writeSummary(); err != nil {
			log.Error(err)
		}
	}()

	tempDir, err := ioutil.TempDir(b.dstDir, "bak")
	if err != nil {
		return err
	}
	b.tempDir = tempDir

	// Reading
	lastFileMap, err := b.getLastFileMap()
	if err != nil {
		return err
	}

	currentFileMap, _, err := GetFileMap(b.srcDirArr, b.hashComparision)
	if err != nil {
		return nil
	}
	b.summary.ReadingTime = time.Now()

	if err := b.compareFileMaps(lastFileMap, currentFileMap); err != nil {
		return err
	}
	b.summary.ComparisonTime = time.Now()

	// Write
	//if err := b.writeFileMap(fileMap); err != nil {
	//	return err
	//}

	//spew.Dump(lastFileMap)
	return nil
}

func (b *Backup) compareFileMaps(lastFileMap, currentFileMap map[string]*File) error {
	for path, current := range currentFileMap {
		b.summary.TotalCount++
		b.summary.TotalSize += current.Size

		if last, had := lastFileMap[path]; had {
			//				last := inf.(*File)

			//
			if last.ModTime.Unix() != current.ModTime.Unix() || last.Size != current.Size {
				log.Debugf("modified: %s", path)
				fi.State = FileModified
				//					atomic.AddUint32(&b.summary.BackupModified, 1)
				//					backupPath, dur, err := b.BackupFile(path)
				//					if err != nil {
				//						atomic.AddUint32(&b.summary.BackupFailure, 1)
				//						log.Error(err)
				//						fi.Message = err.Error()
				//						fi.State = fi.State * -1
				//						//spew.Dump(fi)
				//					} else {
				//						fi.Message = fmt.Sprintf("copy_time=%4.1f", dur)
				//						atomic.AddUint32(&b.summary.BackupSuccess, 1)
				//						atomic.AddUint64(&b.summary.BackupSize, uint64(f.Size()))
				//						os.Chtimes(backupPath, f.ModTime(), f.ModTime())
				//						originMap.Delete(path)
				//					}
			}
			//				originMap.Delete(path)

		} else {
			//				log.Debugf("added: %s", path)
			//				fi.State = FileAdded
			//				atomic.AddUint32(&b.summary.BackupAdded, 1)
			//				backupPath, dur, err := b.BackupFile(path)
			//				if err != nil {
			//					atomic.AddUint32(&b.summary.BackupFailure, 1)
			//					log.Error(err)
			//					fi.Message = err.Error()
			//					fi.State = fi.State * -1
			//					//spew.Dump(fi)
			//				} else {
			//					fi.Message = fmt.Sprintf("copy_time=%4.1f", dur)
			//					atomic.AddUint32(&b.summary.BackupSuccess, 1)
			//					atomic.AddUint64(&b.summary.BackupSize, uint64(f.Size()))
			//					os.Chtimes(backupPath, f.ModTime(), f.ModTime())
			//				}

		}
	}

	//		if !f.IsDir() && f.Mode().IsRegular() {
	//			if inf, ok := originMap.Load(path); ok {
	//			} else {
	//			}
	//			//if fi.State < 0 {
	//			//	log.Debugf("[%d] %s", fi.State, fi.Path)
	//			//}
	//			newMap.Store(path, fi)
	//			i++
	//		}
	//		return nil
	//	})
	//
	//	// Rename directory
	//	lastDir := filepath.Join(b.dstDir, b.summary.Date.Format("20060102"))
	//	err = os.Rename(b.tempDir, lastDir)
	//	if err == nil {
	//		b.summary.DstDir = lastDir
	//	} else {
	//
	//		i := 1
	//		for err != nil && i <= 10 {
	//			altDir := lastDir + "_" + strconv.Itoa(i)
	//			err = os.Rename(b.tempDir, altDir)
	//			if err == nil {
	//				b.summary.DstDir = altDir
	//			}
	//			i += 1
	//		}
	//		if err != nil {
	//			b.summary.Message = err.Error()
	//			b.summary.State = -1
	//			b.summary.DstDir = b.tempDir
	//			os.RemoveAll(b.tempDir)
	//			return err
	//		}
	//	}
	//	b.summary.ComparisonTime = time.Now()
	//
	//	// Write data to database
	//	err = b.writeToDatabase(newMap, originMap)
	//	b.summary.LoggingTime = time.Now()
	//	return err
	return nil
}
