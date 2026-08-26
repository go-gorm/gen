package field_test

import (
	"testing"

	"gorm.io/gen/field"
	"gorm.io/gorm/schema"
)

func TestSerializerPublicMethodSignatures(t *testing.T) {
	legacy := field.NewSerializer("users", "photos")
	var legacyValue func(schema.SerializerValuerInterface) field.AssignExpr = legacy.Value
	if legacyValue == nil {
		t.Fatal("legacy Value method is nil")
	}

	typed := field.NewSerializerField[[]*string]("users", "photos")
	var typedValue func([]*string) field.AssignExpr = typed.Value
	if typedValue == nil {
		t.Fatal("typed Value method is nil")
	}
}
