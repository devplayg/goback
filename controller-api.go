package goback

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (c *Controller) SysInfo(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{
		"time": time.Now(),
	}
	ResponseData(w, r, m)
}

func (c *Controller) GetSummaries(w http.ResponseWriter, r *http.Request) {
	summaries, err := c.server.findSummaries()
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseData(w, r, summaries)
}

func (c *Controller) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := c.server.findStats()
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseData(w, r, stats)
}

func (c *Controller) GetChangesLog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	data, err := c.server.getChangesLog(id)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseZippedRaw(w, r, data)
}

func (c *Controller) GetSummary(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	summary, err := c.server.findSummaryById(id)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseData(w, r, summary)
}

func (c *Controller) DoLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, LoginUri, http.StatusSeeOther)
}

// DoLogin
func (c *Controller) DoLogin(w http.ResponseWriter, r *http.Request) {
	var signIn SignIn
	_, err := c.ParseForm(r, &signIn)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
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
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	ResponseData(w, r, m)
}
