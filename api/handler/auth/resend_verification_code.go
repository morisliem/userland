package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"userland/api/helper"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type ResendVerCodeReq struct {
	Type  string `json:"type"`
	Email string `json:"recipient"`
}

func ResendVerCode(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ResendVerCodeReq
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

		uid, err := userStore.GetUserId(r.Context(), request.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		}

		rn := helper.GenerateRandomNumber()
		go helper.SendEmailVerCode(request.Email, rn)

		err = tokenStore.SetEmailVerificationCode(uid, rn)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tokenStore.SetNewEmail(uid, request.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.Success())
	}
}

func (r *ResendVerCodeReq) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}
	emailErr := validator.ValidateEmail(r.Email)
	if emailErr != nil {
		res["email"] = emailErr.Error()
	}

	if len(strings.TrimSpace(r.Type)) == 0 {
		res["type"] = "type is required"
	}

	if len(res) > 0 {
		return res, errors.New("error")
	} else {
		return res, nil
	}
}
