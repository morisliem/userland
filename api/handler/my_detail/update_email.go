package mydetail

import (
	"encoding/json"
	"errors"
	"net/http"
	"userland/api/helper"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

func UpdateUserEmail(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateEmailRequest

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

		res, err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

		newEmail := store.User{
			Email: request.Email,
		}

		err = userStore.UpdateUserEmail(r.Context(), newEmail, userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		rn := helper.GenerateRandomNumber()
		go helper.SendEmailVerCode(newEmail.Email, rn)

		err = tokenStore.SetEmailVerificationCode(newEmail.Email, rn)
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

func (uer *UpdateEmailRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}
	err := validator.ValidateEmail(uer.Email)
	if err != nil {
		res["email"] = err.Error()
	}

	if len(res) > 0 {
		return res, errors.New("error")
	} else {
		return res, nil
	}
}
