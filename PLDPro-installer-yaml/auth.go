package main

import (
	"net/http"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || !checkAuth(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+cfg.Settings.Realm+`"`)
			http.Error(w, "Acesso não autorizado", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func checkAuth(username, password string) bool {
	expectedUser := cfg.Settings.AuthUserEnv
	expectedPass := cfg.Settings.AuthPassEnv
	return username == expectedUser && password == expectedPass
}
