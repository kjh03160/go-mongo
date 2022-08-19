package errorType

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	err  = DecodeError("", nil, nil, nil, errors.New("test"))
	err1 = NotFoundError("col", nil, nil, nil)
	err2 = DuplicatedKeyError("col", nil, nil, nil, errors.New(""))
	err3 = InternalError("col", nil, nil, nil, errors.New(""))
	err4 = MongoClientError(errors.New(""))
	err5 = TimeoutError("col", nil, nil, nil, errors.New(""))
)

func Test_IsNotFoundErr(t *testing.T) {
	result := IsNotFoundErr(err1)
	assert.True(t, result)

	result = IsNotFoundErr(err)
	assert.False(t, result)
	result = IsNotFoundErr(err2)
	assert.False(t, result)
	result = IsNotFoundErr(err3)
	assert.False(t, result)
	result = IsNotFoundErr(err4)
	assert.False(t, result)
	result = IsNotFoundErr(err5)
	assert.False(t, result)
}

func Test_IsTimeoutErr(t *testing.T) {
	result := IsTimeoutError(err5)
	assert.True(t, result)

	result = IsTimeoutError(err1)
	assert.False(t, result)
	result = IsTimeoutError(err2)
	assert.False(t, result)
	result = IsTimeoutError(err3)
	assert.False(t, result)
	result = IsTimeoutError(err4)
	assert.False(t, result)
	result = IsTimeoutError(err)
	assert.False(t, result)
}

func Test_IsMongoClientErr(t *testing.T) {
	result := IsMongoClientError(err4)
	assert.True(t, result)

	result = IsMongoClientError(err1)
	assert.False(t, result)
	result = IsMongoClientError(err2)
	assert.False(t, result)
	result = IsMongoClientError(err3)
	assert.False(t, result)
	result = IsMongoClientError(err)
	assert.False(t, result)
	result = IsMongoClientError(err5)
	assert.False(t, result)
}

func Test_IsDBInternalErr(t *testing.T) {

	result := IsDBInternalErr(err3)
	assert.True(t, result)
	result = IsDBInternalErr(err3)
	assert.True(t, result)
	result = IsDBInternalErr(err4)
	assert.True(t, result)
	result = IsDBInternalErr(err5)
	assert.True(t, result)

	result = IsDBInternalErr(err)
	assert.False(t, result)
	result = IsDBInternalErr(err1)
	assert.False(t, result)
}

func Test_IsDuplicatedKeyErr(t *testing.T) {

	result := IsDuplicatedKeyErr(err2)
	assert.True(t, result)

	result = IsDuplicatedKeyErr(err)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(err1)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(err3)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(err4)
	assert.False(t, result)
	result = IsDuplicatedKeyErr(err5)
	assert.False(t, result)
}

func Test_IsDecodeError(t *testing.T) {
	result := IsDecodeError(err)
	assert.True(t, result)

	result = IsDecodeError(err1)
	assert.False(t, result)
	result = IsDecodeError(err2)
	assert.False(t, result)
	result = IsDecodeError(err3)
	assert.False(t, result)
	result = IsDecodeError(err4)
	assert.False(t, result)
	result = IsDecodeError(err5)
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
