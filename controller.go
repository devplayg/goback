package goback

import (
	"context"
	"encoding/json"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/himma"
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
	app       *himma.Application
	summaries []*Summary
	server    *Server
}

func NewController(server *Server, dir, addr string, app *himma.Application) *Controller {
	return &Controller{
		server:    server,
		addr:      addr,
		dir:       dir,
		summaries: make([]*Summary, 0),
		version:   app.Version,
		app:       app,
	}
}

func (c *Controller) init() error {
	if err := c.initRouter(); err != nil {
		return err
	}
	uiAssetMap, err := himma.GetAssetMap(
		himma.Bootstrap4,
		himma.BootstrapDatepicker_1_9_0,
		himma.BootstrapSelect_1_13_9,
		himma.BootstrapTable_1_15_5,
		himma.Holder_2_9,
		himma.JqueryMask_1_14_16,
		himma.JqueryValidation_1_19_1,
		himma.JsCookie_2_2_1,
		himma.Moment_2_24_0,
		himma.ProgressbarJs_1_0_1,
		himma.VideoJs_7_7_4,
		himma.WaitMe_31_10_17,
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
	c.router.HandleFunc("/assets/{u1}/{u2}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	c.router.HandleFunc("/assets/{u1}/{u2}/{u3}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	c.router.HandleFunc("/assets/{u1}}/{u2}/{u3}/{u4}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	http.Handle("/", c.router)

	//srv := &http.Server{
	//	Handler:      c.router,
	//	Addr:         c.addr,
	//	WriteTimeout: 15 * time.Second,
	//	ReadTimeout:  15 * time.Second,
	//}
	//go func() {
	//	log.WithFields(log.Fields{}).Infof("listen on %s", c.addr)
	//	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
	//		log.Fatalf("HTTP server ListenAndServe: %v", err)
	//	}
	//
	//}()

	return nil
}

func (c *Controller) Start() error {
	if err := c.init(); err != nil {
		return err
	}
	defer c.Stop()

	srv := &http.Server{
		Handler:      c.router,
		Addr:         c.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	ch := make(chan struct{})
	go func() {
		<-c.server.Ctx.Done()
		c.server.Log.Info("2) got stop signal from engine")
		if err := srv.Shutdown(context.Background()); err != nil {
			c.server.Log.Error(err)
		}
		close(ch)
	}()

	c.server.Log.Debug("1) HTTP server has been started")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		c.server.Log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-ch
	c.server.Log.Debug("3) HTTP server has been stopped")
	return nil

	// ch := make(chan struct{})
	// go func() {
	// 	<-c.server.Ctx.Done()
	// 	c.server.Log.Info("2) got stop sig")
	// 	// Error from closing listeners, or context timeout:
	// 	//log.Printf("HTTP server Shutdown: %v", err)
	// 	//}
	// 	//close(idleConnsClosed)
	// 	if err := srv.Shutdown(context.Background()); err != nil {
	// 		// Error from closing listeners, or context timeout:
	// 		c.server.Log.Error("HTTP server Shutdown: %v", err)
	// 	}
	// 	close(ch)
	//
	// }()
	//
	// c.server.Log.Infof("1) listen on %s", c.addr)
	// if err := srv.ListenAndServe(); err != http.ErrServerClosed {
	// 	log.Fatalf("HTTP server ListenAndServe: %v", err)
	// }
	// <-ch
	// c.server.Log.Info("3) http server has been stopped")
	// return nil
}

func (c *Controller) Stop() error {
	return nil
}

func (c *Controller) loadSummaryDb() error {
	path := filepath.Join(c.dir, SummaryDbName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.server.Log.Warnf("database not found: %s", path)
		return nil
	}
	var summaries []*Summary
	if err := LoadBackupData(path, &summaries, GobEncoding); err != nil {
		return err
	}
	c.summaries = summaries
	// spew.Dump(summaries)
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
