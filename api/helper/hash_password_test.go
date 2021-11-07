package helper

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tt := []struct {
		name     string
		password string
		result   int
	}{
		{
			name:     "tc1",
			password: "helloworld",
			result:   60,
		},
		{
			name:     "tc2",
			password: "abcde12345678910",
			result:   60,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pwd, err := HashPassword(tc.password)
			if err != nil {
				t.Errorf("error occur %v", err)
			}

			if len(pwd) != tc.result {
				t.Errorf("expect %v, got %v", tc.result, len(pwd))
			}
		})
	}
}

func TestCompareHashPassword(t *testing.T) {
	tt := []struct {
		name         string
		password     string
		hashPassword string
		result       bool
	}{
		{
			name:         "success",
			password:     "helloworld",
			hashPassword: "$2a$10$g/xRF9KUwgXocyyCF.qQ9us4LLRcs9wUov5VHemKoU5u0UasWlMSK",
			result:       true,
		},
		{
			name:         "fail",
			password:     "helloworld",
			hashPassword: "$2a$10$g/xRF9KUwgXocyyCF.qQ9us4LLRcs9wUov5VHemKoU5u0UasWlMS",
			result:       false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			same := ComparePasswordHash(tc.password, tc.hashPassword)
			if same != tc.result {
				t.Errorf("expect %v, got %v", tc.result, same)
			}
		})
	}

}
