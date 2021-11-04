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

type deleteAccountReq struct {
	password string
	token    string
}

type deleteAccountMockData struct {
	password string
}

func TestDeleteAccount(t *testing.T) {
	tt := []struct {
		name       string
		input      deleteAccountReq
		dbData     deleteAccountMockData
		statusCode int
		expected   string
	}{
		{
			name: "empty input",
			input: deleteAccountReq{
				password: "",
				token:    "true",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"password":"password is required"}}`),
		}, {
			name: "good input",
			input: deleteAccountReq{
				password: "abcde123D",
				token:    "true",
			},
			dbData: deleteAccountMockData{
				password: "abcde123D",
			},
			statusCode: 201,
			expected:   fmt.Sprintln(`{"status":true}`),
		}, {
			name: "no access token",
			input: deleteAccountReq{
				password: "abcde123D",
				token:    "false",
			},
			dbData: deleteAccountMockData{
				password: "abcde123D",
			},
			statusCode: 401,
			expected:   fmt.Sprintln(`{"code":401,"message":"token is expired"}`),
		}, {
			name: "incorrect password",
			input: deleteAccountReq{
				password: "abcde123D",
				token:    "true",
			},
			dbData: deleteAccountMockData{
				password: "abcde123Dp",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"message":"password incorrect"}}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request, db deleteAccountMockData) {
		var request DeleteAccountRequest

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

		if request.Password != db.password {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			tmp := map[string]string{}
			tmp["message"] = "password incorrect"
			json.NewEncoder(w).Encode(response.UnproccesableEntity(tmp))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response.Success())
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rBody, err := json.Marshal(map[string]string{
				"password": tc.input.password,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r, err := http.NewRequest("POST", "localhost:8080/me/delete", bytes.NewBuffer(rBody))
			r.Header.Add("at", tc.input.token)

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
