package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type ValidateEmailCodeRequest struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}

/*
	Updating the email once the email is verified
*/

func ValidateEmail(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ValidateEmailCodeRequest
		ctx := r.Context()
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		res, err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

		newValidateEmail := store.User{
			Email:   request.Email,
			VerCode: request.Code,
		}

		err = userStore.EmailExist(ctx, newValidateEmail)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		uid, err := userStore.GetUserid(ctx, newValidateEmail.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		code, err := tokenStore.GetEmailVarificationCode(uid)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		if code != newValidateEmail.VerCode {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			tmp := map[string]string{}
			tmp["code"] = "invalid code"
			json.NewEncoder(w).Encode(response.UnproccesableEntity(tmp))
			return
		}

		email, err := tokenStore.GetNewEmail(uid)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		newValidateEmail.Email = email
		newValidateEmail.Id = uid

		err = userStore.ValidateCode(ctx, newValidateEmail)
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

func (vr *ValidateEmailCodeRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}

	emailErr := validator.ValidateEmail(vr.Email)
	if emailErr != nil {
		res["email"] = emailErr.Error()
	}

	if len(strings.TrimSpace(fmt.Sprintf("%v", vr.Code))) == 0 {
		res["code"] = "code cannot be empty"
	}

	if len(res) > 0 {
		return res, errors.New("Error")
	} else {
		return res, nil
	}
}
