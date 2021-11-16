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

// type valEmailReq struct {
// 	email string
// 	code  int
// }

// func TestValidateEmail(t *testing.T) {
// 	tt := []struct {
// 		name       string
// 		input      valEmailReq
// 		statusCode int
// 		expected   string
// 	}{
// 		{
// 			name: "empty input",
// 			input: valEmailReq{
// 				email: "",
// 			},
// 			statusCode: 422,
// 			expected:   fmt.Sprintln(`{"Fields":{"email":"email is required"}}`),
// 		},
// 		{
// 			name: "good input",
// 			input: valEmailReq{
// 				email: "hello@gmail.com",
// 				code:  1234,
// 			},
// 			statusCode: 200,
// 			expected:   fmt.Sprintln(`{"status":true}`),
// 		},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockUserStore := new(store.MockUserStore)
// 			mockTokenStore := new(store.MockTokenStore)

// 			firstReq := ValidateEmailCodeRequest{
// 				Email: tc.input.email,
// 				Code:  tc.input.code,
// 			}

// 			newReq := store.User{
// 				Email:   tc.input.email,
// 				VerCode: tc.input.code,
// 			}

// 			rBody, err := json.Marshal(ValidateEmailCodeRequest{firstReq.Email, firstReq.Code})

// 			if err != nil {
// 				t.Errorf("failed %v", err)
// 			}

// 			r := httptest.NewRequest("POST", "localhost:8080/auth/password/forget", bytes.NewBuffer(rBody))
// 			w := httptest.NewRecorder()

// 			mockUserStore.On("EmailExist", context.Background(), newReq).Return(nil)
// 			mockUserStore.On("GetUserid", context.Background(), newReq.Email).Return("123", nil)
// 			mockTokenStore.On("GetEmailVarificationCode", "123").Return(1234, nil)
// 			mockTokenStore.On("GetNewEmail", "123").Return(newReq.Email, nil)

// 			newData := store.User{
// 				Email:   tc.input.email,
// 				Id:      "123",
// 				VerCode: tc.input.code,
// 			}
// 			mockUserStore.On("ValidateCode", context.Background(), newData).Return(nil)

// 			ValidateEmail(mockUserStore, mockTokenStore).ServeHTTP(w, r)

// 			assert.Equal(t, w.Code, tc.statusCode)
// 			assert.Equal(t, w.Body.String(), tc.expected)

// 		})
// 	}
// }
