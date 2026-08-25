package gen

import (
	"encoding/json"
	"errors"
	"io"

	"gorm.io/gen/internal/diagnostic"
)

// WriteDiagnosticJSON writes err as one newline-terminated JSON value.
// Structured generator diagnostics retain their code, source location, and
// hint; other errors are emitted as an object containing an "error" field.
// A nil error is written as JSON null.
func WriteDiagnosticJSON(w io.Writer, err error) error {
	if err == nil {
		_, writeErr := w.Write([]byte("null\n"))
		return writeErr
	}
	var de *diagnostic.Error
	if errors.As(err, &de) {
		b, marshalErr := json.MarshalIndent(de, "", "  ")
		if marshalErr != nil {
			return marshalErr
		}
		b = append(b, '\n')
		_, writeErr := w.Write(b)
		return writeErr
	}
	b, marshalErr := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	b = append(b, '\n')
	_, writeErr := w.Write(b)
	return writeErr
}
