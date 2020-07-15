package goback

import (
	"encoding/hex"
	"github.com/devplayg/goback/tpl"
	"github.com/devplayg/himma/v2"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"time"
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

func (c *Controller) display(name string, tpl string, w http.ResponseWriter, data interface{}) error {
	tmpl, err := template.New(name).Parse(himma.Base())
	if err != nil {
		return err
	}
	if tmpl, err = tmpl.Parse(tpl); err != nil {
		return err
	}
	if err := tmpl.Execute(w, data); err != nil {
		return err
	}
	return nil
}

func (c *Controller) DisplayLogin(w http.ResponseWriter, r *http.Request) {
	if isLogged(w, r) {
		http.Redirect(w, r, HomeUri, http.StatusSeeOther)
		return
	}
	if err := c.display("login", tpl.Login(), w, c.app); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (c *Controller) DisplayNewAccessKey(w http.ResponseWriter, r *http.Request) {
	_, _, err := c.server.getAccessKeyAndSecretKey()
	if err == nil {
		if isLogged(w, r) {
			http.Redirect(w, r, HomeUri, http.StatusSeeOther)
			return
		}
	}

	if err := c.display("login", tpl.NewAccount(), w, c.app); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (c *Controller) DisplayBackup(w http.ResponseWriter, r *http.Request) {
	if err := c.display("backup", tpl.Backup(), w, c.app); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (c *Controller) DisplayStats(w http.ResponseWriter, r *http.Request) {
	if err := c.display("stats", tpl.Stats(), w, c.app); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (c *Controller) DisplayStatsReport(w http.ResponseWriter, r *http.Request) {
	t, _ := time.Parse("200601", mux.Vars(r)["yyyymm"])
	m := struct {
		*himma.Config
		YYYYMM string
		From   string
		To     string
	}{
		c.app,
		t.Format("Jan, 2006"),
		t.Format("Jan 2"),
		t.AddDate(0, 1, -1).Format("Jan 2, 2006"),
	}

	if err := c.display("report", tpl.Report(), w, m); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (c *Controller) DisplaySettings(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("settings").Funcs(funcMap).Parse(himma.Base())
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	if tmpl, err = tmpl.Funcs(funcMap).Parse(tpl.Settings()); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	checksum, _ := c.server.getDbValue(ConfigBucket, KeyConfigChecksum)
	m := struct {
		*himma.Config
		Settings *Config
		Checksum string
	}{
		c.app, c.server.config, hex.EncodeToString(checksum),
	}
	if err := tmpl.Execute(w, m); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
	}
}
