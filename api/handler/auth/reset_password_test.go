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

type resetPwdReq struct {
	Email            string `json:"email"`
	Code             int    `json:"code"`
	Password         string `json:"password"`
	Confirm_password string `json:"confirm_password"`
}

type resetPwdMockData struct {
	code     int
	password []string
}

func TestResetPassword(t *testing.T) {
	existingPwd := []string{
		"123abcDep", "sfvrcwexs", "123abcDeLD",
	}
	tt := []struct {
		name       string
		input      resetPwdReq
		dbData     resetPwdMockData
		statusCode int
		expected   string
	}{
		{
			name: "empty input",
			input: resetPwdReq{
				Email:            "",
				Code:             -1,
				Password:         "",
				Confirm_password: "",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"confirm_password":"password is required","email":"email is required","password":"password is required"}}`),
		},
		{
			name: "good input",
			input: resetPwdReq{
				Email:            "moris@gmail.com",
				Code:             12345,
				Password:         "123abcDe",
				Confirm_password: "123abcDe",
			},
			dbData: resetPwdMockData{
				code:     12345,
				password: existingPwd,
			},
			statusCode: 201,
			expected:   fmt.Sprintln(`{"status":true}`),
		},
		{
			name: "unmatch password",
			input: resetPwdReq{
				Email:            "moris@gmail.com",
				Code:             12345,
				Password:         "123abcDe",
				Confirm_password: "123abcDee",
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"unmatch_password":"check again your password"}}`),
		}, {
			name: "try to add old password",
			input: resetPwdReq{
				Email:            "moris@gmail.com",
				Code:             12345,
				Password:         "123abcDeLD",
				Confirm_password: "123abcDeLD",
			},
			dbData: resetPwdMockData{
				code:     12345,
				password: existingPwd,
			},
			statusCode: 400,
			expected:   fmt.Sprintln(`{"code":400,"message":"you have user this password before, try another password"}`),
		}, {
			name: "wrong code",
			input: resetPwdReq{
				Email:            "moris@gmail.com",
				Code:             12345,
				Password:         "123abcDeLD",
				Confirm_password: "123abcDeLD",
			},
			dbData: resetPwdMockData{
				code:     12346,
				password: existingPwd,
			},
			statusCode: 400,
			expected:   fmt.Sprintln(`{"code":400,"message":"you entered wrong code"}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request, db resetPwdMockData) {
		var request ResetPasswordRequest
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

		if request.Code == db.code {
			for i := 0; i < 3; i++ {
				if db.password[i] == request.Password {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(response.Bad_request("you have user this password before, try another password"))
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response.Success())

		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request("you entered wrong code"))
			return
		}
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rBody, err := json.Marshal(resetPwdReq{
				Code:             tc.input.Code,
				Email:            tc.input.Email,
				Password:         tc.input.Password,
				Confirm_password: tc.input.Confirm_password,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r := httptest.NewRequest("POST", "localhost:8080/auth/password/reset", bytes.NewBuffer(rBody))
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
