package parser

import (
	"os"
	"path/filepath"
	"testing"
)

const anyInterfaceSrc = `package diy

type TestAny interface {
	M1(v any) error
	M2(vs []any) error
	M3(m map[string]any) error
	M4(v interface{}) error
	M5(v *any) error
	M6(vs ...any) error
	M7() (any, error)
}
`

// a user-defined type named "any" shadows the universe alias and must be
// treated as an ordinary package type, not canonicalized
const anyShadowSrc = `package diy

type any struct{ X int }

type TestShadow interface {
	S(v any) error
}
`

func parseInterface(t *testing.T, src, name string) *InterfaceInfo {
	t.Helper()
	file := filepath.Join(t.TempDir(), "diy.go")
	if err := os.WriteFile(file, []byte(src), 0o400); err != nil {
		t.Fatal(err)
	}
	i := &InterfaceSet{}
	if _, err := i.getInterfaceFromFile(file, name, "gorm.io/x/diy", nil); err != nil {
		t.Fatal(err)
	}
	for idx, info := range i.Interfaces {
		if info.Name == name {
			return &i.Interfaces[idx]
		}
	}
	t.Fatalf("interface %s not found in fixture", name)
	return nil
}

func TestParseAnyAsCanonicalInterface(t *testing.T) {
	info := parseInterface(t, anyInterfaceSrc, "TestAny")
	if len(info.Methods) != 7 {
		t.Fatalf("unexpected method count: %d", len(info.Methods))
	}

	shapes := []struct {
		method       string
		wantType     string
		wantArray    bool
		wantPointer  bool
		wantVariadic bool
	}{
		{method: "M1", wantType: "interface{}"},
		{method: "M2", wantType: "interface{}", wantArray: true},
		{method: "M3", wantType: "map[string]interface{}"},
		{method: "M4", wantType: "interface{}"},
		{method: "M5", wantType: "interface{}", wantPointer: true},
		{method: "M6", wantType: "interface{}", wantArray: true, wantVariadic: true},
		{method: "M7", wantType: "interface{}"},
	}
	for _, s := range shapes {
		var m *Method
		for _, mm := range info.Methods {
			if mm.MethodName == s.method {
				m = mm
				break
			}
		}
		if m == nil {
			t.Fatalf("method %s not found", s.method)
		}
		if s.method == "M7" {
			if len(m.Result) != 2 || m.Result[0].Type != s.wantType || !m.Result[1].IsError() {
				t.Errorf("M7: result = %+v, want (interface{}, error)", m.Result)
			}
			continue
		}
		if len(m.Params) != 1 {
			t.Fatalf("%s: unexpected param count: %d", s.method, len(m.Params))
		}
		p := m.Params[0]
		if p.Type != s.wantType {
			t.Errorf("%s: param type = %q, want %q", s.method, p.Type, s.wantType)
		}
		if p.IsArray != s.wantArray {
			t.Errorf("%s: IsArray = %v, want %v", s.method, p.IsArray, s.wantArray)
		}
		if p.IsPointer != s.wantPointer {
			t.Errorf("%s: IsPointer = %v, want %v", s.method, p.IsPointer, s.wantPointer)
		}
		if p.IsVariadic != s.wantVariadic {
			t.Errorf("%s: IsVariadic = %v, want %v", s.method, p.IsVariadic, s.wantVariadic)
		}
	}
}

func TestParseShadowedAnyStaysUserType(t *testing.T) {
	info := parseInterface(t, anyShadowSrc, "TestShadow")
	p := info.Methods[0].Params[0]
	if p.Type != "any" || p.Package != "UNDEFINED" {
		t.Errorf("shadowed any: Type = %q, Package = %q, want Type = %q, Package = %q",
			p.Type, p.Package, "any", "UNDEFINED")
	}
}

func TestParamIsInterface(t *testing.T) {
	for _, c := range []struct {
		p    Param
		want bool
	}{
		{Param{Type: "interface{}"}, true},
		{Param{Type: "any"}, true},
		{Param{Package: "diy", Type: "any"}, false},
		{Param{Type: "map[string]interface{}"}, false},
	} {
		if got := c.p.IsInterface(); got != c.want {
			t.Errorf("IsInterface(%+v) = %v, want %v", c.p, got, c.want)
		}
	}
}
