package goback

import (
	"github.com/gorilla/sessions"
	"os"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte(os.Getenv(SignInSessionId)))
	store.MaxAge(0)
}
