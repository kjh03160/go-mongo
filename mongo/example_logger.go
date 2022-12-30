package mongo

import (
	"time"

	"github.com/sirupsen/logrus"
)

type MyLogger struct {
	*logrus.Logger
}

func (l *MyLogger) SlowQuery(msg string) {
	l.Error(msg)
}

func (l *MyLogger) GetTimeoutDuration() time.Duration {
	return 10 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfOne() time.Duration {
	return 1 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfMany() time.Duration {
	return 2 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfBulk() time.Duration {
	return 3 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfAggregation() time.Duration {
	return 10 * time.Second
}
