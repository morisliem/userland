package auth

import (
	"encoding/json"
	"net/http"
	"userland/store"
)

func Login(userStore store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Logged in")
	}
}
