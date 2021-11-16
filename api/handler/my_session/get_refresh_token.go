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

type GetRTResponse struct {
	Value      string    `json:"value"`
	Type       string    `json:"type"`
	Expired_at time.Time `json:"expired_at"`
}

func GetRefreshToken(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helper.AuthenticateUserAccessToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		atJti, rtJti, err := jwt.GetAtJtiNRtJtiFromAccessToken(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		res, err := jwt.GenerateRefreshToken(userId, atJti, rtJti, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		response := map[string]GetRTResponse{
			"refresh_token": {
				Value:      res.RefreshToken,
				Expired_at: time.Unix(res.RtExpires, 0),
				Type:       "jwt",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
