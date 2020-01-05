package goback

import (
	"crypto/sha256"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/minio/highwayhash"
	log "github.com/sirupsen/logrus"
	"hash"
	"path/filepath"
	"sync"
	"time"
)

const (
	FileModified = 1 << iota // 1
	FileAdded    = 1 << iota // 2
	FileDeleted  = 1 << iota // 4

	BackupReady    = 1
	BackupRunning  = 2
	BackupFinished = 3

	BackupSuccess = 1
)

var (
	BucketSummary = []byte("summary")
	BucketFiles   = []byte("files")
)

type Backup struct {
	srcDir          string
	dstDir          string
	db              *bolt.DB
	fileDb          *bolt.DB
	tempDir         string
	summary         *Summary
	hashComparision bool
	debug           bool
	key             []byte
	highwayhash     hash.Hash

	//dbOrigin   *sql.DB
	//dbOriginTx *sql.Tx
	//dbLog      *sql.DB
	//dbLogTx    *sql.Tx
}

func NewBackup(srcDir, dstDir string, hashComparision, debug bool) *Backup {
	b := Backup{
		//srcDir:       `filepath.Clean`(srcDir),
		//dstDir:       filepath.Clean(dstDir),
		//dbOriginFile: filepath.Join(filepath.Clean(dstDir), "backup_origin.db"),
		//dbLogFile:    filepath.Join(filepath.Clean(dstDir), "backup_log.db"),
		srcDir:          srcDir,
		dstDir:          dstDir,
		hashComparision: hashComparision,
		//dbOriginFile: filepath.Join(dstDir, "backup_origin.db"),
		//dbLogFile:    filepath.Join(dstDir, "backup_log.db"),
		debug: debug,
	}
	return &b
}

func (b *Backup) init() error {
	if err := b.initDirectories(); err != nil {
		return err
	}

	if err := b.initDatabase(); err != nil {
		return err
	}

	if b.hashComparision {
		key := sha256.Sum256([]byte("Duplicate File Finder"))
		highwayhash, err := highwayhash.New(key[:])
		if err != nil {
			return err
		}

		b.highwayhash = highwayhash
	}

	//b.summary = newSummary(0, b.srcDir)
	return nil
}

func (b *Backup) initDirectories() error {
	b.srcDir = filepath.Clean(b.srcDir)
	if err := isValidDir(b.srcDir); err != nil {
		return errors.New("source directory error: " + err.Error())
	}

	b.dstDir = filepath.Clean(b.dstDir)
	if err := isValidDir(b.dstDir); err != nil {
		return errors.New("destination directory error: " + err.Error())
	}

	//tempDir, err := ioutil.TempDir(b.dstDir, "bak")
	//if err != nil {
	//	return err
	//}
	//b.tempDir = tempDir

	log.WithFields(log.Fields{
		"srcDir": b.srcDir,
		"dstDir": b.dstDir,
	}).Debug("backup directories")

	return nil
}

