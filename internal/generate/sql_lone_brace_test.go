package generate

import (
	"errors"
	"strings"
	"testing"

	"gorm.io/gen/internal/diagnostic"
)

// TestSQLStateCheckAndSplitRejectsLoneBrace covers the silently-dropped lone '{'
// case (and the unaffected normal paths) of sqlStateCheckAndSplit.
func TestSQLStateCheckAndSplitRejectsLoneBrace(t *testing.T) {
	testcases := []struct {
		name string
		SQL  string
		// wantErr true: expect a SQL_INCOMPLETE diagnostic whose message contains
		// wantMsg; wantErr false: expect no error and the literal '{' must survive.
		wantErr               bool
		wantMsg               string
		wantLiteralBraceInSQL bool
	}{
		{
			name:    "lone brace before word",
			SQL:     "select {x from t",
			wantErr: true,
			wantMsg: "'{'",
		},
		{
			name:    "lone brace after template",
			SQL:     "select * from @@table {{where}} id>@id{{end}} {or id<@id",
			wantErr: true,
			wantMsg: "'{'",
		},
		{
			name:    "escaped lone brace is rejected too",
			SQL:     "select \\{ from t",
			wantErr: true,
			wantMsg: "'{'",
		},
		{
			name:                  "brace inside single-quoted string is literal",
			SQL:                   "select * from t where name = '{'",
			wantLiteralBraceInSQL: true,
		},
		{
			name:                  "brace inside double-quoted string is literal",
			SQL:                   `select * from t where name = "{"`,
			wantLiteralBraceInSQL: true,
		},
		{
			name: "normal template still passes",
			SQL:  "select * from @@table {{where}} id>@id{{end}}",
		},
		{
			name:    "brace at end of SQL is incomplete",
			SQL:     "select * from t {",
			wantErr: true,
			wantMsg: "incomplete SQL",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			m := m()
			m.SQLString = tc.SQL
			err := m.sqlStateCheckAndSplit()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("SQL %q: expected error, got nil (Section=%v)", tc.SQL, m.Section.members)
				}
				var de *diagnostic.Error
				if !errors.As(err, &de) {
					t.Fatalf("expected diagnostic error, got %T: %v", err, err)
				}
				if de.Diag.Code != diagnostic.CodeSQLIncomplete {
					t.Fatalf("unexpected code: %s", de.Diag.Code)
				}
				if !strings.Contains(de.Diag.Message, tc.wantMsg) {
					t.Fatalf("message mismatch: want contains %q, got %q", tc.wantMsg, de.Diag.Message)
				}
				return
			}
			if err != nil {
				t.Fatalf("SQL %q: unexpected error: %v", tc.SQL, err)
			}
			if !tc.wantLiteralBraceInSQL {
				return
			}
			// literal '{' inside quotes must survive into the SQL section
			var got strings.Builder
			for _, mem := range m.Section.members {
				got.WriteString(mem.Value)
			}
			if !strings.Contains(got.String(), "{") {
				t.Fatalf("SQL %q: literal '{' lost in members: %v", tc.SQL, m.Section.members)
			}
		})
	}
}

// TestBuildDIYMethod_ReturnsDiagnosticOnLoneBrace verifies the end-to-end path:
// the diagnostic carries the interface/method names and the offending snippet.
func TestBuildDIYMethod_ReturnsDiagnosticOnLoneBrace(t *testing.T) {
	src := `package dal

import "gorm.io/gen"

type UserMethods interface {
	// FindByOddSQL
	//
	// SELECT * FROM users {WHERE id=@id
	FindByOddSQL(id int) gen.T
}
`
	set := parseInterfaceSet(t, src)

	_, err := BuildDIYMethod(set, testMeta(), nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	var de *diagnostic.Error
	if !errors.As(err, &de) {
		t.Fatalf("expected diagnostic error, got %T: %v", err, err)
	}
	if de.Diag.Code != diagnostic.CodeSQLIncomplete {
		t.Fatalf("unexpected code: %s", de.Diag.Code)
	}
	if !strings.Contains(de.Diag.Message, "'{'") {
		t.Fatalf("message should mention the lone brace: %q", de.Diag.Message)
	}
	if de.Diag.Interface != "UserMethods" || de.Diag.Method != "FindByOddSQL" {
		t.Fatalf("unexpected iface/method: %s.%s", de.Diag.Interface, de.Diag.Method)
	}
	if !strings.Contains(de.Diag.Snippet, "SELECT * FROM users {WHERE id=@id") {
		t.Fatalf("snippet should contain the offending SQL: %q", de.Diag.Snippet)
	}
	if de.Diag.File == "" || de.Diag.Line == 0 {
		t.Fatalf("expected file/line, got file=%q line=%d", de.Diag.File, de.Diag.Line)
	}
}
