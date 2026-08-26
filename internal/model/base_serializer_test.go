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
		{name: "structured serializer tag", field: &Field{Type: "map[string]int", GORMTag: field.GormTag{"serializer": {"json"}}}, want: "SerializerField[map[string]int]"},
		{name: "structured json tag", field: &Field{Type: "map[string]int", GORMTag: field.GormTag{"json": {"json"}}}, want: "SerializerField[map[string]int]"},
		{name: "raw gorm tag", field: &Field{Type: "[]*string", Tag: field.Tag{field.TagKeyGorm: "column:photos;serializer:json"}}, want: "SerializerField[[]*string]"},
		{name: "plain slice", field: &Field{Type: "slice"}, want: "Field"},
		{name: "json struct tag", field: &Field{Type: "slice", Tag: field.Tag{field.TagKeyJson: "photos"}}, want: "Field"},
		{name: "custom field type", field: &Field{Type: "[]*string", GORMTag: field.GormTag{"serializer": {"json"}}, CustomGenType: "String"}, want: "String"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.GenType(); got != tt.want {
				t.Fatalf("GenType() = %q, want %q", got, tt.want)
			}
		})
	}
}
