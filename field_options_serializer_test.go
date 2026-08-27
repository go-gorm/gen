package gen

import (
	"testing"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/model"
)

func TestSerializerFieldGenTypeFromPublicOptions(t *testing.T) {
	tests := []struct {
		name string
		opt  model.ModifyFieldOpt
	}{
		{
			name: "structured gorm tag",
			opt: FieldGORMTag("photos", func(tag field.GormTag) field.GormTag {
				tag.Set("serializer", "json")
				return tag
			}),
		},
		{
			name: "raw gorm tag",
			opt:  FieldNewTag("photos", field.Tag{field.TagKeyGorm: "serializer:json"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelField := &model.Field{
				Name:       "Photos",
				Type:       "[]*string",
				ColumnName: "photos",
				Tag:        field.Tag{},
				GORMTag:    field.GormTag{},
			}
			if got := tt.opt(modelField).GenType(); got != "SerializerField[[]*string]" {
				t.Fatalf("GenType() = %q, want SerializerField[[]*string]", got)
			}
		})
	}
}
