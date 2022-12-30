package mongo

import "time"

type Logger interface {
	SlowQuery(msg string)
	GetTimeoutDuration() time.Duration
	GetSlowQueryDurationOfOne() time.Duration
	GetSlowQueryDurationOfMany() time.Duration
	GetSlowQueryDurationOfBulk() time.Duration
	GetSlowQueryDurationOfAggregation() time.Duration
}
