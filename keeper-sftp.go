package goback

import (
	"fmt"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SftpKeeper struct {
	*KeeperDesc
	host      string
	port      int
	dstDir    string
	username  string
	password  string
	conn      *sftp.Client
	active    bool
	date      time.Time
	backupDir string
}

func NewSftpKeeper(storage Storage) *SftpKeeper {
	if storage.Port < 1 {
		storage.Port = 22 // default SSH port
	}
	return &SftpKeeper{
		host:     storage.Host,
		port:     storage.Port,
		dstDir:   storage.Dir,
		username: storage.Username,
		password: storage.Password,
		conn:     nil,
		active:   false,
		KeeperDesc: &KeeperDesc{
			Protocol: Sftp,
			Host:     storage.Host,
			Dir:      storage.Dir,
		},
	}
}

func (k *SftpKeeper) Init(t time.Time) error {
	k.date = t
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}
	auths = append(auths, ssh.Password(k.password))
	config := ssh.ClientConfig{
		User:            k.username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", k.host, k.port), &config)
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP server: %w", err)
	}

	size := 1 << 15
	ftpConn, err := sftp.NewClient(conn, sftp.MaxPacket(size))
	if err != nil {
		return fmt.Errorf("failed to create to SFTP client: %w", err)
	}

	k.conn = ftpConn
	k.active = true
	k.backupDir = k.FindProperBackupDirName(filepath.Join(k.dstDir, k.date.Format("20060102")))
	return nil
}

func (k *SftpKeeper) Active() bool {
	return k.active
}

func (k *SftpKeeper) Close() error {
	return k.conn.Close()
}

func (k *SftpKeeper) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil
}

func (k *SftpKeeper) Description() *KeeperDesc {
	return k.KeeperDesc
}

func (k *SftpKeeper) FindProperBackupDirName(dir string) string {
	dir = filepath.ToSlash(dir)
	i := 0
	for {
		var d string
		if i < 1 {
			d = dir
		} else {
			d = dir + "-" + strconv.Itoa(i)
		}
		if _, err := k.conn.Stat(d); os.IsNotExist(err) {
			return d
		}
		i++
	}
}

func (k *SftpKeeper) keep(path string) (string, float64, error) {
	// time.Sleep(5 * time.Second)
	t := time.Now()

	p := path
	if runtime.GOOS == "windows" {
		//  /BACKUP_DIR/C:/TEMP/DATA => error
		//  /BACKUP_DIR/C/TEMP/DATA => OK
		p = strings.ReplaceAll(path, ":", "")
	}
	dstPath := filepath.ToSlash(filepath.Join(k.backupDir, p))
	dstDir := filepath.ToSlash(filepath.Dir(dstPath))
	if err := k.conn.MkdirAll(dstDir); err != nil {
		return dstPath, 0, err
	}
	dstFile, err := k.conn.OpenFile(dstPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
	if err != nil {
		return "", 0, err
	}
	defer dstFile.Close()

	localFile, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer localFile.Close()
	fi, err := localFile.Stat()
	if err != nil {
		return "", 0, err
	}
	size := fi.Size()
	// log.Debugf("writing %v bytes", size)
	t1 := time.Now()
	n, err := io.Copy(dstFile, io.LimitReader(localFile, size))
	if err != nil {
		return "", 0, err
	}
	if n != size {
		return "", 0, fmt.Errorf("copy: expected %v bytes, got %d", size, n)
	}
	log.Debugf("wrote %v bytes in %s", size, time.Since(t1))
	return dstPath, time.Since(t).Seconds(), nil
}

// func (f *FtpSite) Open() error {
// 	var auths []ssh.AuthMethod
// 	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
// 		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
// 	}
// 	auths = append(auths, ssh.Password(f.Password))
// 	config := ssh.ClientConfig{
// 		User:            f.Username,
// 		Auth:            auths,
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}
// 	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", f.Host, f.Port), &config)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to SFTP server: %w", err)
// 	}
//
// 	size := 1 << 15
// 	ftpConn, err := sftp.NewClient(conn, sftp.MaxPacket(size))
// 	if err != nil {
// 		return fmt.Errorf("failed to create to SFTP client: %w", err)
// 	}
//
// 	f.conn = ftpConn
// 	return nil
// }
//
// func (f *FtpSite) Test() error {
// 	f.conn.Walk(f.Dir)
// 	return nil
// }
//
// func (f FtpSite) Close() error {
// 	if err := f.conn.Close(); err != nil {
// 		return fmt.Errorf("failed to close SFTP connection: %w", err)
// 	}
// 	return nil
// }
//
// func (f FtpSite) Send(src, dst string) error {
// 	d := filepath.ToSlash(filepath.Join("/backup", filepath.Base(src)))
// 	w, err := f.conn.OpenFile(d, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
// 	if err != nil {
// 		return err
// 	}
// 	defer w.Close()
//
// 	file, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	fi, err := file.Stat()
// 	if err != nil {
// 		return err
// 	}
// 	size := fi.Size()
// 	log.Debugf("writing %v bytes", size)
// 	t1 := time.Now()
// 	n, err := io.Copy(w, io.LimitReader(file, size))
// 	if err != nil {
// 		return err
// 	}
// 	if n != size {
// 		return fmt.Errorf("copy: expected %v bytes, got %d", size, n)
// 	}
// 	log.Debug("wrote %v bytes in %s", size, time.Since(t1))
// 	return nil
// }
//
// func (b *Backup) sendChangedFiles() error {
// 	if err := b.ftpSite.Open(); err != nil {
// 		return fmt.Errorf("failed to open remote ftp: %w", err)
// 	}
// 	defer b.ftpSite.Close()
// 	fileGroup, _, err := b.createBackupFileGroup()
// 	if err != nil {
// 		return err
// 	}
//
// 	for i := range fileGroup {
// 		for j := range fileGroup[i] {
// 			fileWrapper := fileGroup[i][j]
// 			// spew.Dump(fileGroup[i][j])
// 			// send := filepath.Join(fileWrapper.dir, fileWrapper)
// 			dst := filepath.Join(b.dstDir)
// 			if err := b.ftpSite.Send(fileWrapper.Path, filepath.Base(fileWrapper.Path)); err != nil {
// 				log.Error(err)
// 				continue
// 			}
// 		}
// 	}
//
// 	// spew.Dump(fileGroup)
//
// 	return nil
// }
