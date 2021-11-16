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

type loginMockData struct {
	email    string
	password string
}

func TestLogin(t *testing.T) {
	tt := []struct {
		name       string
		input      loginReq
		dbData     loginMockData
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
			expected:   fmt.Sprintln(`{"Fields":{"X-Api-ClientId":"X-Api-ClientId is required","email":"email is required","password":"password is required"}}`),
		}, {
			name: "good input",
			input: loginReq{
				email:    "moris@gmail.com",
				password: "abc12345D",
				clientid: "Iphone",
			},
			dbData: loginMockData{
				email:    "moris@gmail.com",
				password: "abc12345D",
			},
			statusCode: 200,
		}, {
			name: "incorrect password",
			input: loginReq{
				email:    "moris@gmail.com",
				password: "abc12345D",
				clientid: "Iphone",
			},
			dbData: loginMockData{
				email:    "moris@gmail.com",
				password: "abc12345F",
			},
			statusCode: 400,
			expected:   fmt.Sprintln(`{"Message":"password incorrect"}`),
		}, {
			name: "incorrect email",
			input: loginReq{
				email:    "moriss@gmail.com",
				password: "abc12345D",
				clientid: "Iphone",
			},
			dbData: loginMockData{
				email:    "moris@gmail.com",
				password: "abc12345D",
			},
			statusCode: 400,
			expected:   fmt.Sprintln(`{"Message":"email is not found"}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request, db loginMockData) {
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

		if request.Email == db.email && request.Password == db.password {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		} else if request.Email != db.email {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response("email is not found"))
			return
		} else if request.Password != db.password {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response("password incorrect"))
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			res := map[string]string{}
			res["email"] = "email is not found"
			res["password"] = "password incorrect"
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

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

			handler(w, r, tc.dbData)

			if w.Code != tc.statusCode {
				t.Errorf("expected %d, got %d", tc.statusCode, w.Code)
			}

			if w.Body.String() != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, w.Body.String())
			}
		})
	}

}
