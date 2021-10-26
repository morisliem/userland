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

type ChangePasswordRequest struct {
	Current_Password string `json:"current_password"`
	Password         string `json:"password"`
	Confirm_Password string `json:"confirm_password"`
}

func ChangeUserPassword(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ChangePasswordRequest
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

		listOfPwd, _ := userStore.GetPasswords(r.Context(), userId)
		reverse(listOfPwd)

		if !helper.ComparePasswordHash(request.Current_Password, listOfPwd[0]) {
			res := map[string]string{
				"message": "password incorrect",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

		loopTo := 0
		if len(listOfPwd) > 3 {
			loopTo = 2
		} else {
			loopTo = len(listOfPwd) - 1
		}

		for i := 0; i < loopTo; i++ {
			if helper.ComparePasswordHash(request.Password, listOfPwd[i]) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnprocessableEntity)
				errMsg := map[string]string{}
				errMsg["password"] = "You have used this password before, try to use other password"
				json.NewEncoder(w).Encode(response.UnproccesableEntity(errMsg))
				return
			}
		}

		hashPassword, err := helper.HashPassword(request.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		newPassword := store.User{
			Password: hashPassword,
		}

		err = userStore.UpdatePassword(r.Context(), userId, newPassword)
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

func (cpr *ChangePasswordRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}

	err := validator.ValidatePassword(cpr.Current_Password)
	if err != nil {
		res["current_password"] = err.Error()
	}

	err = validator.ValidatePassword(cpr.Password)
	if err != nil {
		res["password"] = err.Error()
	}

	err = validator.ValidatePassword(cpr.Confirm_Password)
	if err != nil {
		res["confirm_password"] = err.Error()
	}

	if cpr.Password != cpr.Confirm_Password {
		res["unmatch_password"] = "check again your password"
	}

	if len(res) > 0 {
		return res, errors.New("Error")
	} else {
		return res, nil
	}
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
