package environ

import (
	"fmt"
	"os"
	"strings"
)

const Prefix = "PR_ENV_"

// Get returns a list of environment variables, with their prefix stripped.
//  eg. "FOO=bar,MY_PREFIX_FOO=baz" will result in "FOO=baz" if "MY_PREFIX" is provided.
func Get() []string {
	return filter(os.Environ())
}

// Helper function to filter down a list of environment variables to only prefixed.
func filter(list []string) []string {
	var prefixed []string

	for _, item := range list {
		sl := strings.Split(item, "=")

		if len(sl) != 2 {
			continue
		}

		if strings.HasPrefix(sl[0], Prefix) {
			prefixed = append(prefixed, fmt.Sprintf("%s=%s", strings.TrimPrefix(sl[0], Prefix), sl[1]))
		}
	}

	return prefixed
}
