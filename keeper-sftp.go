package goback

// type SftpKeeper struct {
// 	Protocol int // FTP, SFTP
// 	Host     string
// 	Port     int
// 	Dir      string
// 	Username string
// 	password string
// 	conn     *sftp.Client
// }
//
// func newFtpSite(protocol int, host string, port int, dir, username, password string) *SftpKeeper {
// 	return &SftpKeeper{
// 		Protocol: protocol,
// 		Host:     host,
// 		Port:     port,
// 		Dir:      dir,
// 		Username: username,
// 		password: password,
// 	}
// }

// func (k *SftpKeeper) Open() error {
// 	log.Debug("sftp open")
// 	return nil
// }
// func (k *SftpKeeper) Close() error {
// 	log.Debug("sftp close")
// 	return nil
// }
// func (k *SftpKeeper) Test() error {
// 	log.Debug("sftp test")
// 	return nil
// }
// func (k *SftpKeeper) Keep(srcPath, dstDir string) (string, float64, error) {
// 	log.Debug("sftp keep")
// 	return "", 0, nil
// }

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
