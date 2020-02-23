package goback

import (
	"encoding/json"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/himma/v2"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

func (c *Controller) GetSummaries(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(c.summaries)
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (c *Controller) GetChangesLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	data, err := c.getChangesLog(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Encoding", compress.GZIP)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (c *Controller) findSummaryById(id int) *Summary {
	if c.summaries == nil || len(c.summaries) < 1 {
		return nil
	}
	if len(c.summaries) >= id {
		if c.summaries[id-1].Id == id {
			return c.summaries[id-1]
		}
	}
	for _, s := range c.summaries {
		if s.Id == id {
			return s
		}
	}

	return nil
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
	if err := tmpl.Execute(w, c.app); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
}
