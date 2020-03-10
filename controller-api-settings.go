package goback

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (c *Controller) UpdateJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//id, _ := strconv.Atoi(vars["id"])
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
