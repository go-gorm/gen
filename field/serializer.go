package field

import (
	"context"
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

// Serializer a standard field struct
type Serializer struct{ expr }

// Eq judge equal
func (field Serializer) Eq(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Neq judge not equal
func (field Serializer) Neq(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// In ...
func (field Serializer) In(values ...schema.SerializerValuerInterface) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// Gt ...
func (field Serializer) Gt(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Gte ...
func (field Serializer) Gte(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Lt ...
func (field Serializer) Lt(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Lte ...
func (field Serializer) Lte(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Like ...
func (field Serializer) Like(value schema.SerializerValuerInterface) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: ValuerType{Column: field.ColumnName().String(), Value: value}}}
}

// Value ...
func (field Serializer) Value(value schema.SerializerValuerInterface) AssignExpr {
	return field.value(ValuerType{Column: field.ColumnName().String(), Value: value})
}

// Sum ...
func (field Serializer) Sum() Number[float64] {
	return newNumber[float64](field.sum())
}

// IfNull ...
func (field Serializer) IfNull(value schema.SerializerValuerInterface) Expr {
	return field.ifNull(ValuerType{Column: field.ColumnName().String(), Value: value})
}

func (field Serializer) toSlice(values ...schema.SerializerValuerInterface) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = ValuerType{Column: field.ColumnName().String(), Value: v}
	}
	return slice
}
