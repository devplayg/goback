package goback

type RemoteSite interface {
	Open() error
	Close() error
	Send(srcPath, dstPath string) error
}

type FtpSite struct {
	Protocol int // FTP, SFTP
	Host     string
	Port     int
	Path     string
	Username string
	Password string
}

func newFtpSite(protocol int, host string, port int, path, username, password string) *FtpSite {
	return &FtpSite{
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Path:     path,
		Username: username,
		Password: password,
	}
}
