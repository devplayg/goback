package goback

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/himma"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var (
	uiAssetMap map[string][]byte
)

type Controller struct {
	b      *Backup
	router *mux.Router
	db     *bolt.DB
	fileDb *bolt.DB
	addr   string
	dir    string
}

func NewController(dir, addr string) *Controller {
	return &Controller{
		addr: addr,
		dir:  dir,
	}
}

func (c *Controller) init() error {
	if err := c.initRouter(); err != nil {
		return err
	}
	db, fileDb, err := InitDatabase(c.dir)
	if err != nil {
		return err
	}
	c.db, c.fileDb = db, fileDb

	uiAssetMap, err = himma.GetAssetMap(himma.Bootstrap4)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) initRouter() error {
	c.router = mux.NewRouter()
	c.router.HandleFunc("/summaries", c.GetSummaries)
	c.router.HandleFunc("/assets/{assetType}/{name}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})

	c.router.HandleFunc("/assets/plugins/{pluginName}/{name}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	c.router.HandleFunc("/assets/plugins/{pluginName}/{kind}/{name}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})

	c.router.HandleFunc("/assets/modules/{moduleName}/{name}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	http.Handle("/", c.router)

	srv := &http.Server{
		Handler:      c.router,
		Addr:         c.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
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
	if err := c.db.Close(); err != nil {
		log.Error(err)
	}
	if err := c.fileDb.Close(); err != nil {
		log.Error(err)
	}
	return nil
}

func GetAsset(w http.ResponseWriter, r *http.Request) {
	if content, hasAsset := uiAssetMap[r.RequestURI]; hasAsset {
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