//// Initialize directories
//func (b *Backup) 0() error {
//	tempDir, err := ioutil.TempDir(b.dstDir, "bak")
//	if err != nil {
//		return err
//	}
//	b.tempDir = tempDir
//
//	return nil
//}
//
//// Initialize database
func (b *Backup) initDatabase() error {

	db, err := bolt.Open(filepath.Join(b.dstDir, "backup_log.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(BucketSummary)
		return err
	})
	if err != nil {
		return err
	}
	b.db = db

	fileDb, err := bolt.Open(filepath.Join(b.dstDir, "backup_origin.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	err = fileDb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(BucketFiles)
		return err
	})
	if err != nil {
		return err
	}
	b.fileDb = fileDb

	//dbOriginFile: filepath.Join(filepath.Clean(dstDir), "backup_origin.db"),
	//dbLogFile:    filepath.Join(filepath.Clean(dstDir), "backup_log.db"),
	//	var err error
	//	var query string
	//
	//	// Set databases
	//	b.dbOrigin, err = sql.Open("sqlite3", b.dbOriginFile)
	//	if err != nil {
	//		return err
	//	}
	//	b.dbOriginTx, _ = b.dbOrigin.Begin()
	//	b.dbLog, err = sql.Open("sqlite3", b.dbLogFile)
	//	if err != nil {
	//		return err
	//	}
	//	b.dbLogTx, _ = b.dbLog.Begin()
	//
	//	// Original database
	//	query = `
	//		CREATE TABLE IF NOT EXISTS bak_origin (
	//			path text not null,
	//			size int not null,
	//			mtime text not null
	//		);
	//	`
	//	_, err = b.dbOrigin.Exec(query)
	//	if err != nil {
	//		return err
	//	}
	//
	//	// Log database
	//	query = `
	//		CREATE TABLE IF NOT EXISTS bak_summary (
	//			id integer not null primary key autoincrement,
	//			date integer not null  DEFAULT CURRENT_TIMESTAMP,
	//			src_dir text not null default '',
	//			dst_dir text not null default '',
	//			state integer not null default 0,
	//			total_size integer not null default 0,
	//			total_count integer not null default 0,
	//			backup_modified integer not null default 0,
	//			backup_added integer not null default 0,
	//			backup_deleted integer not null default 0,
	//			backup_success integer not null default 0,
	//			backup_failure integer not null default 0,
	//			backup_size integer not null default 0,
	//			execution_time real not null default 0.0,
	//			message text not null default ''
	//		);
	//
	//		CREATE INDEX IF NOT EXISTS ix_bak_summary ON bak_summary(date);
	//
	//		CREATE TABLE IF NOT EXISTS bak_log(
	//			id int not null,
	//			path text not null,
	//			size int not null,
	//			mtime text not null,
	//			state int not null,
	//			message text not null
	//		);
	//
	//		CREATE INDEX IF NOT EXISTS ix_bak_log_id on bak_log(id);
	//`
	//	_, err = b.dbLog.Exec(query)
	//	if err != nil {
	//		return err
	//	}

	return nil
}

func (b *Backup) getOriginMap(summary *Summary) (sync.Map, int) {
	m := sync.Map{}
	if summary.Id < 1 {
		log.Info("this is first backup")
		return m, 0
	}
	//	log.Infof("recent backup: %s", summary.Date)
	//
	//	//The most recent backup was completed on May 5.
	//	// Recent backups were processed on May 5th.
	//	rows, err := b.dbOrigin.Query("select path, size, mtime from bak_origin")
	//	checkErr(err)
	//
	//	var count = 0
	//	var path string
	//	var size int64
	//	var modTime string
	//	for rows.Next() {
	//		f := newFile("", 0, time.Now())
	//		err = rows.Scan(&path, &size, &modTime)
	//		checkErr(err)
	//		f.Path = path
	//		f.Size = size
	//		f.ModTime, _ = time.Parse(time.RFC3339, modTime)
	//		m.Store(path, f)
	//		count += 1
	//	}
	//	return m, count
	return m, 0
}

func (b *Backup) Start() error {
	defer b.Stop()
	if err := b.init(); err != nil {
		return err
	}

	lastSummary, err := b.getLastSummary()
	if err != nil {
		return err
	}
	if lastSummary == nil {
		if err := b.generateFirstBackupData(); err != nil {
			return err
		}
	}

	if err := b.startBackup(); err != nil {
		return err
	}

	return nil

	// Load last backup data
	//lastSummary, err := b.getLastSummary()
	//if err != nil {
	//	return err
	//}
	//spew.Dump(lastSummary)

	//originMap, originCount := b.getOriginMap(lastSummary)
	//spew.Dump(originMap)
	//spew.Dump(originCount)

	// Write initial data to database
	//newMap := sync.Map{}
	//if originCount < 1 || b.srcDir != lastSummary.SrcDir {
	//	b.summary.State = Running
	//	b.summary.Message = "collecting initialize data"
	//	log.Info(b.summary.Message)
	//
	//	err := filepath.Walk(b.srcDir, func(path string, f os.FileInfo, err error) error {
	//		if f.IsDir() {
	//			return nil
	//		}
	//
	//		if !f.Mode().IsRegular() {
	//			return nil
	//		}
	//
	//		fi := newFile(path, f.Size(), f.ModTime())
	//		newMap.Store(path, fi)
	//		b.summary.TotalCount += 1
	//		b.summary.TotalSize += uint64(f.Size())
	//		return nil
	//	})
	//	if err != nil {
	//		return err
	//	}
	//	os.RemoveAll(b.tempDir)
	//	b.summary.ReadingTime = time.Now()
	//	b.summary.ComparisonTime = b.summary.ReadingTime
	//
	//	log.Infof("writing initial data")
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
	return nil
}

