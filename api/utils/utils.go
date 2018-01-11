package utils

import (
	"fmt"
	"log"
	"net/http"
)

// Param will return a parameter passed by the user.
func Param(r *http.Request, name string) (string, error) {
	value := r.URL.Query().Get(name)

	if value == "" {
		return "", fmt.Errorf("failed to retrieve param: %s", name)
	}

	return value, nil
}

// Fatal is a common function for logging and returning a fatal status code.
func Fatal(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
