package mydetail

import (
	"encoding/json"
	"net/http"
	"userland/api/helper"
	"userland/api/response"
	"userland/store"
)

type UpdateUserRequest struct {
	Fullname string `json:"fullname"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
	Web      string `json:"web"`
}

func UpdateUserDetail(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateUserRequest

		userId, err := helper.AuthenticateUser(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		newUpdate := store.User{
			Fullname: request.Fullname,
			Location: request.Location,
			Bio:      request.Bio,
			Web:      request.Web,
		}

		err = userStore.UpdateUserDetail(r.Context(), newUpdate, userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response.Success())
	}
}
