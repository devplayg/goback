package goback

import (
    "crypto/md5"
    "encoding/hex"
    "errors"
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "github.com/devplayg/golibs/compress"
    "github.com/devplayg/golibs/converter"
    "io/ioutil"
    "os"
    "path/filepath"
)

func (c *Controller) getChangesLog(id int) (map[string]*StatsReportWithList, error) {
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

    b, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to load db: %w", err)
    }

    decompressed, err := compress.Decompress(b, compress.GZIP)
    if err != nil {
        return nil, fmt.Errorf("failed to decompress db: %w", err)
    }
    var changesLog map[string]*StatsReportWithList
    //log.Debug(len(decompressed))
    err = converter.DecodeFromBytes(decompressed, &changesLog)
    if err != nil {
        spew.Dump(decompressed)
        return nil, fmt.Errorf("failed to decode decompressed data : %w", err)
    }
    return changesLog, nil
}
