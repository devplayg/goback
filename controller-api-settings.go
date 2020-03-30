package goback

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

func (c *Controller) ParseForm(r *http.Request, o interface{}) (map[string]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, errors.New("parse error")
	}
	if err := schema.NewDecoder().Decode(o, r.PostForm); err != nil {
		return nil, errors.New("parse error")
	}
	return mux.Vars(r), nil
}

func (c *Controller) UpdateJob(w http.ResponseWriter, r *http.Request) {
	var input Job
	vars, err := c.ParseForm(r, &input)
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
		return
	}
	id, _ := strconv.Atoi(vars["id"])

	for i, job := range c.server.config.Jobs {
		if job.Id == id {
			c.server.config.Jobs[i] = input
			break
		}
	}

	if err := c.server.saveConfig(); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(vars["id"]))
}

func (c *Controller) UpdateStorage(w http.ResponseWriter, r *http.Request) {
	var input Storage
	vars, err := c.ParseForm(r, &input)
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
		return
	}
	id, _ := strconv.Atoi(vars["id"])

	for i, job := range c.server.config.Storages {
		if job.Id == id {
			if id == 1 {
				input.Protocol = LocalDisk
			} else if job.Id == 2 {
				input.Protocol = Sftp
			}
			c.server.config.Storages[i] = input
			break
		}
	}

	if err := c.server.saveConfig(); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(vars["id"]))
}

func (c *Controller) RunBackupJob(w http.ResponseWriter, r *http.Request) {
	var input Job
	vars, err := c.ParseForm(r, &input)
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
		return
	}
	jobId, _ := strconv.Atoi(vars["id"])

	if err := c.server.runBackupJob(jobId); err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}

	w.Write([]byte(vars["id"]))
}
