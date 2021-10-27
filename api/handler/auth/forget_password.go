package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"userland/api/helper"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type ForgetPasswordRequest struct {
	Email string `json:"email"`
}

func ForgetPassword(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ForgetPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		newRequest := store.User{
			Email: request.Email,
		}

		err = userStore.EmailExist(r.Context(), newRequest)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response("email doens't exist"))
			return
		}

		rn := helper.GenerateRandomNumber()
		go helper.SendEmailResetPwdCode(newRequest.Email, rn)

		err = tokenStore.SetEmailVerificationCode(newRequest.Email, rn)
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

func (fpr *ForgetPasswordRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}

	err := validator.ValidateEmail(fpr.Email)
	if err != nil {
		res["message"] = err.Error()
	}

	if len(res) > 0 {
		return res, errors.New("error")
	} else {
		return res, nil
	}
}
