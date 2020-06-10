package goback

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devplayg/himma/v2"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	WebAssetMap himma.AssetMap
)

type Controller struct {
	router    *mux.Router
	summaryDb *os.File
	addr      string
	dbDir     string
	version   string
	app       *himma.Config
	summaries []*Summary
	server    *Server
}

func NewController(server *Server, app *himma.Config) *Controller {
	return &Controller{
		server:    server,
		addr:      server.appConfig.Address,
		dbDir:     server.dbDir,
		summaries: make([]*Summary, 0),
		version:   app.Version,
		app:       app,
	}
}

func (c *Controller) init() error {
	if err := c.setRouter(); err != nil {
		return err
	}
	uiAssetMap, err := himma.GetAssetMap(
		himma.Bootstrap4,
		himma.BootstrapDatepicker_1_9_0,
		himma.BootstrapSelect_1_13_9,
		himma.BootstrapTable_1_16_0,
		himma.Holder_2_9,
		himma.JqueryMask_1_14_16,
		himma.JqueryValidation_1_19_1,
		himma.JsCookie_2_2_1,
		himma.Moment_2_24_0,
		himma.ProgressbarJs_1_0_1,
		himma.WaitMe_31_10_17,
	)
	if err != nil {
		return err
	}
	uiAssetMap.AddZippedAndBase64Encoded("/assets/img/logo.png", LogoImg)
	uiAssetMap.AddRaw("/assets/js/custom.js", customScript())
	uiAssetMap.AddRaw("/assets/css/custom.css", customCss())
	WebAssetMap = uiAssetMap

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
		if err := srv.Shutdown(context.Background()); err != nil {
			c.server.Log.Error(err)
		}
		close(ch)
	}()

	// c.server.Log.Debug("1) HTTP server has been started")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println(err.Error())
		c.server.Log.Error(err)
	}
	<-ch
	// c.server.Log.Debug("3) HTTP server has been stopped")
	return nil
}

func (c *Controller) Stop() error {
	return nil
}

func GetAsset(w http.ResponseWriter, r *http.Request) {
	if content, hasAsset := WebAssetMap[r.RequestURI]; hasAsset {
		w.Header().Set("Content-Type", DetectContentType(filepath.Ext(r.RequestURI)))
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(content)), 10))
		w.Header().Set("Content-Encoding", GZIP)
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

func ResponseErr(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	if statusCode != http.StatusOK {
		log.WithFields(logrus.Fields{
			"ip":     r.RemoteAddr,
			"uri":    r.RequestURI,
			"method": r.Method,
			"length": r.ContentLength,
		}).Error(err)
	}
	w.Header().Add("Content-Type", ApplicationJson)
	b, _ := json.Marshal(map[string]interface{}{
		"error": err.Error(),
	})
	w.WriteHeader(statusCode)
	w.Write(b)
}

func ResponseData(w http.ResponseWriter, r *http.Request, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", ApplicationJson)
	w.Write(json)
}

func ResponseZippedRaw(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Add("Content-Type", ApplicationJson)
	w.Header().Set("Content-Encoding", GZIP)
	w.Write(data)
}
