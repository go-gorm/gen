package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/gen/field"
	"gorm.io/gen/helper"
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

type templateObjectField struct {
	name string
	typ  string
}

func (f templateObjectField) Name() string       { return f.name }
func (f templateObjectField) Type() string       { return f.typ }
func (f templateObjectField) ColumnName() string { return strings.ToLower(f.name) }
func (templateObjectField) GORMTag() string      { return "" }
func (templateObjectField) JSONTag() string      { return "" }
func (templateObjectField) Tag() field.Tag       { return nil }
func (templateObjectField) Comment() string      { return "" }

type templateObject struct {
	structName string
}

func (templateObject) TableName() string        { return "template_users" }
func (o templateObject) StructName() string     { return o.structName }
func (templateObject) FileName() string         { return "" }
func (templateObject) ImportPkgPaths() []string { return nil }
func (o templateObject) Fields() []helper.Field {
	return []helper.Field{templateObjectField{name: "Name", typ: "string"}}
}

func TestGenerateModelFromAppliesWithTemplate(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "query")
	g := NewGenerator(Config{
		OutPath: outPath,
		Mode:    WithDefaultQuery,
	})
	g.WithOpts(WithTemplate("// custom marker for {{.ModelStructName}}"))

	meta := g.GenerateModelFrom(templateObject{structName: "TemplateUser"})
	if meta == nil {
		t.Fatal("GenerateModelFrom() = nil")
	}
	if len(meta.CustomTemplates) != 1 {
		t.Fatalf("CustomTemplates = %v, want the WithTemplate entry", meta.CustomTemplates)
	}

	if panicValue := executePanic(g); panicValue != nil {
		t.Fatalf("Execute() panicked: %v", panicValue)
	}
	modelDir, err := g.getModelOutputPath()
	if err != nil {
		t.Fatalf("resolve model output path: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(modelDir, meta.FileName+".gen.go"))
	if err != nil {
		t.Fatalf("read generated model: %v", err)
	}
	if !strings.Contains(string(content), "// custom marker for TemplateUser") {
		t.Errorf("custom template output missing from generated model file:\n%s", content)
	}
}
