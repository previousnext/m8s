package htpasswd

import "golang.org/x/crypto/bcrypt"

// Hash is for generating a http auth password.
func Hash(pass string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordBytes), nil
}
