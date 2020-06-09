package goback

import (
	"github.com/devplayg/himma/v2"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
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

	c.router.Use(c.authMiddleware)

	http.Handle("/", c.router)

	return nil
}

func (c *Controller) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// is static file ?
		if strings.HasPrefix(r.RequestURI, AssetUriPrefix) {
			next.ServeHTTP(w, r)
			return
		}

		if r.RequestURI == "/" {
			http.Redirect(w, r, LoginUri, http.StatusSeeOther)
			return
		}

		// Check session
		if !isLogged(w, r) && !strings.HasPrefix(r.RequestURI, LoginUri) {
			// http.Error(w, "Forbidden", http.StatusForbidden)
			tmpl, _ := template.New("error").Parse(himma.Base())
			tmpl.Parse(errorPageTpl())
			m := struct {
				*himma.Config
				Status     string
				StatusCode int
			}{
				c.app, "Forbidden", http.StatusForbidden,
			}
			if err := tmpl.Execute(w, m); err != nil {
				ResponseErr(w, r, err, http.StatusInternalServerError)
			}

			tmpl.Execute(w, m)
			return
		}

		next.ServeHTTP(w, r)
	})
}
