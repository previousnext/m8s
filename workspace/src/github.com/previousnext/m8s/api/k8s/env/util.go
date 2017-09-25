package env

import (
	"strings"
)

// Helper function to turn a project into a safe.
func machine(name string) string {
	n := strings.Replace(name, "_", "-", -1)
	n = strings.Replace(n, "/", "", -1)
	return n
}
