package goback

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Controller struct {
	b      *Backup
	r      *mux.Router
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
	return nil
}

func (c *Controller) initRouter() error {
	c.r = mux.NewRouter()
	c.r.HandleFunc("/summaries", c.GetSummaries)

	srv := &http.Server{
		Handler:      c.r,
		Addr:         c.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//go func() {
	log.Fatal(srv.ListenAndServe())
	//}()

	return nil
}

func (c *Controller) Start() error {
	if err := c.init(); err != nil {
		return err
	}

	defer c.Stop()

	srv := &http.Server{
		Handler:      c.r,
		Addr:         c.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	//http.Handle("/", c.r)

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

func (c *Controller) GetSummaries(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("asdfasdf"))
}
