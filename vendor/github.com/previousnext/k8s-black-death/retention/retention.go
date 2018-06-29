package retention

import (
	"time"
	"strconv"
)

// Annotation used for declaring when an environment should be termintated.
const Annotation = "black-death.skpr.io"

// Unix returns a unix string based on a retention period.
func Unix(duration time.Duration) (string, error) {
	return strconv.FormatInt(time.Now().Local().Add(duration).Unix(), 10), nil
}
