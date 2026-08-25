package gen

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"gorm.io/gen/internal/diagnostic"
)

func TestWriteDiagnosticJSON(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var output bytes.Buffer
		if err := WriteDiagnosticJSON(&output, nil); err != nil {
			t.Fatalf("WriteDiagnosticJSON() error = %v", err)
		}
		if got, want := output.String(), "null\n"; got != want {
			t.Fatalf("output = %q, want %q", got, want)
		}
	})

	t.Run("structured", func(t *testing.T) {
		input := diagnostic.New(diagnostic.CodeSQLBuild, "")
		input.Diag.File = "method.go"
		input.Diag.Line = 12
		input.Diag.Method = "FindByName"

		var output bytes.Buffer
		if err := WriteDiagnosticJSON(&output, input); err != nil {
			t.Fatalf("WriteDiagnosticJSON() error = %v", err)
		}

		var got diagnostic.Diagnostic
		if err := json.Unmarshal(output.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal output: %v", err)
		}
		if got.Code != diagnostic.CodeSQLBuild || got.File != "method.go" || got.Line != 12 || got.Method != "FindByName" {
			t.Fatalf("unexpected diagnostic: %+v", got)
		}
	})

	t.Run("plain error", func(t *testing.T) {
		var output bytes.Buffer
		if err := WriteDiagnosticJSON(&output, errors.New("boom")); err != nil {
			t.Fatalf("WriteDiagnosticJSON() error = %v", err)
		}

		var got map[string]string
		if err := json.Unmarshal(output.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal output: %v", err)
		}
		if got["error"] != "boom" {
			t.Fatalf("output error = %q, want boom", got["error"])
		}
	})
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) { return 0, w.err }

func TestWriteDiagnosticJSONPropagatesWriterError(t *testing.T) {
	want := errors.New("write failed")
	if err := WriteDiagnosticJSON(errorWriter{err: want}, errors.New("boom")); !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}
