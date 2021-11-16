package auth

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"userland/store"

// 	"github.com/stretchr/testify/assert"
// )

// type registerReq struct {
// 	fullname         string
// 	email            string
// 	password         string
// 	confirm_password string
// }

// func TestRegister(t *testing.T) {
// 	tt := []struct {
// 		name       string
// 		input      registerReq
// 		statusCode int
// 		expected   string
// 	}{
// 		{
// 			name: "empty input",
// 			input: registerReq{
// 				fullname:         "",
// 				email:            "",
// 				password:         "",
// 				confirm_password: "",
// 			},
// 			statusCode: 422,
// 			expected:   fmt.Sprintln(`{"Fields":{"email":"email is required","fullname":"fullname is required"}}`),
// 		},
// 		{
// 			name: "good request",
// 			input: registerReq{
// 				fullname:         "moris",
// 				email:            "moris@gmail.com",
// 				password:         "123abcdE",
// 				confirm_password: "123abcdE",
// 			},
// 			statusCode: 201,
// 		},
// 	}

// 	mockUserStore := new(store.MockUserStore)
// 	mockTokenStore := new(store.MockTokenStore)

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			rBody, err := json.Marshal(map[string]string{
// 				"fullname":         tc.input.fullname,
// 				"email":            tc.input.email,
// 				"Password":         "tc.input.passwordD1",
// 				"Confirm_password": "tc.input.passwordD1",
// 			})

// 			if err != nil {
// 				t.Errorf("failed %v", err)
// 			}

// 			r, err := http.NewRequest("POST", "localhost:8080/auth/register", bytes.NewBuffer(rBody))

// 			if err != nil {
// 				t.Errorf("failed %v", err)
// 			}

// 			w := httptest.NewRecorder()

// 			newUser := store.User{
// 				Id:       "tc.input",
// 				Fullname: tc.input.fullname,
// 				Email:    tc.input.email,
// 				Password: "tc.input.passwordD1",
// 			}

// 			mockUserStore.On("RegisterUser", context.Background(), newUser, 123).Return(nil)
// 			mockTokenStore.On("SetEmailVerificationCode", "tc.input", 123).Return(nil)
// 			mockTokenStore.On("SetNewEmail", "tc.input", tc.input.email).Return(nil)

// 			Register(mockUserStore, mockTokenStore).ServeHTTP(w, r)

// 			assert.Equal(t, w.Code, tc.statusCode)
// 			assert.Equal(t, w.Body.String(), tc.expected)

// 		})
// 	}
// }
