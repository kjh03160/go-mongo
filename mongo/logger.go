package mongo

import "time"

type Logger interface {
	SlowQuery(msg string)
	GetQueryTimeoutDuration() time.Duration
}
