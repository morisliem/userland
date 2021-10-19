package auth

import (
	"encoding/json"
	"net/http"
	"userland/store"
)

func Register(userStore store.UserStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_ = userStore.GetUser(ctx)
		success := struct{ Success bool `json:"success"`} { Success: true}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(success)
	}
}