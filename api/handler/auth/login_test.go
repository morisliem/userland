package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"userland/api/response"
)

type loginReq struct {
	email    string
	password string
	clientid string
}

func TestLogin(t *testing.T) {
	tt := []struct {
		name       string
		input      loginReq
		statusCode int
		expected   string
	}{
		{
			name: "empty input",
			input: loginReq{
				email:    "",
				password: "",
				clientid: "",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"X-Api-ClientId":"x-api-clientid is required","email":"email is required","password":"password is required"}}`),
		}, {
			name: "good input",
			input: loginReq{
				email:    "moris@gmail.com",
				password: "abc12345D",
				clientid: "Iphone",
			},
			statusCode: 200,
		},
	}

	// Still missing testing the db
	handler := func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest
		// ctx := r.Context()
		json.NewDecoder(r.Body).Decode(&request)
		request.clientid = r.Header.Get("X-Api-Clientid")

		valErr, err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(valErr))
			return
		}

		// newLogin := store.User{
		// 	Email:    request.Email,
		// 	Password: request.Password,
		// }

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rBody, err := json.Marshal(map[string]string{
				"email":    tc.input.email,
				"password": tc.input.password,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r, err := http.NewRequest("POST", "localhost:8080/auth/login", bytes.NewBuffer(rBody))
			r.Header.Add("X-Api-ClientId", tc.input.clientid)

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
