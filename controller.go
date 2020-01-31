package goback

import (
	"encoding/json"
	"fmt"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/himma"
	"github.com/devplayg/himma/asset"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	router    *mux.Router
	summaryDb *os.File
	addr      string
	dir       string
	version   string

	summaries []*Summary
}

func NewController(dir, addr, version string) *Controller {
	return &Controller{
		addr:      addr,
		dir:       dir,
		summaries: make([]*Summary, 0),
		version:   version,
	}
}

func (c *Controller) init() error {
	if err := c.initRouter(); err != nil {
		return err
	}
	uiAssetMap, err := himma.GetAssetMap(
		asset.Bootstrap4,
		asset.BootstrapDatepicker_1_9_0,
		asset.BootstrapSelect_1_13_9,
		asset.BootstrapTable_1_15_5,
		asset.Holder_2_9,
		asset.JqueryMask_1_14_16,
		asset.JqueryValidation_1_19_1,
		asset.JsCookie_2_2_1,
		asset.Moment_2_24_0,
		asset.ProgressbarJs_1_0_1,
		asset.VideoJs_7_7_4,
		asset.WaitMe_31_10_17,
	)
	if err != nil {
		return err
	}
	WebAssetMap = uiAssetMap

	if err := c.loadSummaryDb(); err != nil {
		return err
	}
	return nil
}

func (c *Controller) initRouter() error {
	c.router = mux.NewRouter()
	c.router.HandleFunc("/", c.DisplayBackup)

	c.router.HandleFunc("/summaries", c.GetSummaries)
	c.router.HandleFunc("/summaries/{id:[0-9]+}/changes", c.GetChangesLog)

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
	return nil
}

func (c *Controller) loadSummaryDb() error {
	path := filepath.Join(c.dir, SummaryDbName)
	var summaries []*Summary
	if err := LoadBackupData(path, &summaries, GobEncoding); err != nil {
		return err
	}
	c.summaries = summaries
	//spew.Dump(summaries)
	return nil
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
