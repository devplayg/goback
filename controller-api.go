package goback

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
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

func (c *Controller) GetStatsReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := time.Parse("200601", vars["yyyymm"])
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
	stats, err := c.server.findMonthlySummaries(vars["yyyymm"])
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

	serverAcessKey, serverSecretKey, err := c.server.getAccessKeyAndSecretKey()
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}
	m := map[string]interface{}{
		"logged": false,
		"url":    "",
	}

	ak := sha256.Sum256([]byte(signIn.AccessKey))
	sk := sha256.Sum256([]byte(signIn.SecretKey))
	if !(bytes.Equal(ak[:], serverAcessKey) && bytes.Equal(sk[:], serverSecretKey)) {
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

func (c *Controller) CreateNewAccount(w http.ResponseWriter, r *http.Request) {
	var signIn SignIn
	_, err := c.ParseForm(r, &signIn)
	if err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	ak := sha256.Sum256([]byte(signIn.AccessKey))
	sk := sha256.Sum256([]byte(signIn.SecretKey))
	if err := c.server.setAccessKeyAndSecretKey(ak[:], sk[:]); err != nil {
		ResponseErr(w, r, err, http.StatusInternalServerError)
		return
	}

	log.Info("ACCESS KEY is created")

	http.Redirect(w, r, HomeUri, http.StatusMovedPermanently)
}
