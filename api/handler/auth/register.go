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
		err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
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

		// // Check if the email is in the database
		// err = userStore.EmailExist(ctx, newRegister)
		// if err != nil {
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(response.Response(err.Error()))
		// 	return
		// }

		// err = helper.SendEmail(newRegister.Email, 123)
		// if err != nil {
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	json.NewEncoder(w).Encode(response.Response(err.Error()))
		// 	return
		// }

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

func (rr *RequestRequest) ValidateRequest() error {
	if msg, err := validator.ValidateFullname(rr.Fullname); !err {
		return errors.New(msg)
	}

	if msg, err := validator.ValidateEmail(rr.Email); !err {
		return errors.New(msg)
	}

	if msg, err := validator.ValidatePassword(rr.Password); !err {
		return errors.New(msg)
	}

	if rr.Password != rr.Confirm_password {
		return errors.New("password not match")
	}

	return nil
}
