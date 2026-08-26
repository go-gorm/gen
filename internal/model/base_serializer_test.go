package model

import (
	"testing"

	"gorm.io/gen/field"
)

func TestFieldGenTypeForSerializerTags(t *testing.T) {
	tests := []struct {
		name  string
		field *Field
		want  string
	}{
		{name: "legacy serializer type", field: &Field{Type: "serializer"}, want: "Serializer"},
		{name: "parsed schema serializer", field: &Field{Type: "slice", SchemaSerializer: true}, want: "Serialized"},
		{name: "structured serializer tag", field: &Field{Type: "map", GORMTag: field.GormTag{"serializer": {"json"}}}, want: "Serialized"},
		{name: "structured json tag", field: &Field{Type: "map", GORMTag: field.GormTag{"json": {"json"}}}, want: "Serialized"},
		{name: "raw gorm tag", field: &Field{Type: "slice", Tag: field.Tag{field.TagKeyGorm: "column:photos;serializer:json"}}, want: "Serialized"},
		{name: "plain slice", field: &Field{Type: "slice"}, want: "Field"},
		{name: "json struct tag", field: &Field{Type: "slice", Tag: field.Tag{field.TagKeyJson: "photos"}}, want: "Field"},
		{name: "custom field type", field: &Field{Type: "slice", SchemaSerializer: true, CustomGenType: "String"}, want: "String"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.GenType(); got != tt.want {
				t.Fatalf("GenType() = %q, want %q", got, tt.want)
			}
		})
	}
}
