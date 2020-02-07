package goback

type RemoteSite struct {
	Protocol int // FTP, SFTP
	Host     string
	Port     int
	Path     string
	Username string
	Password string
}

func newRemoteSite(protocol int, host string, port int, path, username, password string) *RemoteSite {
	return &RemoteSite{
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Path:     path,
		Username: username,
		Password: password,
	}
}
