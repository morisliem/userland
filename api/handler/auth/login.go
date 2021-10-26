package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"userland/api/jwt"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest
		ctx := r.Context()
		json.NewDecoder(r.Body).Decode(&request)

		valErr, err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(valErr))
			return
		}

		newLogin := store.User{
			Email:    request.Email,
			Password: request.Password,
		}

		state, err := userStore.GetUserState(ctx, newLogin)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		if state != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			tmp := map[string]string{}
			tmp["email"] = "email still inactive"
			json.NewEncoder(w).Encode(response.UnproccesableEntity(tmp))
			return
		}

		userId, err := userStore.GetUserId(ctx, newLogin)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		ts, err := jwt.GenerateToken(userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		saveErr := jwt.CreateAuth(userId, ts, tokenStore)
		if saveErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(saveErr.Error()))
			return
		}

		cookie1 := &http.Cookie{
			Name:  "access_token",
			Value: ts.AccessToken,
		}

		cookie2 := &http.Cookie{
			Name:  "refresh_token",
			Value: ts.RefreshToken,
		}
		http.SetCookie(w, cookie1)
		http.SetCookie(w, cookie2)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.Success())
	}
}

func (lr *LoginRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}
	emailErr := validator.ValidateEmail(lr.Email)
	if emailErr != nil {
		res["email"] = emailErr.Error()
	}

	pwdErr := validator.ValidatePassword(lr.Password)
	if pwdErr != nil {
		res["password"] = pwdErr.Error()
	}

	if len(res) > 0 {
		return res, errors.New("Error")
	} else {
		return res, nil
	}
}
