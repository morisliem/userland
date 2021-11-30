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
