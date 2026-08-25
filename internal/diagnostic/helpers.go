package diagnostic

// NewCode constructs an Error using the default message and hint for code.
func NewCode(code string) *Error {
	return New(code, "")
}

// WrapCode annotates err using the default message and hint for code.
func WrapCode(err error, code string) *Error {
	return Wrap(err, code, "")
}
