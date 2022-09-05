package errorType

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFoundErr(err error) bool {
	switch err.(type) {
	case *notFoundError:
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
	if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, NotMatchedAnyErr) {
		return NotFoundError(collection, filter, update, doc)
	}

	if mongo.IsDuplicateKeyError(err) {
		return DuplicatedKeyError(collection, filter, update, doc, err)
	}

	if mongo.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
		return TimeoutError(collection, filter, update, doc, err)
	}

	return InternalError(collection, filter, update, doc, err)
}
