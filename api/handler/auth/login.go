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

// Still missing the refresh token
// Have to check user's state before give user access to log in
func Login(userStore store.UserStore) http.HandlerFunc {
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

		userId, err := userStore.GetUserId(ctx, newLogin)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		token, err := jwt.GenerateToken(userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		res := map[string]string{}

		res["Access token"] = token

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
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

	if len(res) > 1 {
		return res, errors.New("Error")
	} else {
		return res, nil
	}
}
