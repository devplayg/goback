package goback

import (
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
	date   time.Time
	dstDir string
	tempDir   string
	backupDir string
}

func (k *LocalKeeper) Open(date time.Time, dstDir string) (string, string, error) {
	k.date = date
	k.dstDir = dstDir
	tempDir, err := ioutil.TempDir(k.dstDir, "backup-")
	if err != nil {
		return "", "", err
	}
	k.tempDir = tempDir
	k.backupDir = FindProperBackupDirName(filepath.Join(k.dstDir, k.date.Format("20060102")))
	return k.tempDir, k.backupDir, nil
}

func (k *LocalKeeper) Close() error {
	return os.Rename(k.tempDir, k.backupDir)
}
func (k *LocalKeeper) Test() error {
	return nil
}

// Copy file
func (k *LocalKeeper) Keep(srcPath, dstDir string) (string, float64, error) {
	// Set source
	t := time.Now()
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", 0.0, err

	}
	defer srcFile.Close()

	// Set destination
	if runtime.GOOS == "windows" {
		//  /BACKUP_DIR/C:/TEMP/DATA => error
		//  /BACKUP_DIR/C/TEMP/DATA => OK
		srcPath = strings.ReplaceAll(srcPath, ":", "")
	}
	dstPath := filepath.Join(dstDir, srcPath)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0644); err != nil {
		return "", 0.0, err
	}
	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return "", 0.0, err
	}
	defer dstFile.Close()

	// Copy
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", 0.0, err
	}

	return dstPath, time.Since(t).Seconds(), err
}
