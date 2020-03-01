package goback

import (
	"encoding/json"
	"github.com/devplayg/golibs/compress"
	"github.com/gorilla/mux"
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
