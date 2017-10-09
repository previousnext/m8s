package cloudwatchlogs

import (
	"sort"
)

func (s Logs) Len() int {
	return len(s)
}

func (s Logs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ByWeight implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Organs value.
type ByTime struct{ Logs }

func (s ByTime) Less(i, j int) bool {
	return s.Logs[i].Timestamp.Before(s.Logs[j].Timestamp)
}

// Helper function for ordering our logs.
func MergeLogs(a Logs, b Logs) Logs {
	var r Logs
	r = append(a, b...)
	sort.Sort(ByTime{r})
	return r
}
