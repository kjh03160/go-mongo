package errorType

import (
	"fmt"

	"github.com/pkg/errors"
)

var SingleResultErr = errors.New("single result is nil")

type basicQueryInfo struct {
	collection string
	filter     interface{}
	update     interface{}
	doc        interface{}
}

type notFoundError struct {
	basicQueryInfo
}

type duplicatedKeyError struct {
	basicQueryInfo
	error
}

type internalError struct {
	basicQueryInfo
	error
}

func NotFoundError(col string, filter, update, doc interface{}) error {
	err := notFoundError{}
	err.filter = filter
	err.collection = col
	err.update = update
	err.doc = doc
	return err
}

func DuplicatedKeyError(col string, filter, update, doc interface{}, mongoErr error) error {
	err := duplicatedKeyError{}
	err.filter = filter
	err.collection = col
	err.update = update
	err.doc = doc
	err.error = mongoErr
	return err
}

func InternalError(col string, filter, update, doc interface{}, mongoErr error) error {
	err := internalError{}
	err.collection = col
	err.filter = filter
	err.update = update
	err.doc = doc
	err.error = mongoErr
	return err
}

func (e notFoundError) Error() string {
	return fmt.Sprintf("%s not found. ", e.collection) + getBasicInfoErrorMsg(e.basicQueryInfo)
}

func (e duplicatedKeyError) Error() string {
	return fmt.Sprintf("%s failed to write due to duplicated key, err: %s ", e.collection, e.error.Error()) + getBasicInfoErrorMsg(e.basicQueryInfo)
}

func (e internalError) Error() string {
	return fmt.Sprintf("mongo internal err: %s ", e.error.Error()) + getBasicInfoErrorMsg(e.basicQueryInfo)
}

func getBasicInfoErrorMsg(e basicQueryInfo) string {
	msg := "| {query info: "
	if e.filter != nil {
		msg += fmt.Sprintf(" filter: %+v", e.filter)
	}

	if e.update != nil {
		msg += fmt.Sprintf(", update: %+v", e.update)
	}

	if e.doc != nil {
		msg += fmt.Sprintf(", doc: %+v", e.doc)
	}
	msg += "}"
	return msg
}
