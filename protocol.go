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

// LocalKeeper saves added or modified files in local disk.
type LocalKeeper struct{}

// Copy file
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

// FtpKeeper saves added or modified files to remote server via FTP
type FtpKeeper struct {
	addr     string
	username string
	password string
}

func NewFtpKeeper(addr, username, password string) *FtpKeeper {
	return &FtpKeeper{
		addr:     addr,
		username: username,
		password: password,
	}
}

func (k FtpKeeper) Keep(srcPath, dstDir string) (string, float64, error) {
	return "", 0, nil
}

// SftpKeeper saves added or modified files to remote server via SFTP
type SftpKeeper struct {
	addr     string
	username string
	password string
}

func NewSftpKeeper(addr, username, password string) *SftpKeeper {
	return &SftpKeeper{
		addr:     addr,
		username: username,
		password: password,
	}
}

func (k SftpKeeper) Keep(srcPath, dstDir string) (string, float64, error) {
	return "", 0, nil
}

// host,port, user, pass
// LocalKeeper.Connected/Save
// FtpSender.Connected/Send
// SftpSender.OK/Send

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
