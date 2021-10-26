package validator

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func ValidateFullname(s string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("fullname is required")
	}
	if len(s) > 64 {
		return errors.New("fullname cannot exceed 64 characters")
	}
	return nil
}

func ValidateEmail(s string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("email is required")
	}
	if len(s) > 128 {
		return errors.New("email cannot exceed 128 characters")
	}

	if !strings.Contains(s, ".") {
		return errors.New("email is missing . ")
	}

	if !strings.Contains(s, "@") {
		return errors.New("email is missing @")
	}

	// Checking email domain
	i := strings.Index(s, "@")
	host := s[i+1:]

	_, err := net.LookupMX(host)
	if err != nil {
		return errors.New("could not find email's domain server")
	}

	return nil
}

func ValidatePassword(s string) error {
	// status := true
	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("password is required")
	}
	if len(s) > 128 {
		return errors.New("password cannot exceed 128 characters")
	}
	if len(s) < 8 {
		return errors.New("password must be more than 8 characters")
	}

	hasUpperCase := false
	hasLowerCase := false
	hasNumber := false

	for _, v := range s {
		if int(v) >= 48 && int(v) <= 57 {
			hasNumber = true
		}
		if int(v) >= 65 && int(v) <= 90 {
			hasUpperCase = true
		}
		if int(v) >= 97 && int(v) <= 122 {
			hasLowerCase = true
		}
	}

	if hasUpperCase && hasLowerCase && hasNumber {
		return nil
	} else {
		if !hasLowerCase {
			return errors.New("password must have a lowercase character")
		}
		if !hasUpperCase {
			return errors.New("password must have a uppercase character")
		}
		if !hasNumber {
			return errors.New("password must have a numerical character")
		}
		return nil
	}
}

func ValidateCode(code int) error {
	tmp := fmt.Sprintf("%v", code)

	if len(strings.TrimSpace(tmp)) == 0 {
		return errors.New("code is required")
	}

	return nil

}
