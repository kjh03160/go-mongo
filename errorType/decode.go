package errorType

import "fmt"

type decodeError struct {
	basicQueryInfo
	error
}

func DecodeError(col string, filter, update, doc interface{}, mongoErr error) error {
	err := internalError{}
	err.setBasicError(col, filter, update, doc)
	err.error = mongoErr
	return err
}

func (e decodeError) Error() string {
	return fmt.Sprintf("decode document err: %s ", e.error.Error()) + getBasicInfoErrorMsg(e.basicQueryInfo)
}
