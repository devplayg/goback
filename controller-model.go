package goback

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (c *Controller) getChangesLog(id int) ([]byte, error) {
	summary := c.findSummaryById(id)
	if summary == nil {
		return nil, errors.New("summary not found")
	}

	h := md5.Sum([]byte(summary.SrcDir))
	suffix := hex.EncodeToString(h[:])
	var logPath string
	if summary.Version == 1 {
		key := fmt.Sprintf("%s-%d", summary.Date.Format("20060102"), summary.BackupId)
		logPath = filepath.Join(c.dir, key, "changes-"+suffix+".db")
	} else {
		// key = fmt.Sprintf("%s-%d", summary.Date.Format("20060102"), summary.BackupId)
		logPath = filepath.Join(summary.BackupDir, "changes-"+suffix+".db")
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil, err
	}

	return ioutil.ReadFile(logPath)
}
