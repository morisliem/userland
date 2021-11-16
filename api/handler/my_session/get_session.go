package mysession

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"userland/api/helper"
	"userland/api/jwt"
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
		var clientInfo ClientInfo
		var SessionResponse GetSessionResponse
		userId, err := helper.AuthenticateUserAccessToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		atJti, err := jwt.GetAtJtiFromAccessToken(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		res, err := userStore.GetUserSession(r.Context(), userId, atJti)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		SessionResponse.Created_at = res.Created_at
		SessionResponse.Ip = res.Ip
		SessionResponse.Is_current = res.Is_current
		SessionResponse.Updated_at = res.Updated_at

		for _, v := range res.Client {
			clientInfo.Name = v.Name
			clientInfo.SessionId = v.SessionId
			SessionResponse.Client = append(SessionResponse.Client, clientInfo)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SessionResponse)
	}
}
