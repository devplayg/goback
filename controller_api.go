package goback

import "net/http"

func (c *Controller) GetSummaries(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("asdfasdf"))
}
