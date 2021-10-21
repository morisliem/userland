package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"userland/api/helper"
	"userland/api/jwt"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type ResetPasswordRequest struct {
	Token            string
	Password         string `json:"password"`
	Confirm_password string `json:"confirm_password"`
}

func ResetPassword(userStore store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ResetPasswordRequest
		ctx := r.Context()

		userId, err := jwt.TokenValid(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
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

		hashPassword, err := helper.HashPassword(request.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		newResetPasswordRequest := store.User{
			Password: hashPassword,
		}

		err = userStore.ResetPassword(ctx, userId, newResetPasswordRequest)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response.Success())
	}
}

func (rpr *ResetPasswordRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}
	errPwd := validator.ValidatePassword(rpr.Password)

	if errPwd != nil {
		res["password"] = errPwd.Error()
	}

	errCPwd := validator.ValidatePassword(rpr.Confirm_password)

	if errCPwd != nil {
		res["confirm_password"] = errCPwd.Error()
	}

	if rpr.Password != rpr.Confirm_password {
		res["unmatch_password"] = "check again your password"
	}

	if len(res) > 1 {
		return res, errors.New("Error")
	} else {
		return res, nil
	}

}
