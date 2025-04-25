package field

import (
	"fmt"

	"gorm.io/gorm/clause"
)

// String string type field
type String = Chars[string]

// Bytes []byte type field
type Bytes = Chars[[]byte]

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

// NotRegexp ...
func (field Chars[T]) NotRegexp(value string) Expr {
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

// Substring https://dev.mysql.com/doc/refman/8.4/en/string-functions.html#function_substring
func (field Chars[T]) Substring(params ...int) Chars[T] {
	if len(params) == 0 {
		return field
	}
	if len(params) == 1 {
		return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{
			SQL:  fmt.Sprintf("SUBSTRING(?,%d)", params[0]),
			Vars: []interface{}{field.RawExpr()},
		}}}}
	}
	return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{
		SQL:  fmt.Sprintf("SUBSTRING(?,%d,%d)", params[0], params[1]),
		Vars: []interface{}{field.RawExpr()},
	}}}}
}

// Substr SUBSTR is a synonym for SUBSTRING
// https://dev.mysql.com/doc/refman/8.4/en/string-functions.html#function_substring
func (field Chars[T]) Substr(params ...int) Chars[T] {
	if len(params) == 0 {
		return field
	}
	if len(params) == 1 {
		return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{
			SQL:  fmt.Sprintf("SUBSTR(?,%d)", params[0]),
			Vars: []interface{}{field.RawExpr()},
		}}}}
	}
	return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{
		SQL:  fmt.Sprintf("SUBSTR(?,%d,%d)", params[0], params[1]),
		Vars: []interface{}{field.RawExpr()},
	}}}}
}

// SubstringIndex SUBSTRING_INDEX
// https://dev.mysql.com/doc/refman/8.0/en/functions.html#function_substring-index
func (field Chars[T]) SubstringIndex(delim string, count int) Chars[T] {
	return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{
		SQL:  fmt.Sprintf("SUBSTRING_INDEX(?,%q,%d)", delim, count),
		Vars: []interface{}{field.RawExpr()},
	}}}}
}

// Field ...
func (field Chars[T]) Field(values ...T) Chars[T] {
	return newChars[T](field.field(values))
}

// Lower converts a string to lower-case.
func (field Chars[T]) Lower() Chars[T] {
	return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{SQL: "LOWER(?)", Vars: []interface{}{field.RawExpr()}}}}}
}

// Upper converts a string to upper-case.
func (field Chars[T]) Upper() Chars[T] {
	return Chars[T]{genericsField: genericsField[T]{expr{e: clause.Expr{SQL: "UPPER(?)", Vars: []interface{}{field.RawExpr()}}}}}
}
