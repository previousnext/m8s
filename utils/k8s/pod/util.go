package pod

import "strings"

const VolumeProjectRoot = "."

func FindCodePath(mounts []string) (bool, string) {
	// Adds the Docker Compose volumes to our Pod object.
	for _, mount := range mounts {
		sl := strings.Split(mount, ":")

		// Ensure we have an volume in the format "/source:/target".
		if len(sl) < 2 {
			continue
		}

		if sl[0] == VolumeProjectRoot {
			return true, sl[1]
		}
	}

	return false, ""
}