package validator

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func ValidateFullname(s string) (string, bool) {
	if len(strings.TrimSpace(s)) == 0 {
		return "Fullname is required", false
	}
	if len(s) > 64 {
		return "Fullname cannot exceed 256 characters", false
	}
	return "", true
}

func ValidateEmail(s string) (string, bool) {
	emailRegex, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if err != nil {
		fmt.Println(err)
		return "Sorry, something went wrong", false
	}

	er := emailRegex.MatchString(s)
	if !er {
		return "Email address is not valid", false
	}

	if len(strings.TrimSpace(s)) == 0 {
		return "Email is required", false
	}
	if len(s) > 128 {
		return "Email cannot exceed 256 characters", false
	}

	// Checking email domain
	i := strings.Index(s, "@")
	host := s[i+1:]

	_, err = net.LookupMX(host)
	if err != nil {
		return "Counld not find email's domain server", false
	}

	return "", true
}

func ValidatePassword(s string) (string, bool) {
	// status := true
	if len(strings.TrimSpace(s)) == 0 {
		return "Password is required", false
	}
	if len(s) > 128 {
		return "Password cannot exceed 128 characters", false
	}
	if len(s) < 8 {
		return "Password must be more than 8 characters", false
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
		return "", true
	} else {
		if !hasLowerCase {
			return "Password must have a lowercase character", false
		}
		if !hasUpperCase {
			return "Password must have a uppercase character", false
		}
		if !hasNumber {
			return "Password must have a numerical character", false
		}
		return "", true
	}
}
