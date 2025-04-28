package utils

type SendError struct {
	msg string
}

func (e *SendError) Error() string {
	return e.msg
}

func NewSendError(msg string) *SendError {
	return &SendError{msg: msg}
}
