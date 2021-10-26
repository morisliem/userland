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

type RequestRequest struct {
	Fullname         string `json:"fullname"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Confirm_password string `json:"confirm_password"`
}

func Register(userStore store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RequestRequest
		ctx := r.Context()
		json.NewDecoder(r.Body).Decode(&request)

		// Validate the request
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
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		// Storing the request to user struct
		newRegister := store.User{
			Fullname: request.Fullname,
			Email:    request.Email,
			Password: hashPassword,
		}

		err = userStore.RegisterUser(ctx, newRegister)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

func (rr *RequestRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}
	nameErr := validator.ValidateFullname(rr.Fullname)
	if nameErr != nil {
		res["fullname"] = nameErr.Error()
	}

	emailErr := validator.ValidateEmail(rr.Email)
	if emailErr != nil {
		res["email"] = emailErr.Error()
	}

	pwdErr := validator.ValidatePassword(rr.Password)
	if pwdErr != nil {
		res["password"] = pwdErr.Error()
	}

	if rr.Password != rr.Confirm_password {
		res["unmatch_password"] = "check again your password"
	}

	if len(res) > 1 {
		return res, errors.New("Error")
	} else {
		return res, nil
	}
}
