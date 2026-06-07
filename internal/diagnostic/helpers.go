package diagnostic

func NewCode(code string) *Error {
	return New(code, "")
}

func WrapCode(err error, code string) *Error {
	return Wrap(err, code, "")
}
