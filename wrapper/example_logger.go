package wrapper

import (
	"time"

	"github.com/sirupsen/logrus"
)

type myLogger struct {
	*logrus.Logger
}

func (l *myLogger) SlowQuery(msg string) {
	l.Error(msg)
}

func (l *myLogger) GetTimeoutDuration() time.Duration {
	return 10 * time.Second
}

func (l *myLogger) GetSlowQueryDurationOfOne() time.Duration {
	return 1 * time.Second
}

func (l *myLogger) GetSlowQueryDurationOfMany() time.Duration {
	return 2 * time.Second
}

func (l *myLogger) GetSlowQueryDurationOfBulk() time.Duration {
	return 3 * time.Second
}

func (l *myLogger) GetSlowQueryDurationOfAggregation() time.Duration {
	return 10 * time.Second
}
