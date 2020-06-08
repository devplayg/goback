package goback

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (c *Controller) setRouter() error {
	c.router = mux.NewRouter()

	c.router.HandleFunc("/", c.DisplayDefault)

	// Assets
	c.router.HandleFunc("/assets/{u1}/{u2}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	c.router.HandleFunc("/assets/{u1}/{u2}/{u3}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	c.router.HandleFunc("/assets/{u1}/{u2}/{u3}/{u4}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})

	// System
	c.router.HandleFunc("/sysInfo", c.SysInfo).Methods(http.MethodGet)

	// Backup
	c.router.HandleFunc("/login", c.DisplayLogin).Methods(http.MethodGet)
	c.router.HandleFunc("/login", c.DoLogin).Methods(http.MethodPost)
	c.router.HandleFunc("/logout", c.DoLogout)
	c.router.HandleFunc("/backup/", c.DisplayBackup)
	c.router.HandleFunc("/summaries", c.GetSummaries)
	c.router.HandleFunc("/summaries/{id:[0-9]+}", c.GetSummary)
	c.router.HandleFunc("/stats", c.GetStats)
	c.router.HandleFunc("/summaries/{id:[0-9]+}/changes", c.GetChangesLog)

	// Statistics
	c.router.HandleFunc("/stats/", c.DisplayStats)

	// Settings
	c.router.HandleFunc("/settings/", c.DisplaySettings)
	c.router.HandleFunc("/settings/job/id/{id:[0-9]+}", c.UpdateJob).Methods(http.MethodPatch)
	c.router.HandleFunc("/settings/storage/id/{id:[0-9]+}", c.UpdateStorage).Methods(http.MethodPatch)
	c.router.HandleFunc("/backup/{id:[0-9]+}/run", c.RunBackupJob).Methods(http.MethodGet)

	c.router.Use(authMiddleware)

	http.Handle("/", c.router)

	return nil
}
