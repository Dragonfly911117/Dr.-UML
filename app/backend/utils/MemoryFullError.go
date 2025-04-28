package utils

type MemoryFullError struct {
	msg string
}

func (e *MemoryFullError) Error() string {
	return e.msg
}

func NewMemoryFullError(msg string) *MemoryFullError {
	return &MemoryFullError{msg: msg}
}
