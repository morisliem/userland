package mydetail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"userland/api/response"
)

type updateEmailReq struct {
	email string
	token string
}

func TestUpdateEmail(t *testing.T) {
	tt := []struct {
		name       string
		input      updateEmailReq
		statusCode int
		expected   string
	}{
		{
			name: "no access token",
			input: updateEmailReq{
				email: "hello@gmail.com",
				token: "false",
			},
			statusCode: 401,
			expected:   fmt.Sprintln(`{"code":401,"message":"token is expired"}`),
		}, {
			name: "good input",
			input: updateEmailReq{
				email: "hello@gmail.com",
				token: "true",
			},
			statusCode: 201,
			expected:   fmt.Sprintln(`{"status":true}`),
		},
		{
			name: "email bad syntax",
			input: updateEmailReq{
				email: "hellogmail.com",
				token: "true",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"email":"email is missing @"}}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var request UpdateEmailRequest
		at, _ := strconv.ParseBool(r.Header.Get("at"))
		if !at {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request("token is expired"))
			return
		}

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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response.Success())
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rBody, err := json.Marshal(map[string]string{
				"email": tc.input.email,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r, err := http.NewRequest("POST", "localhost:8080/me/email", bytes.NewBuffer(rBody))
			r.Header.Add("at", tc.input.token)

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
