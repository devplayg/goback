package goback

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Keeper interface {
	// Connected() (bool, error)
	Keep(string, string) (string, float64, error)
	// Get()
}

func NewKeeper(dst string) Keeper {
	d := strings.ToLower(dst)
	if strings.HasPrefix(d, "ftp://") {
		return nil
	}
	if strings.HasPrefix(d, "sftp://") {
		return nil
	}
	return LocalKeeper{}
}

// host,port, user, pass
// LocalKeeper.Connected/Save
// FtpSender.Connected/Send
// SftpSender.OK/Send
type LocalKeeper struct {}

func (k LocalKeeper) Keep(srcPath, dstDir string) (string, float64, error) {
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

type FtpKeeper struct{}

type SftpKeeper struct{}



// func BackupFile(srcPath, destDir string) (string, float64, error) {
// 	// Set source
// 	t := time.Now()
// 	srcFile, err := os.Open(srcPath)
// 	if err != nil {
// 		return "", 0.0, err
//
// 	}
// 	defer srcFile.Close()
//
// 	// Set destination
// 	if runtime.GOOS == "windows" {
// 		//  /BACKUP_DIR/C:/TEMP/DATA => error
// 		//  /BACKUP_DIR/C/TEMP/DATA => OK
// 		srcPath = strings.ReplaceAll(srcPath, ":", "")
// 	}
// 	dstPath := filepath.Join(destDir, srcPath)
// 	if err := os.MkdirAll(filepath.Dir(dstPath), 0644); err != nil {
// 		return "", 0.0, err
// 	}
// 	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, 0600)
// 	if err != nil {
// 		return "", 0.0, err
// 	}
// 	defer dstFile.Close()
//
// 	// Copy
// 	_, err = io.Copy(dstFile, srcFile)
// 	if err != nil {
// 		return "", 0.0, err
// 	}
//
// 	return dstPath, time.Since(t).Seconds(), err
// }
