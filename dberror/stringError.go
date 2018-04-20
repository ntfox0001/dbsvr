package dberror

type StringError struct {
	errorStr string
}

func (ce StringError) Error() string {
	return ce.errorStr
}

func NewStringErr(str string) StringError {
	return StringError{errorStr: str}
}