func (b *Backup) Stop() error {
	if err := b.db.Close(); err != nil {
		log.Error(err)
	}

	if err := b.fileDb.Close(); err != nil {
		log.Error(err)
	}

	return nil
}

func (b *Backup) writeToDatabase(newMap sync.Map, originMap sync.Map) error {
	//log.Info("writing to database")
	//
	//rs, err := b.dbLogTx.Exec("insert into bak_summary(date,src_dir,dst_dir,state,total_size,total_count,backup_modified,backup_added,backup_deleted,backup_success,backup_failure,backup_size,execution_time,message) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
	//	b.summary.Date.Format(time.RFC3339),
	//	b.summary.SrcDir,
	//	b.summary.DstDir,
	//	b.summary.State,
	//	b.summary.TotalSize,
	//	b.summary.TotalCount,
	//	b.summary.BackupModified,
	//	b.summary.BackupAdded,
	//	b.summary.BackupDeleted,
	//	b.summary.BackupSuccess,
	//	b.summary.BackupFailure,
	//	b.summary.BackupSize,
	//	b.summary.ExecutionTime,
	//	b.summary.Message,
	//)
	//if err != nil {
	//	return err
	//}
	//
	//b.summary.Id, _ = rs.LastInsertId()
	//log.Infof("backup_id=%d", b.summary.ID)
	//
	//var maxInsertSize uint32 = 500
	//var lines []string
	//var eventLines []string
	//var i uint32 = 0
	//var j uint32 = 0
	//
	//// Delete original data
	//b.dbOriginTx.Exec("delete from bak_origin")
	//
	//// Modified or added files
	//newMap.Range(func(key, value interface{}) bool {
	//	f := value.(*File)
	//	path := strings.Replace(f.Path, "'", "''", -1)
	//	lines = append(lines, fmt.Sprintf("select '%s', %d, '%s'", path, f.Size, f.ModTime.Format(time.RFC3339)))
	//
	//	i += 1
	//
	//	if i%maxInsertSize == 0 || i == b.summary.TotalCount {
	//		err := b.insertIntoOrigin(lines)
	//		checkErr(err)
	//		lines = nil
	//	}
	//
	//	if f.State != 0 {
	//		eventLines = append(eventLines, fmt.Sprintf("select %d, '%s', %d, '%s', %d, '%s'", b.summary.ID, path, f.Size, f.ModTime.Format(time.RFC3339), f.State, f.Message))
	//		j += 1
	//
	//		if j%maxInsertSize == 0 {
	//			err := b.insertIntoLog(eventLines)
	//			checkErr(err)
	//			eventLines = nil
	//		}
	//	}
	//	return true
	//})
	//if len(eventLines) > 0 {
	//	err := b.insertIntoLog(eventLines)
	//	checkErr(err)
	//	eventLines = nil
	//}
	//
	//// Deleted files
	//eventLines = make([]string, 0)
	//j = 0
	//originMap.Range(func(key, value interface{}) bool {
	//	atomic.AddUint32(&b.summary.BackupSuccess, 1)
	//	f := value.(*File)
	//	log.Debugf("deleted: %s", f.Path)
	//	f.State = FileDeleted
	//	path := strings.Replace(f.Path, "'", "''", -1)
	//	eventLines = append(eventLines, fmt.Sprintf("select %d, '%s', %d, '%s', %d, '%s'", b.summary.ID, path, f.Size, f.ModTime.Format(time.RFC3339), f.State, f.Message))
	//	j += 1
	//
	//	if j%maxInsertSize == 0 {
	//		err := b.insertIntoLog(eventLines)
	//		checkErr(err)
	//		eventLines = nil
	//	}
	//	return true
	//})
	//if len(eventLines) > 0 {
	//	err := b.insertIntoLog(eventLines)
	//	checkErr(err)
	//	eventLines = nil
	//}
	//atomic.AddUint32(&b.summary.BackupDeleted, j)

	return nil
}

