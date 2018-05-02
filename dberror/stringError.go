package dberror

type StringError struct {
	errStr string
}

func (ce StringError) Error() string {
	return ce.errStr
}

func NewStringErr(str string) StringError {
	return StringError{errStr: str}
}
