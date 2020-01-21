package goback

import (
    "encoding/json"
    "fmt"
    "github.com/devplayg/golibs/compress"
    "github.com/devplayg/golibs/converter"
	"github.com/devplayg/himma"
	"github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
    "mime"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

var (
    WebAssetMap map[string][]byte
)

type Controller struct {
    router *mux.Router
    //summaryDbPath string
    summaryDb *os.File
    //fileMapDbPath string
    //fileMapDb     *os.File
    addr string
    dir  string

    summaries []*Summary
}

func NewController(dir, addr string) *Controller {
    return &Controller{
        addr:      addr,
        dir:       dir,
        summaries: make([]*Summary, 0),
    }
}

func (c *Controller) init() error {
    if err := c.initRouter(); err != nil {
        return err
    }
    //c.summaryDbPath = filepath.Join(c.dir, SummaryDbName)
    //c.fileMapDbPath = filepath.Join(c.dir, FileMapDbName)
    //summaryDb, fileMapDb, err := InitDatabase(c.summaryDbPath, c.fileMapDbPath)
    // if err != nil {
    // 	return err
    // }
    // c.summaryDb, c.fileMapDb = summaryDb, fileMapDb
    //
    uiAssetMap, err := himma.GetAssetMap(himma.Bootstrap4, himma.Plugins)
    if err != nil {
    	return err
    }
	WebAssetMap = uiAssetMap

    summaries, err := c.loadSummaryDb()
    if err != nil {
        return err
    }
    c.summaries = summaries
    return nil
}

func (c *Controller) initRouter() error {
    c.router = mux.NewRouter()
    c.router.HandleFunc("/", c.DisplayBackup)
    c.router.HandleFunc("/summaries", c.GetSummaries)
    c.router.HandleFunc("/assets/{assetType}/{name}", func(w http.ResponseWriter, r *http.Request) {
        GetAsset(w, r)
    })

    c.router.HandleFunc("/assets/{u1}/{u2}/{u3}", func(w http.ResponseWriter, r *http.Request) {
        GetAsset(w, r)
    })
    c.router.HandleFunc("/assets/{u1}}/{u2}/{u3}/{u4}", func(w http.ResponseWriter, r *http.Request) {
        GetAsset(w, r)
    })
    // c.router.HandleFunc("/assets/plugins/{pluginName}/{kind}/{name}", func(w http.ResponseWriter, r *http.Request) {
    // 	GetAsset(w, r)
    // })
    //
    // c.router.HandleFunc("/assets/modules/{moduleName}/{name}", func(w http.ResponseWriter, r *http.Request) {
    // 	GetAsset(w, r)
    // })
    http.Handle("/", c.router)

    srv := &http.Server{
        Handler:      c.router,
        Addr:         c.addr,
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    go func() {
        log.WithFields(log.Fields{}).Infof("listen on %s", c.addr)
        log.Fatal(srv.ListenAndServe())
    }()

    return nil
}

func (c *Controller) Start() error {
    if err := c.init(); err != nil {
        return err
    }

    defer c.Stop()

    fmt.Scanln()
    return nil
}

func (c *Controller) Stop() error {
    //if err := c.summaryDb.Close(); err != nil {
    //    log.Error(err)
    //}
    return nil
}

func (c *Controller) loadSummaryDb() ([]*Summary, error) {
    data, err := ioutil.ReadFile(filepath.Join(c.dir, SummaryDbName))
    if err != nil {
        return nil, err
    }

    var summaries []*Summary
    if err := converter.DecodeFromBytes(data, &summaries); err != nil {
        return nil, fmt.Errorf("failed to load summary database: %w", err)
    }
    return summaries, nil
}

func GetAsset(w http.ResponseWriter, r *http.Request) {
    if content, hasAsset := WebAssetMap[r.RequestURI]; hasAsset {
        w.Header().Set("Content-Type", DetectContentType(filepath.Ext(r.RequestURI)))
        w.Header().Set("Content-Length", strconv.FormatInt(int64(len(content)), 10))
        w.Header().Set("Content-Encoding", compress.GZIP)
        w.Write(content)
    }
}

func DetectContentType(ext string) string {
    ctype := mime.TypeByExtension(filepath.Ext(ext))
    if ctype == "" {
        return "application/octet-stream"
    }
    return ctype
}

func Response(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
    if statusCode != http.StatusOK {
        log.WithFields(log.Fields{
            "ip":     r.RemoteAddr,
            "uri":    r.RequestURI,
            "method": r.Method,
            "length": r.ContentLength,
        }).Error(err)
    }
    w.Header().Add("Content-Type", "application/json")
    b, _ := json.Marshal(err)
    w.WriteHeader(statusCode)
    w.Write(b)
}
