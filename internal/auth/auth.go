package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey from header
// Example:
// Authorization: ApiKey <token>
func GetAPIKey(header http.Header) (string, error) {
	headerVal := header.Get("Authorization")
	if headerVal == "" {
		return "", errors.New("authorization header not found")
	}

	values := strings.Split(headerVal, " ")
	if len(values) != 2 {
		return "", errors.New("invalid authorization header")
	}

	if values[0] != "ApiKey" {
		return "", errors.New("invalid authorization header value")
	}

	return values[1], nil
}
