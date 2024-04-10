package field

import (
	"golang.org/x/exp/constraints"
	"gorm.io/gorm/clause"
)

// ======================== number =======================

// newNumber build number type field
func newNumber[T constraints.Integer | constraints.Float](e expr) Number[T] {
	return Number[T]{genericsField: newGenerics[T](e)}
}

// Number int type field
type Number[T constraints.Integer | constraints.Float] struct {
	genericsField[T]
}

// NotIn ...
func (field Number[T]) NotIn(values ...T) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field Number[T]) Between(left T, right T) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Number[T]) NotBetween(left T, right T) Expr {
	return Not(field.Between(left, right))
}

// Add ...
func (field Number[T]) Add(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.add(value)}}
}

// Sub ...
func (field Number[T]) Sub(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.sub(value)}}
}

// Mul ...
func (field Number[T]) Mul(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.mul(value)}}
}

// Div ...
func (field Number[T]) Div(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.div(value)}}
}

// Mod ...
func (field Number[T]) Mod(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.mod(value)}}
}

// FloorDiv ...
func (field Number[T]) FloorDiv(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.floorDiv(value)}}
}

// RightShift ...
func (field Number[T]) RightShift(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.rightShift(value)}}
}

// LeftShift ...
func (field Number[T]) LeftShift(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.leftShift(value)}}
}

// BitXor ...
func (field Number[T]) BitXor(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.bitXor(value)}}
}

// BitAnd ...
func (field Number[T]) BitAnd(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.bitAnd(value)}}
}

// BitOr ...
func (field Number[T]) BitOr(value T) Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.bitOr(value)}}
}

// BitFlip ...
func (field Number[T]) BitFlip() Number[T] {
	return Number[T]{genericsField: genericsField[T]{field.bitFlip()}}
}

// Zero set zero value
func (field Number[T]) Zero() AssignExpr {
	return field.value(0)
}
