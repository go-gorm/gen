package diagnostic

import "testing"

func TestDefaultMessageKnownCodes(t *testing.T) {
	cases := []struct {
		code string
		want string
	}{
		{CodeSQLIncomplete, "incomplete SQL"},
		{CodeSQLVar, "variable parse error"},
		{CodeTemplateParse, "template parse error"},
		{CodeSQLBuild, "build SQL error"},
	}
	for _, c := range cases {
		if got := DefaultMessage(c.code); got != c.want {
			t.Fatalf("code=%s got=%q want=%q", c.code, got, c.want)
		}
		if got := DefaultHint(c.code); got == "" {
			t.Fatalf("code=%s expected non-empty hint", c.code)
		}
	}
}

func TestNew_FillsDefaultsWhenEmpty(t *testing.T) {
	e := New(CodeSQLIncomplete, "")
	if e.Diag.Message == "" || e.Diag.Message != DefaultMessage(CodeSQLIncomplete) {
		t.Fatalf("unexpected message: %q", e.Diag.Message)
	}
	if e.Diag.Hint == "" {
		t.Fatalf("expected hint")
	}
}

func TestWrap_FillsDefaultsWhenEmpty(t *testing.T) {
	e := Wrap(assertErr{}, CodeSQLBuild, "")
	if e.Diag.Message == "" || e.Diag.Message != DefaultMessage(CodeSQLBuild) {
		t.Fatalf("unexpected message: %q", e.Diag.Message)
	}
	if e.Diag.Hint == "" {
		t.Fatalf("expected hint")
	}
}

type assertErr struct{}

func (assertErr) Error() string { return "x" }
