package step

import (
	"errors"
	"fmt"
	"strings"
)

type Logger struct {
	Logs []string
}

func (l *Logger) Add(m string) {
	l.Logs = append(l.Logs, m)
}

func (l *Logger) String() string {
	return strings.Join(l.Logs, "\n")
}

type Step interface {
	Name() string
	Run(*Logger) bool
}

func Run(steps []Step) (*Logger, error) {
	l := &Logger{}

	for _, s := range steps {
		if pass := s.Run(l); !pass {
			return l, errors.New(fmt.Sprintf("The step %s has failed", s.Name()))
		}
	}

	return l, nil
}
