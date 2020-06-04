package goback

import (
	"encoding/hex"
	"github.com/devplayg/himma/v2"
	"html/template"
	"net/http"
)

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{
		"DirExists": DirExists,
	}
}

func (c *Controller) DisplayDefault(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/backup/", http.StatusSeeOther)
}

func (c *Controller) DisplayBackup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("backup")
	tmpl, err := tmpl.Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	//if tmpl, err = tmpl.Parse(DisplayWithLocalFile("backup")); err != nil {
	if tmpl, err = tmpl.Parse(displayBackup()); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, c.app); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplayStats(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("stats")
	tmpl, err := tmpl.Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	//if tmpl, err = tmpl.Parse(DisplayWithLocalFile("stats")); err != nil {
	if tmpl, err = tmpl.Parse(displayStats()); err != nil {
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

	tpl := template.New("settings").Funcs(funcMap)
	tmpl, err := tpl.Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	// if tmpl, err = tmpl.Parse(DisplayWithLocalFile("settings")); err != nil {
	// 	Response(w, r, err, http.StatusInternalServerError)
	// }
	//if tmpl, err = tmpl.Funcs(funcMap).Parse(DisplayWithLocalFile("settings")); err != nil {
	if tmpl, err = tmpl.Funcs(funcMap).Parse(displaySettings()); err != nil {
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
