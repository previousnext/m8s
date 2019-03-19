package env

import (
	"strconv"
	"strings"
	"time"
)

// Helper function to turn a project into a safe.
func machine(name string) string {
	n := strings.Replace(name, "_", "-", -1)
	n = strings.Replace(n, "/", "", -1)
	n = strings.Replace(n, ".", "", -1)
	return n
}

// Helper function to convert retention into a future unix timestamp.
func retentionToUnix(now time.Time, retention string) (string, error) {
	duration, err := time.ParseDuration(retention)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(now.Local().Add(duration).Unix(), 10), nil
}
