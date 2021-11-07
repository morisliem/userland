package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"userland/api/helper"
	"userland/api/response"
	"userland/store"

	"github.com/gofrs/uuid"
)

type registerReq struct {
	fullname         string
	email            string
	password         string
	confirm_password string
}

func TestRegister(t *testing.T) {
	tt := []struct {
		name       string
		input      registerReq
		statusCode int
		expected   string
	}{
		{
			name: "password syntex error",
			input: registerReq{
				fullname:         "moris",
				email:            "moris@gmail.com",
				password:         "abcde123",
				confirm_password: "abcde123",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"password":"password must have a uppercase character"}}`),
		},
		{
			name: "good request",
			input: registerReq{
				fullname:         "moris",
				email:            "moris@gmail.com",
				password:         "123abcdE",
				confirm_password: "123abcdE",
			},
			statusCode: 201,
		},
		{
			name: "unmatch password",
			input: registerReq{
				fullname:         "moris",
				email:            "moris@gmail.com",
				password:         "123abcdP",
				confirm_password: "123abcdE",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"unmatch_password":"check again your password"}}`),
		},
		{
			name: "empty input",
			input: registerReq{
				fullname:         "",
				email:            "",
				password:         "",
				confirm_password: "",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"email":"email is required","fullname":"fullname is required","password":"password is required"}}`)},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		// var userStore store.UserStore
		// var tokenStore store.TokenStore
		var request RegisterRequest
		// ctx := r.Context()
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

		go helper.SendEmailVerCode(newRegister.Email, rn)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			rBody, err := json.Marshal(map[string]string{
				"fullname":         tc.input.fullname,
				"email":            tc.input.email,
				"Password":         tc.input.password,
				"Confirm_password": tc.input.confirm_password,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r, err := http.NewRequest("POST", "localhost:8080/auth/register", bytes.NewBuffer(rBody))

			if err != nil {
				t.Errorf("failed %v", err)
			}

			w := httptest.NewRecorder()

			handler(w, r)

			if w.Code != tc.statusCode {
				t.Errorf("expected %d, got %d", tc.statusCode, w.Code)
			}

			if w.Body.String() != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, w.Body.String())
			}

		})
	}
}
