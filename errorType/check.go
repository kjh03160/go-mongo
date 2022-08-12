package errorType

import (
	"reflect"

	"mongo-orm/util"

	"go.mongodb.org/mongo-driver/mongo"
)

func IsErrorTypeOf(err error, v interface{}) bool {
	t := util.GetInterfaceType(v)
	errorType := reflect.TypeOf(err)

	if t == errorType {
		return true
	}
	return false
}

func IsNotFoundErr(err error) bool {
	switch err.(type) {
	case *notFoundError:
		return true
	default:
		return false
	}
}

func IsNotModifiedAnyErr(err error) bool {
	switch err.(type) {
	case *notModifiedError:
		return true
	default:
		return false
	}
}

func IsDuplicatedKeyErr(err error) bool {
	switch err.(type) {
	case *duplicatedKeyError:
		return true
	default:
		return false
	}
}

func IsTimeoutError(err error) bool {
	switch err.(type) {
	case *timeoutError:
		return true
	default:
		return false
	}
}

func IsMongoClientError(err error) bool {
	switch err.(type) {
	case *mongoClientError:
		return true
	default:
		return false
	}
}

func IsDBInternalErr(err error) bool {
	switch err.(type) {
	case *internalError,
		*timeoutError,
		*mongoClientError:
		return true
	default:
		return false
	}
}

func IsDecodeError(err error) bool {
	switch err.(type) {
	case *decodeError:
		return true
	default:
		return false
	}
}

func ParseAndReturnDBError(err error, collection string, filter, update, doc interface{}) error {
	if err == mongo.ErrNoDocuments {
		return NotFoundError(collection, filter, update, doc)
	}

	if err == NotMatchedAnyErr {
		return NotFoundError(collection, filter, update, doc)
	}

	if err == NotModifiedAnyErr {
		return NotModifiedError(collection, filter, update, doc, err)
	}

	if mongo.IsDuplicateKeyError(err) {
		return DuplicatedKeyError(collection, filter, update, doc, err)
	}

	if mongo.IsTimeout(err) {
		return TimeoutError(collection, filter, update, doc, err)
	}

	return InternalError(collection, filter, update, doc, err)
}
