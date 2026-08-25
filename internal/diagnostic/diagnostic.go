package diagnostic

import (
	"encoding/json"
	"fmt"
)

// Diagnostic is the machine-readable description of a generator failure.
type Diagnostic struct {
	// Code classifies the failure for programmatic handling.
	Code string `json:"code"`
	// Message is the human-readable summary.
	Message string `json:"message"`
	// File identifies the source file associated with the failure.
	File string `json:"file,omitempty"`
	// Line is the 1-based source line, when known.
	Line int `json:"line,omitempty"`
	// Column is the 1-based byte column, when known.
	Column int `json:"column,omitempty"`
	// Interface identifies the DIY interface being processed.
	Interface string `json:"interface,omitempty"`
	// Method identifies the DIY interface method being processed.
	Method string `json:"method,omitempty"`
	// Snippet contains a rendered source frame around the failure.
	Snippet string `json:"snippet,omitempty"`
	// Hint contains optional remediation guidance.
	Hint string `json:"hint,omitempty"`
}

// Error combines a structured Diagnostic with an optional underlying cause.
type Error struct {
	// Diag is the structured error description returned to users and tools.
	Diag Diagnostic
	// Err is the wrapped cause, if any.
	Err error
}

// Error returns a location- and code-prefixed diagnostic message.
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

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error { return e.Err }

// MarshalJSON serializes the diagnostic and includes the cause as a string.
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

// New constructs an Error and fills an empty message and hint from code.
func New(code, message string) *Error {
	if message == "" {
		message = DefaultMessage(code)
	}
	return &Error{Diag: Diagnostic{Code: code, Message: message, Hint: DefaultHint(code)}}
}

// Wrap annotates err with code and message.
// When err is already an Error, non-empty arguments replace its code and
// message while preserving its existing source and method context.
func Wrap(err error, code, message string) *Error {
	if err == nil {
		return New(code, message)
	}
	if message == "" {
		message = DefaultMessage(code)
	}
	if e, ok := err.(*Error); ok {
		if code != "" {
			e.Diag.Code = code
		}
		if message != "" {
			e.Diag.Message = message
		}
		if e.Diag.Hint == "" {
			e.Diag.Hint = DefaultHint(e.Diag.Code)
		}
		return e
	}
	return &Error{Diag: Diagnostic{Code: code, Message: message, Hint: DefaultHint(code)}, Err: err}
}

// WithLocation adds source coordinates to err without overwriting coordinates
// already present on a structured Error. A nil err remains nil.
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

// WithMethod adds interface and method context without overwriting context
// already present on a structured Error. A nil err remains nil.
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
