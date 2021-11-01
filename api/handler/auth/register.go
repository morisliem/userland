package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"userland/api/helper"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"

	"github.com/gofrs/uuid"
)

type RegisterRequest struct {
	Fullname         string `json:"fullname"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Confirm_password string `json:"confirm_password"`
}

func Register(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RegisterRequest
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

		userId, err := uuid.NewV4()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		newRegister := store.User{
			Id:       userId.String(),
			Fullname: request.Fullname,
			Email:    request.Email,
			Password: hashPassword,
		}
		rn := helper.GenerateRandomNumber()

		err = userStore.RegisterUser(ctx, newRegister, rn)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		go helper.SendEmailVerCode(newRegister.Email, rn)

		err = tokenStore.SetEmailVerificationCode(newRegister.Id, rn)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		err = tokenStore.SetNewEmail(newRegister.Id, newRegister.Email)
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

func (rr *RegisterRequest) ValidateRequest() (map[string]string, error) {
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

	if len(res) > 0 {
		return res, errors.New("error")
	} else {
		return res, nil
	}
}
