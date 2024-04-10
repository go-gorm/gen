package field

import (
	"fmt"

	"gorm.io/gorm/clause"
)

// String string type field
type String struct{ Chars[string] }

// Bytes []byte type field
type Bytes struct{ Chars[[]byte] }

// ======================== string =======================

// newChars create field for chars
func newChars[T ~string | ~[]byte](e expr) Chars[T] {
	return Chars[T]{genericsField: newGenerics[T](e)}
}

// Chars string type field
type Chars[T ~string | ~[]byte] struct {
	genericsField[T]
}

// Between ...
func (field Chars[T]) Between(left T, right T) Expr {
	return field.between([]interface{}{left, right})
}

// NotBetween ...
func (field Chars[T]) NotBetween(left T, right T) Expr {
	return Not(field.Between(left, right))
}

// NotIn ...
func (field Chars[T]) NotIn(values ...T) Expr {
	return expr{e: clause.Not(field.In(values...).expression())}
}

// Regexp ...
func (field Chars[T]) Regexp(value string) Expr {
	return field.regexp(value)
}

// NotRegxp ...
func (field Chars[T]) NotRegxp(value string) Expr {
	return expr{e: clause.Not(field.Regexp(value).expression())}
}

// Zero ...
func (field Chars[T]) Zero() AssignExpr {
	return field.value("")
}

// FindInSet equal to FIND_IN_SET(field_name, input_string_list)
func (field Chars[T]) FindInSet(targetList string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{field.RawExpr(), targetList}}}
}

// FindInSetWith equal to FIND_IN_SET(input_string, field_name)
func (field Chars[T]) FindInSetWith(target string) Expr {
	return expr{e: clause.Expr{SQL: "FIND_IN_SET(?,?)", Vars: []interface{}{target, field.RawExpr()}}}
}

// Replace ...
func (field Chars[T]) Replace(from, to string) Chars[T] {
	return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{SQL: "REPLACE(?,?,?)", Vars: []interface{}{field.RawExpr(), from, to}}}}}
}

// Concat ...
func (field Chars[T]) Concat(before, after string) Chars[T] {
	switch {
	case before != "" && after != "":
		return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{SQL: "CONCAT(?,?,?)", Vars: []interface{}{before, field.RawExpr(), after}}}}}
	case before != "":
		return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{SQL: "CONCAT(?,?)", Vars: []interface{}{before, field.RawExpr()}}}}}
	case after != "":
		return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{SQL: "CONCAT(?,?)", Vars: []interface{}{field.RawExpr(), after}}}}}
	default:
		return field
	}
}

// SubstringIndex SUBSTRING_INDEX
// https://dev.mysql.com/doc/refman/8.0/en/functions.html#function_substring-index
func (field Chars[T]) SubstringIndex(delim string, count int) Chars[T] {
	return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{
		SQL:  fmt.Sprintf("SUBSTRING_INDEX(?,%q,%d)", delim, count),
		Vars: []interface{}{field.RawExpr()},
	}}}}
}
