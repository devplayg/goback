package goback

import (
	"net/http"
	"strings"
)

type SignIn struct {
	AccessKey string
	SecretKey string
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// is static file ?
		if strings.HasPrefix(r.RequestURI, AssetUriPrefix) {
			next.ServeHTTP(w, r)
			return
		}

		if r.RequestURI == "/" {
			http.Redirect(w, r, LoginUri, http.StatusSeeOther)
			return
		}

		// Check session
		if !isLogged(w, r) && !strings.HasPrefix(r.RequestURI, LoginUri) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isLogged(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		return false
	}

	if len(session.Values) < 1 {
		return false
	}

	return true
}
