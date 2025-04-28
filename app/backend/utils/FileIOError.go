package utils

type FileIOError struct {
	msg string
}

func (e *FileIOError) Error() string {
	return e.msg
}

func NewFileIOError(msg string) *FileIOError {
	return &FileIOError{msg: msg}
}
