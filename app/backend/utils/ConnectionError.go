package utils

type ConnectionError struct {
	msg string
}

func (e *ConnectionError) Error() string {
	return e.msg
}

func NewConnectionError(msg string) *ConnectionError {
	return &ConnectionError{msg: msg}
}
