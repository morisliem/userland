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

type ResetPasswordRequest struct {
	Email            string `json:"email"`
	Code             int    `json:"code"`
	Password         string `json:"password"`
	Confirm_password string `json:"confirm_password"`
}

func ResetPassword(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ResetPasswordRequest
		ctx := r.Context()

		err := json.NewDecoder(r.Body).Decode(&request)
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

		code, err := tokenStore.GetEmailVarificationCode(request.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		if code != request.Code {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response("invalid code"))
			return
		}

		userId, err := userStore.GetUserid(ctx, request.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		listOfPwd, _ := userStore.GetPasswords(ctx, userId)

		// checking if user has used the password before
		for i := 0; i < len(listOfPwd); i++ {
			if helper.ComparePasswordHash(request.Password, listOfPwd[i]) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response.Bad_request("You have used this password before, try to use other password"))
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

		newResetPasswordRequest := store.User{
			Password: hashPassword,
		}

		err = userStore.UpdatePassword(ctx, userId, newResetPasswordRequest)
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

	errEmail := validator.ValidateEmail(rpr.Email)
	if errEmail != nil {
		res["email"] = errEmail.Error()
	}

	errCode := validator.ValidateCode(rpr.Code)
	if errCode != nil {
		res["code"] = errCode.Error()
	}

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

	if len(res) > 0 {
		return res, errors.New("Error")
	} else {
		return res, nil
	}
}
