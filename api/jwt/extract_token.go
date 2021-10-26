package jwt

import (
	"net/http"
	"strings"
)

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	splitBearerToken := strings.Split(bearerToken, " ")

	if len(splitBearerToken) == 2 {
		return splitBearerToken[1]
	}

	return ""
}
