package validator

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

func ValidateFullname(s string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("fullname is required")
	}
	if len(s) > 64 {
		return errors.New("fullname cannot exceed 256 characters")
	}
	return nil
}

// Got bug sometime it works, some time it doesn't work
func ValidateEmail(s string) error {
	emailRegex, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if err != nil {
		return errors.New("sorry, something went wrong")
	}

	er := emailRegex.MatchString(s)
	if !er {
		return errors.New("email address is not valid")
	}

	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("email is required")
	}
	if len(s) > 128 {
		return errors.New("email cannot exceed 256 characters")
	}

	// Checking email domain
	i := strings.Index(s, "@")
	host := s[i+1:]

	_, err = net.LookupMX(host)
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
