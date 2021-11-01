package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"userland/api/helper"
	"userland/api/jwt"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	clientid string
}

type GetATResponse struct {
	Value      string    `json:"value"`
	Type       string    `json:"type"`
	Expired_at time.Time `json:"expired_at"`
}

func Login(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest
		ctx := r.Context()
		json.NewDecoder(r.Body).Decode(&request)
		request.clientid = r.Header.Get("X-Api-Clientid")

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

		is_active, err := userStore.EmailActive(ctx, newLogin)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		// To check if the user has activated their email or not
		if !is_active {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request("email still inactive"))
			return
		}

		userId, err := userStore.GetUserId(ctx, newLogin)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		// Generate access token but still missing the refresh token id
		ts, err := jwt.GenerateAccessToken(userId, "", "")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		saveErr := jwt.CreateATAuth(userId, ts, tokenStore)
		if saveErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(saveErr.Error()))
			return
		}

		ip, err := helper.GetUserIp()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		// device := r.Header["X-Api-Clientid"]

		// Add session here
		err = userStore.SetUserSession(ctx, ts, userId, ip, request.clientid)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		at := &GetATResponse{
			Value:      ts.AccessToken,
			Type:       "jwt",
			Expired_at: time.Unix(ts.AtExpires, 0),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(at)
	}
}

func (lr *LoginRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}

	if len(strings.TrimSpace(lr.clientid)) == 0 {
		res["X-Api-ClientId"] = "required x-api-clientid"
	}

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
