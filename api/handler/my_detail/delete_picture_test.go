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

type deletePictureReq struct {
	token string
}

type deletePictMockData struct {
	picture string
}

func TestDeletePicture(t *testing.T) {
	tt := []struct {
		name       string
		input      deletePictureReq
		dbData     deletePictMockData
		statusCode int
		expected   string
	}{
		{
			name: "no access token",
			input: deletePictureReq{
				token: "false",
			},
			statusCode: 401,
			expected:   fmt.Sprintln(`{"code":401,"message":"token is expired"}`),
		}, {
			name: "good input",
			input: deletePictureReq{
				token: "true",
			},
			dbData: deletePictMockData{
				picture: "gopher.png",
			},
			statusCode: 200,
			expected:   fmt.Sprintln(`{"status":true}`),
		}, {
			name: "no picture exist",
			input: deletePictureReq{
				token: "true",
			},
			statusCode: 500,
			expected:   fmt.Sprintln(`{"Message":"unable to find picture"}`),
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request, db deletePictMockData) {
		at, _ := strconv.ParseBool(r.Header.Get("at"))
		if !at {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request("token is expired"))
			return
		}

		if len(db.picture) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Response("unable to find picture"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.Success())
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest("POST", "localhost:8080/me/delete", nil)
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
