package goback

import (
	"encoding/hex"
	"github.com/devplayg/goback/tpl"
	"github.com/devplayg/himma/v2"
	"html/template"
	"net/http"
	"strings"
)

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{
		"DirExists": DirExists,
	}
}

func (c *Controller) DisplayDefault(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, HomeUri, http.StatusSeeOther)
}

type SignIn struct {
	AccessKey string
	SecretKey string
}

//
// func checkAuth( w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r,  SignInSessionName)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	spew.Dump(session)
// }

func authMiddleware(next http.Handler) http.Handler {
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
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isLogged(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		return false
	}

	if len(session.Values) < 1 {
		return false
	}

	return true
}

func (c *Controller) DisplayLogin(w http.ResponseWriter, r *http.Request) {
	if isLogged(w, r) {
		http.Redirect(w, r, HomeUri, http.StatusSeeOther)
		return
	}
	tmpl, err := template.New("login").Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	// if tmpl, err = tmpl.Parse(DisplayWithLocalFile("backup")); err != nil {
	if tmpl, err = tmpl.Parse(tpl.Login()); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, c.app); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplayBackup(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("backup").Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	// if tmpl, err = tmpl.Parse(DisplayWithLocalFile("backup")); err != nil {
	if tmpl, err = tmpl.Parse(tpl.Backup()); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, c.app); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplayStats(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("stats").Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	// if tmpl, err = tmpl.Parse(DisplayWithLocalFile("stats")); err != nil {
	if tmpl, err = tmpl.Parse(tpl.Stats()); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, c.app); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplaySettings(w http.ResponseWriter, r *http.Request) {
	// config, err := loadConfig()
	// if err != nil {
	//	log.Error(err)
	//	Response(w, r, errors.New("failed to load settings"), http.StatusInternalServerError)
	// }

	// testTemplate, err = template.New("hello.gohtml").Funcs(template.FuncMap{
	// 	"hasPermission": func(feature string) bool {
	// 		return false
	// 	},
	// }).ParseFiles("hello.gohtml")

	tmpl, err := template.New("settings").Funcs(funcMap).Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	// if tmpl, err = tmpl.Parse(DisplayWithLocalFile("settings")); err != nil {
	// 	Response(w, r, err, http.StatusInternalServerError)
	// }
	// if tmpl, err = tmpl.Funcs(funcMap).Parse(DisplayWithLocalFile("settings")); err != nil {
	if tmpl, err = tmpl.Funcs(funcMap).Parse(tpl.Settings()); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}

	checksum, _ := c.server.getDbValue(ConfigBucket, KeyConfigChecksum)
	aa := struct {
		*himma.Config
		Settings *Config
		Checksum string
	}{
		c.app, c.server.config, hex.EncodeToString(checksum),
	}
	if err := tmpl.Execute(w, aa); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}
