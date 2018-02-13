package auth

import (
	"encoding/base64"
	"encoding/json"

	"docker.io/go-docker/api/types"
)

// Base64 converts a username and password into a base64 string.
func Base64(username, password string) (string, error) {
	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}

	authBytes, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(authBytes), nil
}
