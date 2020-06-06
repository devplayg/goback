package goback

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

func (c *Controller) GetSummaries(w http.ResponseWriter, r *http.Request) {
	summaries, err := c.server.findSummaries()
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}

	b, _ := json.Marshal(summaries)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (c *Controller) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := c.server.findStats()
	if err != nil {
		Response(w, r, err, http.StatusInternalServerError)
	}

	b, _ := json.Marshal(stats)
	w.Header().Add("Content-Type", ApplicationJson)
	w.Write(b)
}

func (c *Controller) GetChangesLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	data, err := c.server.getChangesLog(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Encoding", GZIP)
	w.Header().Add("Content-Type", ApplicationJson)
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

func (c *Controller) DoLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, LoginUri, http.StatusSeeOther)
}

func (c *Controller) DoLogin(w http.ResponseWriter, r *http.Request) {
	var signIn SignIn
	_, err := c.ParseForm(r, &signIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// inputPwdHash := sha256.Sum256([]byte(SecretKey))
	serverPwdHash := sha256.Sum256([]byte(os.Getenv(SecretKey)))

	m := map[string]interface{}{
		"logged": false,
		"url":    "",
	}
	if !(os.Getenv(AccessKey) == signIn.AccessKey && signIn.SecretKey == hex.EncodeToString(serverPwdHash[:])) {
		b, _ := json.Marshal(m)
		w.Header().Add("Content-Type", ApplicationJson)
		w.Write(b)
		return
	}
	m["logged"] = true
	m["url"] = HomeUri

	session.Values[AccessKey] = r.RemoteAddr
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(m)
	w.Header().Add("Content-Type", ApplicationJson)
	w.Write(b)
	//http.Redirect(w, r, HomeUri, http.StatusSeeOther)
}
