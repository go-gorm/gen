package field

import (
	"context"
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ValuerType delays conversion of a serializer value until GORM builds the statement.
type ValuerType struct {
	// Column identifies the schema field whose serializer should perform the conversion.
	Column string
	// Value is the application value passed to the schema serializer.
	Value schema.SerializerValuerInterface
}

// GormValue converts Value with the serializer registered for Column.
// Conversion errors are recorded on db and the returned expression still uses
// a bind variable so statement construction remains parameterized.
func (v ValuerType) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) {
	stmt := db.Statement.Schema
	field := stmt.LookUpField(v.Column)
	newValue, err := v.Value.Value(ctx, field, reflect.ValueOf(v.Value), v.Value)
	_ = db.AddError(err)
	return clause.Expr{SQL: "?", Vars: []interface{}{newValue}}
}

// SerializerField is a typed field whose values are converted by a GORM serializer.
type SerializerField[T any] struct{ expr }

// Serializer preserves the original field type for values that implement
// GORM's serializer interface.
type Serializer = SerializerField[schema.SerializerValuerInterface]

// Eq judge equal
func (field SerializerField[T]) Eq(value T) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Neq judge not equal
func (field SerializerField[T]) Neq(value T) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// In ...
func (field SerializerField[T]) In(values ...T) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// Gt ...
func (field SerializerField[T]) Gt(value T) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Gte ...
func (field SerializerField[T]) Gte(value T) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Lt ...
func (field SerializerField[T]) Lt(value T) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Lte ...
func (field SerializerField[T]) Lte(value T) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Like ...
func (field SerializerField[T]) Like(value T) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: field.wrap(value)}}
}

// Value ...
func (field SerializerField[T]) Value(value T) AssignExpr {
	return field.value(field.wrap(value))
}

// Sum ...
func (field SerializerField[T]) Sum() Number[float64] {
	return newNumber[float64](field.sum())
}

// IfNull ...
func (field SerializerField[T]) IfNull(value T) Expr {
	return field.ifNull(field.wrap(value))
}

func (field SerializerField[T]) wrap(value T) interface{} {
	column := field.ColumnName().String()
	if valuer, ok := any(value).(schema.SerializerValuerInterface); ok {
		return ValuerType{Column: column, Value: valuer}
	}
	return schemaSerializerValue{Column: column, Value: value}
}

func (field SerializerField[T]) toSlice(values ...T) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = field.wrap(v)
	}
	return slice
}

type schemaSerializerValue struct {
	Column string
	Value  interface{}
}

func (v schemaSerializerValue) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	expr := clause.Expr{SQL: "?", Vars: []interface{}{v.Value}}
	if db == nil || db.Statement == nil || db.Statement.Schema == nil {
		if db != nil {
			_ = db.AddError(fmt.Errorf("field: schema is unavailable for serialized column %q", v.Column))
		}
		return expr
	}

	schemaField := db.Statement.Schema.LookUpField(v.Column)
	if schemaField == nil {
		_ = db.AddError(fmt.Errorf("field: serialized column %q was not found in schema", v.Column))
		return expr
	}

	valuer := schemaField.Serializer
	if valuer == nil {
		_ = db.AddError(fmt.Errorf("field: column %q has no schema serializer", v.Column))
		return expr
	}

	destination := db.Statement.ReflectValue
	if !destination.IsValid() {
		if db.Statement.Schema.ModelType == nil {
			_ = db.AddError(fmt.Errorf("field: schema model type is unavailable for serialized column %q", v.Column))
			return expr
		}
		destination = reflect.New(db.Statement.Schema.ModelType)
	}
	newValue, err := valuer.Value(ctx, schemaField, destination, v.Value)
	_ = db.AddError(err)
	expr.Vars[0] = newValue
	return expr
}
