package field

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm/clause"
)

// ScanValuer interface for Field
type ScanValuer interface {
	Scan(src interface{}) error   // sql.Scanner
	Value() (driver.Value, error) // driver.Valuer
}

func NewGenerics[T any](e expr) GenericsField[T] {
	return GenericsField[T]{e}
}

// GenericsField a generics field struct
type GenericsField[T any] struct{ expr }

// Eq judge equal
func (field GenericsField[T]) Eq(value T) Expr {
	return expr{e: clause.Eq{Column: field.RawExpr(), Value: value}}
}

// Neq judge not equal
func (field GenericsField[T]) Neq(value T) Expr {
	return expr{e: clause.Neq{Column: field.RawExpr(), Value: value}}
}

// In ...
func (field GenericsField[T]) In(values ...T) Expr {
	return expr{e: clause.IN{Column: field.RawExpr(), Values: field.toSlice(values...)}}
}

// Gt ...
func (field GenericsField[T]) Gt(value T) Expr {
	return expr{e: clause.Gt{Column: field.RawExpr(), Value: value}}
}

// Gte ...
func (field GenericsField[T]) Gte(value T) Expr {
	return expr{e: clause.Gte{Column: field.RawExpr(), Value: value}}
}

// Lt ...
func (field GenericsField[T]) Lt(value T) Expr {
	return expr{e: clause.Lt{Column: field.RawExpr(), Value: value}}
}

// Lte ...
func (field GenericsField[T]) Lte(value T) Expr {
	return expr{e: clause.Lte{Column: field.RawExpr(), Value: value}}
}

// Like ...
func (field GenericsField[T]) Like(value string) Expr {
	return expr{e: clause.Like{Column: field.RawExpr(), Value: value}}
}

// NotLike ...
func (field GenericsField[T]) NotLike(value string) Expr {
	return expr{e: clause.Not(field.Like(value).expression())}
}

// Value ...
func (field GenericsField[T]) Value(value T) AssignExpr {
	return field.value(value)
}

// Sum ...
func (field GenericsField[T]) Sum() GenericsField[T] {
	return GenericsField[T]{field.sum()}
}

// IfNull ...
func (field GenericsField[T]) IfNull(value T) Expr {
	return field.ifNull(value)
}

func (field GenericsField[T]) toSlice(values ...T) []interface{} {
	slice := make([]interface{}, len(values))
	for i, v := range values {
		slice[i] = v
	}
	return slice
}

func NewGenericsInt[T any](e expr) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: NewGenerics[T](e)}
}

// GenericsInt int type field
type GenericsInt[T any] struct {
	GenericsField[T]
}

// NotIn ...
func (field GenericsInt[T]) NotIn(values ...T) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Between ...
func (field GenericsInt[T]) Between(left T, right T) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field GenericsInt[T]) NotBetween(left T, right T) Expr {
	return Not(field.Between(left, right))
}

// Add ...
func (field GenericsInt[T]) Add(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.add(value)}}
}

// Sub ...
func (field GenericsInt[T]) Sub(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.sub(value)}}
}

// Mul ...
func (field GenericsInt[T]) Mul(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.mul(value)}}
}

// Div ...
func (field GenericsInt[T]) Div(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.div(value)}}
}

// Mod ...
func (field GenericsInt[T]) Mod(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.mod(value)}}
}

// FloorDiv ...
func (field GenericsInt[T]) FloorDiv(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.floorDiv(value)}}
}

// RightShift ...
func (field GenericsInt[T]) RightShift(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.rightShift(value)}}
}

// LeftShift ...
func (field GenericsInt[T]) LeftShift(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.leftShift(value)}}
}

// BitXor ...
func (field GenericsInt[T]) BitXor(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.bitXor(value)}}
}

// BitAnd ...
func (field GenericsInt[T]) BitAnd(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.bitAnd(value)}}
}

// BitOr ...
func (field GenericsInt[T]) BitOr(value T) GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.bitOr(value)}}
}

// BitFlip ...
func (field GenericsInt[T]) BitFlip() GenericsInt[T] {
	return GenericsInt[T]{GenericsField: GenericsField[T]{field.bitFlip()}}
}

// Zero set zero value
func (field GenericsInt[T]) Zero() AssignExpr {
	return field.value(0)
}

func NewGenericsString[T any](e expr) GenericsString[T] {
	return GenericsString[T]{GenericsField: NewGenerics[T](e)}
}

// GenericsString string type field
type GenericsString[T any] struct {
	GenericsField[T]
}

// Between ...
func (field GenericsString[T]) Between(left T, right T) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field GenericsString[T]) NotBetween(left T, right T) Expr {
	return Not(field.Between(left, right))
}

// NotIn ...
func (field GenericsString[T]) NotIn(values ...T) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Regexp ...
func (field GenericsString[T]) Regexp(value string) Expr {
	return field.regexp(value)
}

// NotRegxp ...
func (field GenericsString[T]) NotRegxp(value string) Expr {
	return expr{e: clause.Not(field.Regexp(value).expression())}
}

// Zero ...
func (field GenericsString[T]) Zero() AssignExpr {
	return field.value("")
}

// FindInSet equal to FIND_IN_SET(field_name, input_string_list)
func (field GenericsString[T]) FindInSet(targetList string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{field.RawExpr(), targetList}}}
}

// FindInSetWith equal to FIND_IN_SET(input_string, field_name)
func (field GenericsString[T]) FindInSetWith(target string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{target, field.RawExpr()}}}
}

// Replace ...
func (field GenericsString[T]) Replace(from, to string) GenericsString[T] {
	return GenericsString[T]{GenericsField: GenericsField[T]{expr{e: clause.Expr{SQL: "REPLACE(?,?,?)", Vars: []interface{}{field.RawExpr(), from, to}}}}}
}

// Concat ...
func (field GenericsString[T]) Concat(before, after string) GenericsString[T] {
	switch {
	case before != "" && after != "":
		return GenericsString[T]{GenericsField: GenericsField[T]{expr{e: clause.Expr{SQL: "CONCAT(?,?,?)", Vars: []interface{}{before, field.RawExpr(), after}}}}}
	case before != "":
		return GenericsString[T]{GenericsField: GenericsField[T]{expr{e: clause.Expr{SQL: "CONCAT(?,?)", Vars: []interface{}{before, field.RawExpr()}}}}}
	case after != "":
		return GenericsString[T]{GenericsField: GenericsField[T]{expr{e: clause.Expr{SQL: "CONCAT(?,?)", Vars: []interface{}{field.RawExpr(), after}}}}}
	default:
		return field
	}
}

// SubstringIndex SUBSTRING_INDEX
// https://dev.mysql.com/doc/refman/8.0/en/functions.html#function_substring-index
func (field GenericsString[T]) SubstringIndex(delim string, count int) GenericsString[T] {
	return GenericsString[T]{GenericsField: GenericsField[T]{expr{e: clause.Expr{
		SQL:  fmt.Sprintf("SUBSTRING_INDEX(?,%q,%d)", delim, count),
		Vars: []interface{}{field.RawExpr()},
	}}}}
}
