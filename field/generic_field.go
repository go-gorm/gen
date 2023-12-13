package field

import (
	"gorm.io/gorm/clause"
)

// GenericField a standard generic field struct
type GenericField[T any] struct{ expr }

// Eq judge equal
func (field GenericField[T]) Eq(value T) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq judge not equal
func (field GenericField[T]) Neq(value T) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field GenericField[T]) In(values ...T) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field GenericField[T]) NotIn(values ...T) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Gt ...
func (field GenericField[T]) Gt(value T) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte ...
func (field GenericField[T]) Gte(value T) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt ...
func (field GenericField[T]) Lt(value T) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte ...
func (field GenericField[T]) Lte(value T) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// Like ...
func (field GenericField[T]) Like(value T) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// Value ...
func (field GenericField[T]) Value(value T) AssignExpr {
	return field.value(value)
}

// Sum ...
func (field GenericField[T]) Sum() GenericField[T] {
	return GenericField[T]{field.sum()}
}

// IfNull ...
func (field GenericField[T]) IfNull(value T) Expr {
	return field.ifNull(value)
}

// Field ...
func (field GenericField[T]) Field(value ...T) Expr {
	return field.field(field.toSlice(value...))
}

// ToInt convert to int field
func (field GenericField[T]) ToInt() Int {
	return Int{field.expr}
}

// ToFloat64 convert to float64 field
func (field GenericField[T]) ToFloat64() Float64 {
	return Float64{field.expr}
}

// ToString convert to string field
func (field GenericField[T]) ToString() String {
	return String{field.expr}
}

func (field GenericField[T]) toSlice(values ...T) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
