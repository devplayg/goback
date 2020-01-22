package goback

import (
    "crypto/md5"
    "encoding/hex"
    "errors"
    "fmt"
    "github.com/dustin/go-humanize"
    "github.com/minio/highwayhash"
    "io"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "time"
)

var ErrorBucketNotFound = errors.New("bucket not found")

func IsValidDir(dir string) (string, error) {
    absDir, err := filepath.Abs(dir)
    if err != nil {
        return "", err
    }

    if _, err := os.Stat(absDir); os.IsNotExist(err) {
        return "", errors.New("directory not found: " + absDir)
    }

    fi, err := os.Lstat(absDir)
    if err != nil {
        return "", err
    }

    if !fi.Mode().IsDir() {
        return "", errors.New("invalid source directory: " + fi.Name())
    }

    return absDir, nil
}

func GetFileHash(path string) (string, error) {
    highwayhash, err := highwayhash.New(HashKey[:])
    if err != nil {
        return "", err
    }
    file, err := os.Open(path) // specify your file here
    if err != nil {
        return "", err
    }
    defer file.Close()

    if _, err = io.Copy(highwayhash, file); err != nil {
        return "", err
    }

    checksum := highwayhash.Sum(nil)
    return hex.EncodeToString(checksum), nil
}

//
//func GetFileMap(dir string, hashComparision bool) (*sync.Map, *FilesReport, int64, uint64, error) {
//    fileMap := sync.Map{}
//    report := NewFilesReport()
//
//    var size uint64
//    var count int64
//
//    err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
//        if file.IsDir() {
//            return nil
//        }
//
//        if !file.Mode().IsRegular() {
//            return nil
//        }
//
//        fi := NewFileWrapper(path, file.Size(), file.ModTime())
//        if hashComparision {
//            h, err := GetFileHash(path)
//            if err != nil {
//                return err
//            }
//            fi.Hash = h
//        }
//
//        // Statistics
//        report.addExtension(file.Name(), file.Size())
//        report.addSize(file)
//        size += uint64(fi.Size)
//        count++
//
//        fileMap.Store(path, fi)
//        return nil
//    })
//
//    return &fileMap, report, count, size, err
//}
//
//func GetFileSizeCategory(size int64) int64 {
//    for i := range fileSizeCategories {
//        if size <= fileSizeCategories[i] {
//            return fileSizeCategories[i]
//        }
//    }
//    return -1
//
//}

func BackupFile(srcPath, tempDir string) (string, float64, error) {
    // Set source
    t := time.Now()
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return "", 0.0, err

    }
    defer srcFile.Close()

    // Set destination
    if runtime.GOOS == "windows" {
        srcPath = strings.ReplaceAll(srcPath, ":", "")
    }
    dstPath := filepath.Join(tempDir, srcPath)
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

func GetHumanizedSize(size uint64) string {
    humanized := humanize.Bytes(size)

    str := fmt.Sprintf("%d B", size)
    if humanized == str {
        return str
    }
    return fmt.Sprintf("%s (%s)", str, humanized)
}

func LoadOrCreateDatabase(path string) (*os.File, error) {
    db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return nil, err
    }
    return db, nil
}

func FillMd5ValueOfStringKey(m map[string]string, key string) {
    k := strings.TrimSpace(key)
    sum := md5.Sum([]byte(key))
    v := hex.EncodeToString(sum[:])
    m[k] = v
}



func NewSizeDistribution() map[int64]int64 {
    m := make(map[int64]int64)
    for _, size := range fileSizeCategories {
        m[size] = 0
    }
    return m
}