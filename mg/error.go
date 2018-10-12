package mg

import "fmt"

type ErrorCode int

const (
	IOError ErrorCode = 0
	ParseError
)

type MagnanimousError struct {
	Code    ErrorCode
	message string
}

func (e *MagnanimousError) Error() string {
	return e.message
}

func (e ErrorCode) String() string {
	names := [3]string{"IOError", "ParseError"}
	return names[e]
}

func (e *MagnanimousError) String() string {
	return fmt.Sprintf("%s: %s", e.Code.String(), e.message)
}

func NewParseError(location Location, message string) *MagnanimousError {
	return &MagnanimousError{
		message: fmt.Sprintf("(%s) %s", location.String(), message),
		Code:    ParseError,
	}
}
