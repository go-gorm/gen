package diagnostic

import (
	"strings"
	"testing"
)

func TestCodeFrame_Basic(t *testing.T) {
	src := []byte("a\nb\nc\n")
	out := CodeFrame(src, 2, 1, 1)
	if !strings.Contains(out, ">2 | b") {
		t.Fatalf("unexpected output:\n%s", out)
	}
	if !strings.Contains(out, "  | ^") {
		t.Fatalf("missing caret:\n%s", out)
	}
}

func TestCodeFrame_ClampsOutOfRange(t *testing.T) {
	src := []byte("a\nb\n")
	out := CodeFrame(src, 100, 100, 1)
	if !strings.Contains(out, ">2 | b") {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestCodeFrame_TabAlignment(t *testing.T) {
	src := []byte("\tX\n")
	out := CodeFrame(src, 1, 2, 0)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("unexpected output:\n%s", out)
	}
	caretLine := lines[1]
	if !strings.HasSuffix(caretLine, "^") {
		t.Fatalf("missing caret:\n%s", out)
	}
	i := strings.Index(caretLine, "| ")
	if i < 0 {
		t.Fatalf("unexpected caret line:\n%s", out)
	}
	after := strings.TrimSuffix(caretLine[i+2:], "^")
	leading := len(after)
	if leading != 4 {
		t.Fatalf("unexpected caret column: %d\n%s", leading, out)
	}
}

func TestCodeFrame_UTF8Clamping(t *testing.T) {
	src := []byte("你a\n")
	out := CodeFrame(src, 1, 2, 0)
	if !strings.Contains(out, "^") {
		t.Fatalf("missing caret:\n%s", out)
	}
}
