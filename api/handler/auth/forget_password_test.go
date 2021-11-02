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

type forget_password_req struct {
	email string
}

func TestForgetPassword(t *testing.T) {
	tt := []struct {
		name       string
		input      forget_password_req
		statusCode int
		expected   string
	}{
		{
			name: "empty email",
			input: forget_password_req{
				email: "",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"message":"email is required"}}`),
		},
		{
			name: "good input",
			input: forget_password_req{
				email: "moris@gmail.com",
			},
			statusCode: 200,
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var request ForgetPasswordRequest
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

	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rBody, err := json.Marshal(map[string]string{
				"email": tc.input.email,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r := httptest.NewRequest("POST", "localhost:8080/auth/password/forget", bytes.NewBuffer(rBody))
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
