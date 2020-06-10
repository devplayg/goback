package goback

import (
	"encoding/hex"
	"github.com/devplayg/goback/tpl"
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
	http.Redirect(w, r, HomeUri, http.StatusSeeOther)
}

func (c *Controller) display(name string, tpl string, w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.New(name).Parse(himma.Base())
	if err != nil {
		return err
	}
	if tmpl, err = tmpl.Parse(tpl); err != nil {
		return err
	}
	if err := tmpl.Execute(w, c.app); err != nil {
		return err
	}
	return nil
}

func (c *Controller) DisplayLogin(w http.ResponseWriter, r *http.Request) {
	if isLogged(w, r) {
		http.Redirect(w, r, HomeUri, http.StatusSeeOther)
		return
	}
	if c.server.appConfig.DeveloperMode {
		if err := c.display("login", DisplayWithLocalFile("login"), w, r); err != nil {
			ResponseErr(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	if err := c.display("login", tpl.Login(), w, r); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplayBackup(w http.ResponseWriter, r *http.Request) {
	if c.server.appConfig.DeveloperMode {
		if err := c.display("login", DisplayWithLocalFile("backup"), w, r); err != nil {
			ResponseErr(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	if err := c.display("login", tpl.Backup(), w, r); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplayStats(w http.ResponseWriter, r *http.Request) {
	if c.server.appConfig.DeveloperMode {
		if err := c.display("login", DisplayWithLocalFile("stats"), w, r); err != nil {
			ResponseErr(w, r, err, http.StatusInternalServerError)
		}
		return
	}
	if err := c.display("login", tpl.Stats(), w, r); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
	}
}

func (c *Controller) DisplaySettings(w http.ResponseWriter, r *http.Request) {
	// config, err := loadConfig()
	// if err != nil {
	//	log.Error(err)
	//	ResponseErr(w, r, errors.New("failed to load settings"), http.StatusInternalServerError)
	// }

	// testTemplate, err = template.New("hello.gohtml").Funcs(template.FuncMap{
	// 	"hasPermission": func(feature string) bool {
	// 		return false
	// 	},
	// }).ParseFiles("hello.gohtml")

	tmpl, err := template.New("settings").Funcs(funcMap).Parse(himma.Base())
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
	}
	if c.server.appConfig.DeveloperMode {
		if tmpl, err = tmpl.Funcs(funcMap).Parse(DisplayWithLocalFile("settings")); err != nil {
			ResponseErr(w, r, err, http.StatusInternalServerError)
		}
	} else {
		if tmpl, err = tmpl.Funcs(funcMap).Parse(tpl.Settings()); err != nil {
			ResponseErr(w, r, err, http.StatusInternalServerError)
		}
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
