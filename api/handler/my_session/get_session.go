package mysession

import (
	"encoding/json"
	"net/http"
	"time"
	"userland/api/helper"
	"userland/api/response"
	"userland/store"
)

type ClientInfo struct {
	SessionId string `json:"id"`
	Name      string `json:"name"`
}

type GetSessionResponse struct {
	Is_current bool         `json:"is_current"`
	Ip         string       `json:"ip"`
	Client     []ClientInfo `json:"clients"`
	Created_at time.Time    `json:"created_at"`
	Updated_at time.Time    `json:"updated_at"`
}

func GetUserSession(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helper.AuthenticateUserAccessToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		res, err := userStore.GetUserSession(r.Context(), userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
