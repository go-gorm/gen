package field

import (
	"golang.org/x/exp/constraints"
	"gorm.io/gorm/clause"
)

// Int type field
type Int = Number[int]

// Int8 type field
type Int8 = Number[int8]

// Int16 type field
type Int16 = Number[int16]

// Int32 type field
type Int32 = Number[int32]

// Int64 type field
type Int64 = Number[int64]

// Uint type field
type Uint = Number[uint]

// Uint8 type field
type Uint8 = Number[uint8]

// Uint16 type field
type Uint16 = Number[uint16]

// Uint32 type field
type Uint32 = Number[uint32]

// Uint64 type field
type Uint64 = Number[uint64]

// Float32 type field
type Float32 = Number[float32]

// Float64 type field
type Float64 = Number[float64]

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
	return newNumber[T](field.add(value))
}

// Sub ...
func (field Number[T]) Sub(value T) Number[T] {
	return newNumber[T](field.sub(value))
}

// Mul ...
func (field Number[T]) Mul(value T) Number[T] {
	return newNumber[T](field.mul(value))
}

// Div ...
func (field Number[T]) Div(value T) Number[T] {
	return newNumber[T](field.div(value))
}

// Mod ...
func (field Number[T]) Mod(value T) Number[T] {
	return newNumber[T](field.mod(value))
}

// FloorDiv ...
func (field Number[T]) FloorDiv(value T) Number[T] {
	return newNumber[T](field.floorDiv(value))
}

// RightShift ...
func (field Number[T]) RightShift(value T) Number[T] {
	return newNumber[T](field.rightShift(value))
}

// LeftShift ...
func (field Number[T]) LeftShift(value T) Number[T] {
	return newNumber[T](field.leftShift(value))
}

// BitXor ...
func (field Number[T]) BitXor(value T) Number[T] {
	return newNumber[T](field.bitXor(value))
}

// BitAnd ...
func (field Number[T]) BitAnd(value T) Number[T] {
	return newNumber[T](field.bitAnd(value))
}

// BitOr ...
func (field Number[T]) BitOr(value T) Number[T] {
	return newNumber[T](field.bitOr(value))
}

// BitFlip ...
func (field Number[T]) BitFlip() Number[T] {
	return newNumber[T](field.bitFlip())
}

// Floor ...
func (field Number[T]) Floor() Number[T] {
	return newNumber[T](field.floor())
}

// Zero set zero value
func (field Number[T]) Zero() AssignExpr {
	return field.value(0)
}

// Field ...
func (field Number[T]) Field(values ...T) Number[T] {
	return newNumber[T](field.field(values))
}
