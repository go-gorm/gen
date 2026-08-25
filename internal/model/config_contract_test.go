package model

import (
	"reflect"
	"testing"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

func TestConfigPreprocessAndNames(t *testing.T) {
	modify := ModifyFieldOpt(func(input *Field) *Field { return input })
	filter := FilterFieldOpt(func(input *Field) *Field { return input })
	create := CreateFieldOpt(func(input *Field) *Field { return input })
	method := AddMethodOpt(func() []interface{} { return []interface{}{"method"} })

	cfg := (&Config{
		ModelPkg:    "/tmp/project/models",
		TablePrefix: "pre_",
		TableName:   "users",
		ModelName:   "User",
		ModelOpts:   []Option{modify, filter, create, method},
		NameStrategy: NameStrategy{
			TableNameNS: func(string) string { return "archived_users" },
			ModelNameNS: func(string) string { return "ArchivedUser" },
			FileNameNS:  func(string) string { return "user_archive" },
		},
	}).Preprocess()

	if cfg.ModelPkg != "models" {
		t.Fatalf("ModelPkg = %q, want models", cfg.ModelPkg)
	}
	if len(cfg.ModifyOpts) != 1 || len(cfg.FilterOpts) != 1 || len(cfg.CreateOpts) != 1 || len(cfg.MethodOpts) != 1 {
		t.Fatalf("unexpected option split: modify=%d filter=%d create=%d method=%d",
			len(cfg.ModifyOpts), len(cfg.FilterOpts), len(cfg.CreateOpts), len(cfg.MethodOpts))
	}

	tableName, structName, fileName := cfg.GetNames()
	if tableName != "pre_archived_users" || structName != "ArchivedUser" || fileName != "user_archive" {
		t.Fatalf("GetNames() = (%q, %q, %q)", tableName, structName, fileName)
	}
	if got := cfg.GetModelMethods(); !reflect.DeepEqual(got, []interface{}{"method"}) {
		t.Fatalf("GetModelMethods() = %#v", got)
	}
}

func TestConfigDefaultsAndSchemaName(t *testing.T) {
	cfg := (&Config{}).Preprocess()
	if cfg.ModelPkg != DefaultModelPkg {
		t.Fatalf("ModelPkg = %q, want %q", cfg.ModelPkg, DefaultModelPkg)
	}
	if got := (*Config)(nil).GetModelMethods(); got != nil {
		t.Fatalf("nil config methods = %#v, want nil", got)
	}
	if got := (*Config)(nil).GetSchemaName(nil); got != "" {
		t.Fatalf("nil config schema = %q, want empty", got)
	}

	cfg.SchemaNameOpts = []SchemaNameOpt{
		func(*gorm.DB) string { return "" },
		func(*gorm.DB) string { return "tenant" },
		func(*gorm.DB) string { return "ignored" },
	}
	if got := cfg.GetSchemaName(nil); got != "tenant" {
		t.Fatalf("GetSchemaName() = %q, want tenant", got)
	}
}

func TestGroupByColumn(t *testing.T) {
	index := migrator.Index{NameValue: "idx_tenant_user", ColumnList: []string{"tenant_id", "user_id"}}
	got := GroupByColumn([]gorm.Index{nil, index})
	if len(got) != 2 {
		t.Fatalf("GroupByColumn() columns = %d, want 2", len(got))
	}
	if got["tenant_id"][0].Priority != 1 || got["user_id"][0].Priority != 2 {
		t.Fatalf("unexpected priorities: tenant=%d user=%d", got["tenant_id"][0].Priority, got["user_id"][0].Priority)
	}
	if empty := GroupByColumn(nil); len(empty) != 0 {
		t.Fatalf("empty GroupByColumn() = %#v", empty)
	}
}

func TestFieldAndSQLBufferContracts(t *testing.T) {
	modelField := &Field{Name: "Select", Type: "*int64"}
	if got := modelField.EscapeKeyword(); got.Name != "Select_" {
		t.Fatalf("EscapeKeyword().Name = %q", got.Name)
	}
	if got := modelField.GenType(); got != "Int64" {
		t.Fatalf("GenType() = %q, want Int64", got)
	}

	relationField := &Field{Type: "[]User", Relation: &field.Relation{}}
	if got := relationField.GenType(); got != "[]User" {
		t.Fatalf("relation GenType() = %q", got)
	}

	var buffer SQLBuffer
	for _, b := range []byte("  SELECT\n\t*  FROM users") {
		buffer.WriteSQL(b)
	}
	if got, want := buffer.Dump(), " SELECT * FROM users"; got != want {
		t.Fatalf("Dump() = %q, want %q", got, want)
	}
	if buffer.Len() != 0 {
		t.Fatalf("Dump() did not reset buffer, len=%d", buffer.Len())
	}
}
