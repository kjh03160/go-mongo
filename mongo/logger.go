package mongo

import "time"

type Logger interface {
	SlowQuery(msg string)
	GetTimeoutDuration() time.Duration // context cancel timeout setting
	// slow query time settings
	GetSlowQueryDurationOfOne() time.Duration
	GetSlowQueryDurationOfMany() time.Duration
	GetSlowQueryDurationOfBulk() time.Duration
	GetSlowQueryDurationOfAggregation() time.Duration
}
