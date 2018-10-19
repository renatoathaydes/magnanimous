package mg

import "fmt"

type ErrorCode int

const (
	IOError ErrorCode = 0
	ParseError
	InclusionCycleError
)

type MagnanimousError struct {
	Code    ErrorCode
	message string
}

func (e *MagnanimousError) Error() string {
	return e.message
}

func (e ErrorCode) String() string {
	names := [2]string{"IOError", "ParseError"}
	return names[e]
}

func (e *MagnanimousError) String() string {
	return fmt.Sprintf("(%s) %s", e.Code.String(), e.message)
}

func NewError(location Location, code ErrorCode, message string) error {
	return &MagnanimousError{
		message: fmt.Sprintf("(%s) %s", location.String(), message),
		Code:    code,
	}
}