//func (b *Backup) insertIntoLog(rows []string) error {
//	query := fmt.Sprintf("insert into bak_log(id, path, size, mtime, state, message) %s", strings.Join(rows, " union all "))
//	_, err := b.dbLogTx.Exec(query)
//	return err
//}
//
//func (b *Backup) insertIntoOrigin(rows []string) error {
//	query := fmt.Sprintf("insert into bak_origin(path, size, mtime) %s", strings.Join(rows, " union all "))
//	_, err := b.dbOriginTx.Exec(query)
//	defer func() {
//		if r := recover(); r != nil {
//			if err != nil {
//				log.Println(query)
//			}
//		}
//	}()
//	checkErr(err)
//	return err
//}
//
//func (b *Backup) Close() error {
//	b.summary.ExecutionTime = b.summary.LoggingTime.Sub(b.summary.Date).Seconds()
//	b.summary.Message += fmt.Sprintf("reading: %3.1fs, comparing: %3.1fs, writing: %3.1fs",
//		b.summary.ReadingTime.Sub(b.summary.Date).Seconds(),
//		b.summary.ComparisonTime.Sub(b.summary.ReadingTime).Seconds(),
//		b.summary.LoggingTime.Sub(b.summary.ComparisonTime).Seconds(),
//	)
//	b.dbLogTx.Exec("update bak_summary set backup_deleted = ?, execution_time = ?, message = ? where id = ?",
//		b.summary.BackupDeleted,
//		b.summary.ExecutionTime,
//		b.summary.Message,
//		b.summary.ID,
//	)
//
//	b.dbLogTx.Commit()
//	b.dbOriginTx.Commit()
//	b.dbOrigin.Close()
//	b.dbLog.Close()
//
//	if b.summary.ID > 1 { // ID 1 is about initializing data
//		log.WithFields(log.Fields{
//			"modified": b.summary.BackupModified,
//			"added":    b.summary.BackupAdded,
//			"deleted":  b.summary.BackupDeleted,
//		}).Infof("files: %d", b.summary.BackupModified+b.summary.BackupAdded+b.summary.BackupDeleted)
//		log.WithFields(log.Fields{
//			"success": b.summary.BackupSuccess,
//			"failure": b.summary.BackupFailure,
//		}).Infof("backup result")
//		log.Infof("backup size: %d(%s)", b.summary.BackupSize, humanize.Bytes(b.summary.BackupSize))
//	}
//	log.WithFields(log.Fields{
//		"files": b.summary.TotalCount,
//		"size":  fmt.Sprintf("%d(%s)", b.summary.TotalSize, humanize.Bytes(b.summary.TotalSize)),
//	}).Info("source directory")
//
//	log.WithFields(log.Fields{
//		"reading":    fmt.Sprintf("%3.1fs", b.summary.ReadingTime.Sub(b.summary.Date).Seconds()),
//		"comparison": fmt.Sprintf("%3.1fs", b.summary.ComparisonTime.Sub(b.summary.ReadingTime).Seconds()),
//		"writing":    fmt.Sprintf("%3.1fs", b.summary.LoggingTime.Sub(b.summary.ComparisonTime).Seconds()),
//	}).Infof("execution time: %3.1fs", b.summary.ExecutionTime)
//
//	return nil
//}
//
//func (b *Backup) BackupFile(path string) (string, float64, error) {
//	// Set source
//	t := time.Now()
//	from, err := os.Open(path)
//	if err != nil {
//		return "", time.Since(t).Seconds(), err
//
//	}
//	defer from.Close()
//
//	// Set destination
//	dst := filepath.Join(b.tempDir, path[len(b.srcDir):])
//	err = os.MkdirAll(filepath.Dir(dst), 0644)
//	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
//	if err != nil {
//		return "", time.Since(t).Seconds(), err
//	}
//	defer to.Close()
//
//	// Copy
//	_, err = io.Copy(to, from)
//	if err != nil {
//		return "", time.Since(t).Seconds(), err
//	}
//
//	return dst, time.Since(t).Seconds(), err
//}
//
//func checkErr(err error) {
//	if err != nil {
//		log.Errorf("[Error] %s", err.Error())
//	}
//}

func (b *Backup) issueSummaryId() (int64, error) {
	var summaryId int64
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketSummary)
		if b == nil {
			return errors.New("invalid bucket: " + string(BucketSummary))
		}
		id, _ := b.NextSequence()
		summaryId = int64(id)
		return b.Put(Int64ToBytes(summaryId), nil)
	})
	return summaryId, err
}
