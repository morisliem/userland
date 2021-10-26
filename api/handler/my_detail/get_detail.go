package mydetail

import (
	"encoding/json"
	"net/http"
	"time"
	"userland/api/helper"
	"userland/api/response"
	"userland/store"
)

type UserDetailResponse struct {
	Id         string    `json:"id"`
	Fullname   string    `json:"fullname"`
	Location   string    `json:"location"`
	Bio        string    `json:"bio"`
	Web        string    `json:"web"`
	Picture    string    `json:"picture"`
	Created_at time.Time `json:"created_at"`
}

func GetUserDetail(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, err := helper.AuthenticateUser(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		res, err := userStore.GetUserDetail(r.Context(), userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		userDetail := &UserDetailResponse{
			Id:         res.Id,
			Fullname:   res.Fullname,
			Location:   res.Location,
			Bio:        res.Bio,
			Picture:    res.Picture,
			Created_at: res.Created_at,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userDetail)
	}
}
