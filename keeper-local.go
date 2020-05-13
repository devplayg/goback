package goback

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LocalKeeper saves added or modified files in local disk.
type LocalKeeper struct {
	*KeeperDesc
	started   time.Time
	dstDir    string
	tempDir   string
	backupDir string
}

func NewLocalKeeper(dstDir string) *LocalKeeper {
	return &LocalKeeper{
		dstDir: dstDir,
		KeeperDesc: &KeeperDesc{
			Protocol: LocalDisk,
			Host:     "local",
			Dir:      dstDir,
		},
	}
}

func (k *LocalKeeper) Init(t time.Time) error {
	if len(k.dstDir) < 1 {
		return errors.New("storage directory required")
	}

	if !DirExists(k.dstDir) {
		return fmt.Errorf("storage directory does not exist: %s", k.dstDir)
	}

	k.started = t
	tempDir, err := ioutil.TempDir(k.dstDir, "backup-")
	if err != nil {
		return err
	}
	k.tempDir = tempDir
	k.backupDir = FindProperBackupDirName(filepath.Join(k.dstDir, k.started.Format("20060102")))
	return nil
}

func (k *LocalKeeper) Active() bool {
	return true
}

func (k *LocalKeeper) Close() error {
	return os.Rename(k.tempDir, k.backupDir)
}

func (k *LocalKeeper) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (k *LocalKeeper) Description() *KeeperDesc {
	return k.KeeperDesc
}

// Copy file
func (k *LocalKeeper) keep(path string) (string, float64, error) {
	t := time.Now()

	// Set destination
	p := path
	if runtime.GOOS == "windows" {
		//  /BACKUP_DIR/C:/TEMP/DATA => error
		//  /BACKUP_DIR/C/TEMP/DATA => OK
		p = strings.ReplaceAll(path, ":", "")
	}
	dstPath := filepath.Join(k.tempDir, p)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0644); err != nil {
		return "", 0.0, err
	}
	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return "", 0.0, err
	}
	defer dstFile.Close()

	// Set source
	srcFile, err := os.Open(path)
	if err != nil {
		return "", 0.0, err

	}
	defer srcFile.Close()

	// Copy
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", 0.0, err
	}

	return dstPath, time.Since(t).Seconds(), err
}
