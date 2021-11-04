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

type validate_email_req struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}

type validateEmailMockData struct {
	code int
}

func TestValidateEmail(t *testing.T) {
	tt := []struct {
		name       string
		input      validate_email_req
		dbData     validateEmailMockData
		statusCode int
		expected   string
	}{
		{
			name: "empty input",
			input: validate_email_req{
				Email: "",
				Code:  -1,
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"email":"email is required"}}`),
		},
		{
			name: "good input",
			input: validate_email_req{
				Email: "moris@gmail.com",
				Code:  1234456,
			},
			dbData: validateEmailMockData{
				code: 1234456,
			},
			statusCode: 200,
		},
		{
			name: "wrong code",
			input: validate_email_req{
				Email: "moris@gmail.com",
				Code:  123456,
			},
			dbData: validateEmailMockData{
				code: 1234456,
			},
			statusCode: 422,
			expected:   fmt.Sprintln(`{"Fields":{"code":"invalid code"}}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request, db validateEmailMockData) {
		var request ValidateEmailCodeRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Response(err.Error()))
			return
		}

		res, err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

		if request.Code != db.code {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			tmp := map[string]string{}
			tmp["code"] = "invalid code"
			json.NewEncoder(w).Encode(response.UnproccesableEntity(tmp))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rBody, err := json.Marshal(resetPwdReq{
				Code:  tc.input.Code,
				Email: tc.input.Email,
			})

			if err != nil {
				t.Errorf("failed %v", err)
			}

			r := httptest.NewRequest("POST", "localhost:8080/auth/register/validate", bytes.NewBuffer(rBody))
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
