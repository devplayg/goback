package goback

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
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
	//vars := mux.Vars(r)
	var job Job
	vars, err := c.ParseForm(r, &job)
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}
	jobId, _ := strconv.Atoi(vars["id"])

	//spew.Dump(job)
	// r.PostForm is a map of our POST form values
	//err := decoder.Decode(&person, r.PostForm)
	//if err != nil {
	// Handle error
	//}

	//
	//decoder := json.NewDecoder(r.Body)
	//
	//var job Job
	//err := decoder.Decode(job)
	//if err != nil {
	//	panic(err)
	//}

	//if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
	//	log.Error(err)
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//spew.Dump(job)

	//var job Job
	//_ = json.NewDecoder(r.G).Decode(&job)
	//spew.Dump(job)
	//body, _ := ioutil.ReadAll(r.Body)
	//id := r.FormValue("id")
	//spew.Dump(id)

	//owner := r.Form.Get("id")
	//
	//spew.Dump(jobId)
	//spew.Dump(owner)
	//spew.Dump(r)

	//decoder := json.NewDecoder(r.Body)
	//var job Job
	//err := decoder.Decode(&job)
	//if err != nil {
	//	panic(err)
	//}

	//spew.Dump(job)

	for _, j := range c.server.config.Jobs {
		if j.Id == jobId {
			spew.Dump(job)
			spew.Dump(j)
		}
	}

	//data, err := c.getChangesLog(id)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//w.Header().Set("Content-Encoding", compress.GZIP)
	//w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(vars["id"]))
	c.server.config.Save()
}
