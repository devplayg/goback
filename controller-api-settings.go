package goback

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (c *Controller) ParseForm(r *http.Request, o interface{}) (map[string]string, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	if err := schema.NewDecoder().Decode(o, r.PostForm); err != nil {
		return nil, err
	}
	return mux.Vars(r), nil
}

func (c *Controller) UpdateJob(w http.ResponseWriter, r *http.Request) {

	// Parse form
	var input Job
	vars, err := c.ParseForm(r, &input)
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
		return
	}
	jobId, _ := strconv.Atoi(vars["id"])
	input.Tune()

	// Find job
	job := c.server.findJobById(jobId)
	if job == nil {
		Response(w, r, errors.New("job not found"), http.StatusInternalServerError)
		return
	}
	if job.running {
		Response(w, r, errors.New("job is running"), http.StatusInternalServerError)
		return
	}

	// Parse cron
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional)
	_, err = cronParser.Parse(input.Schedule)
	if err != nil {
		Response(w, r, fmt.Errorf("failed to parse scheduler; %w", err), http.StatusInternalServerError)
		return
	}

	// Update job
	job.Schedule = input.Schedule
	job.SrcDirs = input.SrcDirs
	job.Enabled = input.Enabled
	job.BackupType = input.BackupType

	c.server.rwMutex.Lock()
	defer func() {
		c.server.rwMutex.Unlock()
	}()
	if job.cronEntryId != nil {
		log.WithFields(logrus.Fields{
			"jobId":       job.Id,
			"cronEntryId": job.cronEntryId,
		}).Info("SCHEDULE REMOVED")
		c.server.cron.Remove(*job.cronEntryId)
		job.cronEntryId = nil
	}
	if job.Enabled {
		entryId, err := c.server.cron.AddFunc(job.Schedule, func() {
			log.WithFields(logrus.Fields{
				"jobId": job.Id,
			}).Info("RUN SCHEDULER")
			if err := c.server.runBackupJob(jobId); err != nil {
				log.Error(err)
			}
		})
		if err != nil {
			Response(w, r, err, http.StatusInternalServerError)
			return
		}
		job.cronEntryId = &entryId
		log.WithFields(logrus.Fields{
			"schedule":    job.Schedule,
			"jobId":       job.Id,
			"cronEntryId": entryId,
		}).Info("SCHEDULE RESERVED")
	}

	if err := c.server.saveConfig(input.Checksum); err != nil {
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

	storageId, _ := strconv.Atoi(vars["id"])
	input.Tune()

	storage := c.server.findStorageById(storageId)
	if storage == nil {
		Response(w, r, errors.New("storage not found"), http.StatusInternalServerError)
		return
	}
	storage.Id = input.Id
	storage.Dir = input.Dir
	storage.Username = input.Username
	storage.Password = input.Password
	storage.Port = input.Port
	storage.Host = input.Host

	if err := c.server.saveConfig(input.Checksum); err != nil {
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
		return
	}

	w.Write([]byte(vars["id"]))
}
