package environ

import (
	"strings"
)

// Get returns a list of environment variables, with their prefix stripped.
//  eg. "FOO=bar,MY_PREFIX_FOO=baz" will result in "FOO=baz" if "MY_PREFIX" is provided.
func GetWithPrefix(prefix string, list []string) map[string]string {
	return filter(prefix, list)
}

// Helper function to filter down a list of environment variables to only prefixed.
func filter(prefix string, list []string) map[string]string {
	prefixed := make(map[string]string)

	for _, item := range list {
		sl := strings.Split(item, "=")

		if len(sl) != 2 {
			continue
		}

		if strings.HasPrefix(sl[0], prefix) {
			name := strings.TrimPrefix(sl[0], prefix)
			prefixed[name]= sl[1]
		}
	}

	return prefixed
}
