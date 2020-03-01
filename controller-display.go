package goback

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/devplayg/himma/v2"
	"html/template"
	"net/http"
)

func (c *Controller) DisplayDefault(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/backup/", http.StatusSeeOther)
}

func (c *Controller) DisplayBackup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("backup")
	tmpl, err := tmpl.Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if tmpl, err = tmpl.Parse(DisplayWithLocalFile("backup")); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, c.app); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplaySettings(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfig()
	if err != nil {
		log.Error(err)
		Response(w, r, errors.New("failed to load settings"), http.StatusInternalServerError)
	}

	tpl := template.New("settings")
	tmpl, err := tpl.Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if tmpl, err = tmpl.Parse(DisplayWithLocalFile("settings")); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}

	aa := struct {
		*himma.Config
		Settings *Config
	}{
		c.app, config,
	}
	spew.Dump(aa)
	if err := tmpl.Execute(w, aa); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}
