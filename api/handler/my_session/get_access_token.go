package mysession

import (
	"encoding/json"
	"net/http"
	"time"
	"userland/api/helper"
	"userland/api/jwt"
	"userland/api/response"
	"userland/store"
)

type GetATResponse struct {
	Value      string    `json:"value"`
	Type       string    `json:"type"`
	Expired_at time.Time `json:"expired_at"`
}

func GetAccessToken(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helper.AuthenticateUserRefreshToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		atJti, err := jwt.GetAtJtiFromRefreshToken(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		deleted, err := jwt.DeleteATAuth(atJti, tokenStore)
		if err != nil || deleted == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		deleted, err = jwt.DeleteRTAuth(atJti, tokenStore)
		if err != nil || deleted == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := jwt.GenerateAccessToken(userId, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		err = userStore.UpdateUserSession(r.Context(), atJti, res.AccessUuid)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := map[string]GetATResponse{
			"access_token": {
				Value:      res.AccessToken,
				Expired_at: time.Unix(res.AtExpires, 0),
				Type:       "jwt",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
