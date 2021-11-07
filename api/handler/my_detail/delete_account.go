package mydetail

import (
	"encoding/json"
	"net/http"
	"userland/api/helper"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type DeleteAccountRequest struct {
	Password string `json:"password"`
}

func DeleteUserAccount(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request DeleteAccountRequest
		userId, err := helper.AuthenticateUserAccessToken(r, tokenStore)
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

		res, err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

		userPwd, _ := userStore.GetPassword(r.Context(), userId)

		if !helper.ComparePasswordHash(request.Password, userPwd) {
			res := map[string]string{
				"message": "password incorrect",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

		err = userStore.DeleteAccount(r.Context(), userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.Success())
	}
}

func (dar *DeleteAccountRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}

	err := validator.ValidatePassword(dar.Password)
	if err != nil {
		res["password"] = err.Error()
	}

	if len(res) > 0 {
		return res, err
	} else {
		return res, nil
	}
}
