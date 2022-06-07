package errorType

import (
	"fmt"
	"github.com/pkg/errors"
)

var SingleResultErr = errors.New("single result is nil")

type basicQueryInfo struct {
	collection string
	filter     interface{}
}
type notFoundError struct {
	basicQueryInfo
}

type internalError struct {
	basicQueryInfo
	update interface{}
	error
}

func NotFoundError(col string, filter interface{}) error {
	err := notFoundError{}
	err.collection = col
	err.filter = filter
	return err
}

func InternalError(col string, filter interface{}, mongoErr error, args ...interface{}) error {
	err := internalError{}
	err.collection = col
	err.filter = filter
	err.update = args
	err.error = mongoErr
	return err
}

func (e notFoundError) Error() string {
	return fmt.Sprintf("%s not found. filter: %+v", e.collection, e.filter)
}

func (e internalError) Error() string {
	msg := fmt.Sprintf("mongo internal err: %s", e.error.Error())

	if e.filter != nil {
		msg += fmt.Sprintf(" filter: %+v", e.filter)
	}
	if e.update != nil {
		msg += fmt.Sprintf(" update: %+v", e.update)
	}
	return msg
}
