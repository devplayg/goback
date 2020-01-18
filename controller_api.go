package goback

import (
	"encoding/json"
	"github.com/devplayg/himma"
	"html/template"
	"net/http"
)

func (c *Controller) GetSummaries(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(c.summaries)
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (c *Controller) DisplayBackup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("streams")
	tmpl, err := tmpl.Parse(himma.Base())
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if tmpl, err = tmpl.Parse(DisplayBackup()); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}
