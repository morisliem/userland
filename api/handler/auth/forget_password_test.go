package auth

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http/httptest"
// 	"testing"
// 	"userland/store"

// 	"github.com/stretchr/testify/assert"
// )

// func TestForgetPassword(t *testing.T) {
// 	tt := []struct {
// 		name       string
// 		input      string
// 		statusCode int
// 		expected   string
// 	}{
// 		{
// 			name:       "empty email",
// 			input:      "",
// 			statusCode: 422,
// 			expected:   fmt.Sprintln(`{"Fields":{"message":"email is required"}}`),
// 		}, {
// 			name:       "good input",
// 			input:      "moris@gmail.com",
// 			statusCode: 200,
// 			expected:   fmt.Sprintln(`{"status":true}`),
// 		},
// 	}

// 	mockUserStore := new(store.MockUserStore)
// 	mockTokenStore := new(store.MockTokenStore)

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			rBody, err := json.Marshal(map[string]string{
// 				"email": tc.input,
// 			})

// 			if err != nil {
// 				t.Errorf("failed %v", err)
// 			}

// 			r := httptest.NewRequest("POST", "localhost:8080/auth/password/forget", bytes.NewBuffer(rBody))
// 			w := httptest.NewRecorder()

// 			newReq := store.User{
// 				Email: tc.input,
// 			}

// 			mockUserStore.On("EmailExist", context.Background(), newReq).Return(nil)
// 			mockTokenStore.On("SetEmailVerificationCode", newReq.Email, 10).Return(nil)

// 			ForgetPassword(mockUserStore, mockTokenStore).ServeHTTP(w, r)

// 			assert.Equal(t, w.Code, tc.statusCode)
// 			assert.Equal(t, w.Body.String(), tc.expected)
// 		})
// 	}
// }

// // func TestForgetPassword(t *testing.T) {
// // 	tt := []struct {
// // 		name       string
// // 		input      forget_password_req
// // 		dbData     forgetPwdMockData
// // 		statusCode int
// // 		expected   string
// // 	}{
// // 		{
// // 			name: "empty email",
// // 			input: forget_password_req{
// // 				email: "",
// // 			},
// // 			statusCode: 422,
// // 			expected:   fmt.Sprintln(`{"Fields":{"message":"email is required"}}`),
// // 		}, {
// // 			name: "good input",
// // 			input: forget_password_req{
// // 				email: "moris@gmail.com",
// // 			},
// // 			dbData: forgetPwdMockData{
// // 				email: "moris@gmail.com",
// // 			},
// // 			statusCode: 200,
// // 		}, {
// // 			name: "unmatch email",
// // 			input: forget_password_req{
// // 				email: "moris@gmail.com",
// // 			},
// // 			dbData: forgetPwdMockData{
// // 				email: "moris9@gmail.com",
// // 			},
// // 			statusCode: 400,
// // 			expected:   fmt.Sprintln(`{"code":400,"message":"unable to find user"}`),
// // 		},
// // 	}

// // 	handler := func(w http.ResponseWriter, r *http.Request, db forgetPwdMockData) {
// // 		var request ForgetPasswordRequest
// // 		err := json.NewDecoder(r.Body).Decode(&request)
// // 		if err != nil {
// // 			w.Header().Set("Content-Type", "application/json")
// // 			w.WriteHeader(http.StatusBadRequest)
// // 			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
// // 			return
// // 		}

// // 		res, err := request.ValidateRequest()
// // 		if err != nil {
// // 			w.Header().Set("Content-Type", "application/json")
// // 			w.WriteHeader(http.StatusUnprocessableEntity)
// // 			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
// // 			return
// // 		}

// // 		if request.Email != db.email {
// // 			w.Header().Set("Content-Type", "application/json")
// // 			w.WriteHeader(http.StatusBadRequest)
// // 			json.NewEncoder(w).Encode(response.Bad_request("unable to find user"))
// // 			return
// // 		}

// // 		w.Header().Set("Content-Type", "application/json")
// // 		w.WriteHeader(http.StatusOK)

// // 	}

// // 	for _, tc := range tt {
// // 		t.Run(tc.name, func(t *testing.T) {
// // 			rBody, err := json.Marshal(map[string]string{
// // 				"email": tc.input.email,
// // 			})

// // 			if err != nil {
// // 				t.Errorf("failed %v", err)
// // 			}

// // 			r := httptest.NewRequest("POST", "localhost:8080/auth/password/forget", bytes.NewBuffer(rBody))
// // 			w := httptest.NewRecorder()

// // 			handler(w, r, tc.dbData)
// // 			if w.Code != tc.statusCode {
// // 				t.Errorf("expected %d, got %d", tc.statusCode, w.Code)
// // 			}

// // 			if w.Body.String() != tc.expected {
// // 				t.Errorf("expected %s, got %s", tc.expected, w.Body.String())
// // 			}
// // 		})
// // 	}
// // }
