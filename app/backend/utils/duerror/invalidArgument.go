package duerror

type InvalidArgumentError struct {
	msg string
}

func (e *InvalidArgumentError) Error() string {
	return e.msg
}

func NewInvalidArgumentError(msg string) *InvalidArgumentError {
	return &InvalidArgumentError{msg: msg}
}
