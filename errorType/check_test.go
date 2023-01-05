package errorType

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	decodeErr   = DecodeError("", nil, nil, nil, errors.New("test"))
	notFoundErr = NotFoundError("col", nil, nil, nil)
	dupKeyErr   = DuplicatedKeyError("col", nil, nil, nil, errors.New(""))
	internalErr = InternalError("col", nil, nil, nil, errors.New(""))
	clientErr   = MongoClientError(errors.New(""))
	timeoutErr  = TimeoutError("col", nil, nil, nil, errors.New(""))
)

func Test_IsNotFoundErr(t *testing.T) {
	result := IsNotFoundErr(notFoundErr)
	assert.True(t, result)
	result = IsNotFoundErr(errors.Wrap(notFoundErr, ""))
	assert.True(t, result)

	result = IsNotFoundErr(decodeErr)
	assert.False(t, result)
	result = IsNotFoundErr(dupKeyErr)
	assert.False(t, result)
	result = IsNotFoundErr(internalErr)
	assert.False(t, result)
	result = IsNotFoundErr(clientErr)
	assert.False(t, result)
	result = IsNotFoundErr(timeoutErr)
	assert.False(t, result)
}

func Test_IsTimeoutErr(t *testing.T) {
	result := IsTimeoutError(timeoutErr)
	assert.True(t, result)
	result = IsTimeoutError(errors.Wrap(timeoutErr, ""))
	assert.True(t, result)

	result = IsTimeoutError(notFoundErr)
	assert.False(t, result)
	result = IsTimeoutError(dupKeyErr)
	assert.False(t, result)
	result = IsTimeoutError(internalErr)
	assert.False(t, result)
	result = IsTimeoutError(clientErr)
	assert.False(t, result)
	result = IsTimeoutError(decodeErr)
	assert.False(t, result)
}

func Test_IsMongoClientErr(t *testing.T) {
	result := IsMongoClientError(clientErr)
	assert.True(t, result)
	result = IsMongoClientError(errors.Wrap(clientErr, ""))
	assert.True(t, result)

	result = IsMongoClientError(notFoundErr)
	assert.False(t, result)
	result = IsMongoClientError(dupKeyErr)
	assert.False(t, result)
	result = IsMongoClientError(internalErr)
	assert.False(t, result)
	result = IsMongoClientError(decodeErr)
	assert.False(t, result)
	result = IsMongoClientError(timeoutErr)
	assert.False(t, result)
}

func Test_IsDBInternalErr(t *testing.T) {

	result := IsDBInternalErr(internalErr)
	assert.True(t, result)
	result = IsDBInternalErr(clientErr)
	assert.True(t, result)
	result = IsDBInternalErr(timeoutErr)
	assert.True(t, result)
	result = IsDBInternalErr(errors.Wrap(clientErr, ""))
	assert.True(t, result)
	result = IsDBInternalErr(errors.Wrap(internalErr, ""))
	assert.True(t, result)
	result = IsDBInternalErr(errors.Wrap(timeoutErr, ""))
	assert.True(t, result)

	result = IsDBInternalErr(decodeErr)
	assert.False(t, result)
	result = IsDBInternalErr(notFoundErr)
	assert.False(t, result)
}

func Test_IsDuplicatedKeyErr(t *testing.T) {

	result := IsDuplicatedKeyErr(dupKeyErr)
	assert.True(t, result)
	result = IsDuplicatedKeyErr(errors.Wrap(dupKeyErr, ""))
	assert.True(t, result)

	result = IsDuplicatedKeyErr(decodeErr)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(notFoundErr)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(internalErr)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(clientErr)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(timeoutErr)
	assert.False(t, result)
}

func Test_IsDecodeError(t *testing.T) {
	result := IsDecodeError(decodeErr)
	assert.True(t, result)
	result = IsDecodeError(errors.Wrap(decodeErr, ""))
	assert.True(t, result)

	result = IsDecodeError(notFoundErr)
	assert.False(t, result)
	result = IsDecodeError(dupKeyErr)
	assert.False(t, result)
	result = IsDecodeError(internalErr)
	assert.False(t, result)
	result = IsDecodeError(clientErr)
	assert.False(t, result)
	result = IsDecodeError(timeoutErr)
	assert.False(t, result)
}

func Test_ParseAndReturnDBError(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		err := mongo.ErrNoDocuments

		parsedErr := ParseAndReturnDBError(err, "", nil, nil, nil)
		assert.Error(t, parsedErr)
		assert.IsType(t, &notFoundError{}, parsedErr)
	})

	t.Run("not Matched Any", func(t *testing.T) {
		err := NotMatchedAnyErr

		parsedErr := ParseAndReturnDBError(err, "", nil, nil, nil)
		assert.Error(t, parsedErr)
		assert.IsType(t, &notFoundError{}, parsedErr)
	})

	t.Run("duplicated Key", func(t *testing.T) {
		err := mongo.WriteError{Code: 11000}

		parsedErr := ParseAndReturnDBError(err, "", nil, nil, nil)
		assert.Error(t, parsedErr)
		assert.IsType(t, &duplicatedKeyError{}, parsedErr)
	})

	t.Run("Timeout", func(t *testing.T) {
		err := context.DeadlineExceeded

		parsedErr := ParseAndReturnDBError(err, "", nil, nil, nil)
		assert.Error(t, parsedErr)
		assert.IsType(t, &timeoutError{}, parsedErr)
	})

	t.Run("Internal", func(t *testing.T) {
		err := errors.New("unexpected err")

		parsedErr := ParseAndReturnDBError(err, "", nil, nil, nil)
		assert.Error(t, parsedErr)
		assert.IsType(t, &internalError{}, parsedErr)
	})
}
