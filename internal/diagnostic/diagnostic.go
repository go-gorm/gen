package diagnostic

import (
	"encoding/json"
	"fmt"
)

type Diagnostic struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	File      string `json:"file,omitempty"`
	Line      int    `json:"line,omitempty"`
	Column    int    `json:"column,omitempty"`
	Interface string `json:"interface,omitempty"`
	Method    string `json:"method,omitempty"`
	Snippet   string `json:"snippet,omitempty"`
	Hint      string `json:"hint,omitempty"`
}

type Error struct {
	Diag Diagnostic
	Err  error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	loc := ""
	if e.Diag.File != "" {
		if e.Diag.Line > 0 {
			if e.Diag.Column > 0 {
				loc = fmt.Sprintf("%s:%d:%d: ", e.Diag.File, e.Diag.Line, e.Diag.Column)
			} else {
				loc = fmt.Sprintf("%s:%d: ", e.Diag.File, e.Diag.Line)
			}
		} else {
			loc = fmt.Sprintf("%s: ", e.Diag.File)
		}
	}
	code := ""
	if e.Diag.Code != "" {
		code = e.Diag.Code + ": "
	}
	return loc + code + e.Diag.Message
}

func (e *Error) Unwrap() error { return e.Err }

func (e *Error) MarshalJSON() ([]byte, error) {
	type payload struct {
		Diagnostic
		Cause string `json:"cause,omitempty"`
	}
	p := payload{Diagnostic: e.Diag}
	if e.Err != nil {
		p.Cause = e.Err.Error()
	}
	return json.Marshal(p)
}

func New(code, message string) *Error {
	return &Error{Diag: Diagnostic{Code: code, Message: message}}
}

func Wrap(err error, code, message string) *Error {
	if err == nil {
		return New(code, message)
	}
	if e, ok := err.(*Error); ok {
		if code != "" {
			e.Diag.Code = code
		}
		if message != "" {
			e.Diag.Message = message
		}
		return e
	}
	return &Error{Diag: Diagnostic{Code: code, Message: message}, Err: err}
}

func WithLocation(err error, file string, line, column int) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		if e.Diag.File == "" && file != "" {
			e.Diag.File = file
		}
		if e.Diag.Line == 0 && line != 0 {
			e.Diag.Line = line
		}
		if e.Diag.Column == 0 && column != 0 {
			e.Diag.Column = column
		}
		return e
	}
	return &Error{Diag: Diagnostic{File: file, Line: line, Column: column, Message: err.Error()}, Err: err}
}

func WithMethod(err error, iface, method string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		if e.Diag.Interface == "" && iface != "" {
			e.Diag.Interface = iface
		}
		if e.Diag.Method == "" && method != "" {
			e.Diag.Method = method
		}
		return e
	}
	return &Error{Diag: Diagnostic{Interface: iface, Method: method, Message: err.Error()}, Err: err}
}
