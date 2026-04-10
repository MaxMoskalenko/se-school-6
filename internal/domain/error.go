package domain

type Error struct {
	code int
	error
}

func NewError(code int, err error) *Error {
	return &Error{
		code:  code,
		error: err,
	}
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	return e.Error()
}
