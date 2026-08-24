package gen

import (
	"strings"
	"testing"
)

func TestFormatUseAny(t *testing.T) {
	src := `package query

// value swaps interface{} spellings (comment mention).
func value(v interface{}, list []interface{}, m map[string]interface{}) interface{} {
	var x interface{ M() } = nil
	_ = x
	return v
}
`
	g := NewGenerator(Config{OutPath: t.TempDir()})

	t.Run("default keeps interface{} spelling", func(t *testing.T) {
		out, err := g.format("value.go", []byte(src))
		if err != nil {
			t.Fatal(err)
		}
		got := string(out)
		if !strings.Contains(got, "v interface{}") || !strings.Contains(got, "map[string]interface{}") {
			t.Errorf("default output should keep the interface{} spelling:\n%s", got)
		}
		if strings.Contains(got, "v any") || strings.Contains(got, "map[string]any") {
			t.Errorf("default output must not rewrite to any:\n%s", got)
		}
	})

	t.Run("UseAny rewrites empty interfaces only", func(t *testing.T) {
		g.UseAny = true
		out, err := g.format("value.go", []byte(src))
		if err != nil {
			t.Fatal(err)
		}
		got := string(out)
		if strings.Contains(got, "interface{}") {
			t.Errorf("empty interface should be rewritten to any:\n%s", got)
		}
		for _, want := range []string{"v any", "[]any", "map[string]any", ") any {"} {
			if !strings.Contains(got, want) {
				t.Errorf("rewritten output missing %q:\n%s", want, got)
			}
		}
		if !strings.Contains(got, "interface{ M() }") {
			t.Errorf("non-empty interface must survive the rewrite:\n%s", got)
		}
	})
}
