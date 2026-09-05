package gen

import (
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/gen/internal/generate"
	"gorm.io/gen/internal/parser"
)

func TestGenerateModelFileWrapsCustomTemplateError(t *testing.T) {
	g := NewGenerator(Config{
		OutPath: filepath.Join(t.TempDir(), "query"),
		Mode:    WithDefaultQuery,
	})
	g.models = map[string]*generate.QueryStructMeta{
		"User": {
			Generated:       true,
			FileName:        "users",
			TableName:       "users",
			ModelStructName: "User",
			QueryStructName: "user",
			S:               "u",
			StructInfo:      parser.Param{Type: "User"},
			CustomTemplates: []string{"// ok {{.ModelStructName}}", "{{invalid"},
		},
	}

	err := g.generateModelFile()
	if err == nil {
		t.Fatal("generateModelFile() = nil, want error for invalid custom template")
	}
	msg := err.Error()
	if !strings.Contains(msg, "custom template #2") {
		t.Errorf("error message missing custom template index: %s", msg)
	}
	if !strings.Contains(msg, "User") {
		t.Errorf("error message missing model name: %s", msg)
	}
}
