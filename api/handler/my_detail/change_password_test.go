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

type changePwdReq struct {
	token            string
	current_password string
	password         string
	confirm_password string
}

type changePwdMockData struct {
	current_password string
	password         []string
}

func TestChangePassword(t *testing.T) {
	existingPwd := []string{"123abcDep", "sfvrcwexs", "123abcDeLD"}

	tt := []struct {
		name       string
		input      changePwdReq
		dbData     changePwdMockData
		statusCode int
		expected   string
	}{
		{
			name: "empty input",
			input: changePwdReq{
				current_password: "",
				password:         "",
				confirm_password: "",
				token:            "true",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"confirm_password":"password is required","current_password":"password is required","password":"password is required"}}`),
		}, {
			name: "unmatch password",
			input: changePwdReq{
				current_password: "abcdeF123",
				password:         "abcdeF123L",
				confirm_password: "abcdeF123d",
				token:            "true",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"unmatch_password":"check again your password"}}`),
		}, {
			name: "good input",
			input: changePwdReq{
				current_password: "abcdeF123",
				password:         "abcdeF123D",
				confirm_password: "abcdeF123D",
				token:            "true",
			},
			dbData: changePwdMockData{
				current_password: "abcdeF123",
				password:         existingPwd,
			},
			statusCode: 201,
			expected:   fmt.Sprintln(`{"status":true}`),
		}, {
			name: "adding old password",
			input: changePwdReq{
				current_password: "abcdeF123",
				password:         "123abcDep",
				confirm_password: "123abcDep",
				token:            "true",
			},
			dbData: changePwdMockData{
				current_password: "abcdeF123",
				password:         existingPwd,
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"password":"You have used this password before, try to use other password"}}`),
		},
		{
			name: "no access token",
			input: changePwdReq{
				current_password: "abcdeF123",
				password:         "123abcDep",
				confirm_password: "123abcDep",
				token:            "false",
			},
			dbData: changePwdMockData{
				current_password: "abcdeF123",
				password:         existingPwd,
			},
			statusCode: 401,
			expected:   fmt.Sprintln(`{"code":401,"message":"token is expired"}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request, db changePwdMockData) {
		var request ChangePasswordRequest

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

		if request.Current_Password != db.current_password {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			tmp := map[string]string{}
			tmp["message"] = "password incorrect"
			json.NewEncoder(w).Encode(response.UnproccesableEntity(tmp))
			return
		}

		for i := 0; i < 3; i++ {
			if db.password[i] == request.Password {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnprocessableEntity)
				errMsg := map[string]string{}
				errMsg["password"] = "You have used this password before, try to use other password"
				json.NewEncoder(w).Encode(response.UnproccesableEntity(errMsg))
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response.Success())
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rBody, err := json.Marshal(map[string]string{
				"current_password": tc.input.current_password,
				"password":         tc.input.password,
				"confirm_password": tc.input.confirm_password,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r, err := http.NewRequest("POST", "localhost:8080/me/password", bytes.NewBuffer(rBody))
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
