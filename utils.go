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
    "sort"
    "strings"
    "sync"
    "time"
)

const (
    DefaultDateFormat = "2006-01-02 15:04:05"
)

var ErrorBucketNotFound = errors.New("bucket not found")

const (
    kB = 1000
    MB = 1000000
    GB = 1000000000
)

var fileSizeCategories = []int64{
    1 * kB,
    5 * kB,
    10 * kB,
    50 * kB,
    100 * kB,
    500 * kB,

    1 * MB,
    5 * MB,
    10 * MB,
    50 * MB,
    100 * MB,
    500 * MB,

    1 * GB,
    5 * GB,
    10 * GB,
    50 * GB,
    100 * GB,
    500 * GB,
}

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

func GetFileMap(dir string, hashComparision bool) (*sync.Map, map[string]int64, map[int64]int64, int64, uint64, error) {
    fileMap := sync.Map{}
    extensionMap := make(map[string]int64)
    sizeDistribution := make(map[int64]int64)

    var size uint64
    var count int64

    err := filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
        if file.IsDir() {
            return nil
        }

        if !file.Mode().IsRegular() {
            return nil
        }

        fi := NewFileWrapper(path, file.Size(), file.ModTime())
        if hashComparision {
            h, err := GetFileHash(path)
            if err != nil {
                return err
            }
            fi.Hash = h
        }

        // Statistics
        ext := strings.ToLower(filepath.Ext(file.Name()))
        if len(ext) > 0 {
            extensionMap[ext]++
        } else {
            extensionMap["__NO_EXT__"]++
        }
        sizeDistribution[GetFileSizeCategory(file.Size())]++
        size += uint64(fi.Size)
        count++

        fileMap.Store(path, fi)
        return nil
    })

    return &fileMap, extensionMap, sizeDistribution, count, size, err
}

func GetFileSizeCategory(size int64) int64 {
    for i := range fileSizeCategories {
        if size <= fileSizeCategories[i] {
            return fileSizeCategories[i]
        }
    }
    return -1

}

func IsEqualStringSlices(a, b []string) bool {
    sort.Strings(a)
    sort.Strings(b)
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

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

// func InitDatabase(summaryDbPath, fileMapDbPath string) (*os.File, *os.File, error) {
// 	summaryDb, err := os.OpenFile(summaryDbPath, os.O_RDWR|os.O_CREATE, 0644)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	fileMapDb, err := os.OpenFile(fileMapDbPath, os.O_RDWR|os.O_CREATE, 0644)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	return summaryDb, fileMapDb, nil
// }

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

func CreateSrcDirsHashMap(dirs []string) map[string]string {
    m := make(map[string]string)
    for _, u := range dirs {
        FillMd5ValueOfStringKey(m, u)
    }
    return m
}
