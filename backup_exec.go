package goback

import (
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
)

func (b *Backup) startBackup() error {
	tempDir, err := ioutil.TempDir(b.dstDir, "bak")
	if err != nil {
		return err
	}
	b.tempDir = tempDir

	lastFileMap, err := b.getLastFileMap()
	if err != nil {
		return err
	}
	//
	//currentFileMap, size, err := GetFileMap(b.srcDir, b.hashComparision)
	//if err != nil {
	//    return nil
	//}
	//
	//
	//if err := b.compareFiles(); err != nil {
	//	return err
	//}

	spew.Dump(lastFileMap)
	return nil
}

func aaa() {
	//	b.writeToDatabase(newMap, sync.Map{})
	//	b.summary.LoggingTime = time.Now()
	//	return nil
	//}
	//	b.summary.ReadingTime = time.Now()
	//
	//	// Search files and compare with previous data
	//	log.Infof("comparing old and new")
	//	b.summary.State = 3
	//	i := 1
	//	err := filepath.Walk(b.srcDir, func(path string, f os.FileInfo, err error) error {
	//		if !f.IsDir() && f.Mode().IsRegular() {
	//
	//			log.Debugf("Start checking: [%d] %s (%d)", i, path, f.Size())
	//			atomic.AddUint32(&b.summary.TotalCount, 1)
	//			atomic.AddUint64(&b.summary.TotalSize, uint64(f.Size()))
	//			fi := newFile(path, f.Size(), f.ModTime())
	//
	//			if inf, ok := originMap.Load(path); ok {
	//				last := inf.(*File)
	//
	//				if last.ModTime.Unix() != f.ModTime().Unix() || last.Size != f.Size() {
	//					log.Debugf("modified: %s", path)
	//					fi.State = FileModified
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
	//				}
	//				originMap.Delete(path)
	//			} else {
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
	//			}
	//			//if fi.State < 0 {
	//			//	log.Debugf("[%d] %s", fi.State, fi.Path)
	//			//}
	//			newMap.Store(path, fi)
	//			i++
	//		}
	//		return nil
	//	})
	//	checkErr(err)
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
}
