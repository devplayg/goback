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
	key := hex.EncodeToString(h[:])
	logPath := filepath.Join(c.dbDir, fmt.Sprintf(ChangesDbName, key, summary.BackupId))
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil, err
	}

	return ioutil.ReadFile(logPath)
}
