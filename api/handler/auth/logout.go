package auth

import (
	"encoding/json"
	"net/http"
	"userland/api/jwt"
	"userland/api/response"
	"userland/store"
)

func Logout(tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		au, err := jwt.ExtractTokenMetadata(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		deleted, err := jwt.DeleteAuth(au.AccessUuid, tokenStore)
		if err != nil || deleted == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response.Response("successfully logged out"))
	}
}
