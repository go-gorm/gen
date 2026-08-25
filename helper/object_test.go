package helper

import (
	"testing"

	"gorm.io/gen/field"
)

type testObject struct {
	structName string
	fields     []Field
}

func (testObject) TableName() string        { return "objects" }
func (o testObject) StructName() string     { return o.structName }
func (testObject) FileName() string         { return "object" }
func (testObject) ImportPkgPaths() []string { return nil }
func (o testObject) Fields() []Field        { return o.fields }

type testObjectField struct {
	name     string
	typeName string
}

func (f testObjectField) Name() string     { return f.name }
func (f testObjectField) Type() string     { return f.typeName }
func (testObjectField) ColumnName() string { return "column" }
func (testObjectField) GORMTag() string    { return "" }
func (testObjectField) JSONTag() string    { return "" }
func (testObjectField) Tag() field.Tag     { return nil }
func (testObjectField) Comment() string    { return "" }

func TestCheckObject(t *testing.T) {
	tests := []struct {
		name    string
		object  Object
		wantErr bool
	}{
		{name: "empty struct name", object: testObject{}, wantErr: true},
		{
			name: "empty field name",
			object: testObject{structName: "User", fields: []Field{
				testObjectField{typeName: "string"},
			}},
			wantErr: true,
		},
		{
			name: "empty field type",
			object: testObject{structName: "User", fields: []Field{
				testObjectField{name: "Name"},
			}},
			wantErr: true,
		},
		{
			name: "valid",
			object: testObject{structName: "User", fields: []Field{
				testObjectField{name: "Name", typeName: "string"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckObject(tt.object)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
