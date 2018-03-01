package utils

import (
	"strings"
)

// Machine returns a machine safe name.
func Machine(name string) string {
	n := strings.Replace(name, "_", "-", -1)
	n = strings.Replace(n, "/", "", -1)
	n = strings.Replace(n, ".", "", -1)
	n = strings.Replace(n, ":", "", -1)
	return n
}
