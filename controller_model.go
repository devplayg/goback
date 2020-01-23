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
        return nil, errors.New("summry not found")
    }
    key := fmt.Sprintf("%s-%d", summary.Date.Format("20060102"), summary.BackupId)
    h := md5.Sum([]byte(summary.SrcDir))
    suffix := hex.EncodeToString(h[:])
    path := filepath.Join(c.dir, key, "changes-"+suffix+".db")
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, err
    }

    return ioutil.ReadFile(path)
}
