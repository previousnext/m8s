package main

import (
	"fmt"
	"strconv"

	"github.com/previousnext/k8s-black-death/retention"
)

// Helper function for returning the black death value.
func getBlackDeath(annotations map[string]string) (int64, error) {
	if val, ok := annotations[retention.Annotation]; ok {
		return strconv.ParseInt(val, 10, 64)
	}

	return 0, fmt.Errorf("annotation does not exist")
}
