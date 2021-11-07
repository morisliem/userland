package mydetail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"userland/api/response"
)

type setPictureReq struct {
	token string
}

func TestSetPicture(t *testing.T) {
	tt := []struct {
		name       string
		input      setPictureReq
		statusCode int
		expected   string
	}{
		{
			name: "no access token",
			input: setPictureReq{
				token: "false",
			},
			statusCode: 401,
			expected:   fmt.Sprintln(`{"code":401,"message":"token is expired"}`),
		}, {
			name: "good input",
			input: setPictureReq{
				token: "true",
			},
			statusCode: 201,
			expected:   fmt.Sprintln(`{"status":true}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		at, _ := strconv.ParseBool(r.Header.Get("at"))
		if !at {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request("token is expired"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response.Success())
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest("POST", "localhost:8080/me/picture", nil)
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
