package field

import (
	"gorm.io/gorm/clause"
)

// newGenerics create new generic field type
func newGenerics[T any](e expr) genericsField[T] {
	return genericsField[T]{e}
}

// genericsField a generics field struct
// serving as a base field type, offers a suite of fundamental methods/functions for database operations."
type genericsField[T any] struct{ expr }

// Eq judge equal
func (field genericsField[T]) Eq(value T) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq judge not equal
func (field genericsField[T]) Neq(value T) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field genericsField[T]) In(values ...T) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// NotIn ...
func (field genericsField[T]) NotIn(values ...T) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Gt ...
func (field genericsField[T]) Gt(value T) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte ...
func (field genericsField[T]) Gte(value T) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt ...
func (field genericsField[T]) Lt(value T) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte ...
func (field genericsField[T]) Lte(value T) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// Like ...
func (field genericsField[T]) Like(value string) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field genericsField[T]) NotLike(value string) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Value ...
func (field genericsField[T]) Value(value T) AssignExpr {
	return field.value(value)
}

// Sum ...
func (field genericsField[T]) Sum() genericsField[T] {
	return genericsField[T]{field.sum()}
}

// IfNull ...
func (field genericsField[T]) IfNull(value T) Expr {
	return field.ifNull(value)
}

// Field ...
func (field genericsField[T]) Field(value []interface{}) Expr {
	return field.field(value)
}

func (field genericsField[T]) toSlice(values ...T) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}
