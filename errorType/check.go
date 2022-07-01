package errorType

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFoundErr(err error) bool {
	switch err.(type) {
	case notFoundError,
		*notFoundError:
		return true
	default:
		return false
	}
}

func IsDBInternalErr(err error) bool {
	switch err.(type) {
	case internalError,
		*internalError:
		return true
	default:
		return false
	}
}

func ParseAndReturnDBError(err error, collection string, filter, update, doc interface{}) error {
	if err == mongo.ErrNoDocuments {
		return NotFoundError(collection, filter, update, doc)
	}

	if mongo.IsDuplicateKeyError(err) {
		return DuplicatedKeyError(collection, filter, update, doc, err)
	}

	return InternalError(collection, filter, update, doc, err)
}
