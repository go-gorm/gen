package field

import (
	"database/sql/driver"
	"fmt"
	"golang.org/x/exp/constraints"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

var (
	// Star a symbol of "*"
	Star = NewAsterisk("")
	// ALL same with Star
	ALL = Star
)

// ScanValuer interface for Field
type ScanValuer interface {
	Scan(src interface{}) error   // sql.Scanner
	Value() (driver.Value, error) // driver.Valuer
}

// Option field option
type Option func(clause.Column) clause.Column

var (
	banColumnRaw Option = func(col clause.Column) clause.Column {
		col.Raw = false
		return col
	}
)

// ======================== generic field =======================

// NewField create new field
func NewField(table, column string, opts ...Option) Field {
	return Field(genericsField[driver.Valuer]{expr{col: toColumn(table, column, opts...)}})
}

// NewUnsafeFieldRaw create new field by native sql
//
// Warning: Using NewUnsafeFieldRaw with raw SQL exposes your application to SQL injection vulnerabilities.
// Always validate/sanitize inputs and prefer parameterized queries or NewField methods for field construction.
// Use this low-level function only when absolutely necessary, and ensure any embedded values are properly escaped.
func NewUnsafeFieldRaw(rawSQL string, vars ...interface{}) Field {
	return Field{expr: expr{e: clause.Expr{SQL: rawSQL, Vars: vars}}}
}

// NewSerializer create new field2
func NewSerializer(table, column string, opts ...Option) Serializer {
	return Serializer{expr: expr{col: toColumn(table, column, opts...)}}
}

// NewAsterisk create new * field
func NewAsterisk(table string, opts ...Option) Asterisk {
	return Asterisk{asteriskExpr: asteriskExpr{expr{col: toColumn(table, "*", opts...)}}}
}

// ======================== integer =======================

// NewNumber build number type field
func NewNumber[T constraints.Integer | constraints.Float](table, column string, opts ...Option) Number[T] {
	return newNumber[T](expr{col: toColumn(table, column, opts...)})
}

// NewInt create new field for int
func NewInt(table, column string, opts ...Option) Int {
	return NewNumber[int](table, column, opts...)
}

// NewInt8 create new field for int8
func NewInt8(table, column string, opts ...Option) Int8 {
	return NewNumber[int8](table, column, opts...)
}

// NewInt16 create new field for int16
func NewInt16(table, column string, opts ...Option) Int16 {
	return NewNumber[int16](table, column, opts...)
}

// NewInt32 create new field for int32
func NewInt32(table, column string, opts ...Option) Int32 {
	return NewNumber[int32](table, column, opts...)
}

// NewInt64 create new field for int64
func NewInt64(table, column string, opts ...Option) Int64 {
	return NewNumber[int64](table, column, opts...)
}

// NewUint create new field for uint
func NewUint(table, column string, opts ...Option) Uint {
	return NewNumber[uint](table, column, opts...)
}

// NewUint8 create new field for uint8
func NewUint8(table, column string, opts ...Option) Uint8 {
	return NewNumber[uint8](table, column, opts...)
}

// NewUint16 create new field for uint16
func NewUint16(table, column string, opts ...Option) Uint16 {
	return NewNumber[uint16](table, column, opts...)
}

// NewUint32 create new field for uint32
func NewUint32(table, column string, opts ...Option) Uint32 {
	return NewNumber[uint32](table, column, opts...)
}

// NewUint64 create new field for uint64
func NewUint64(table, column string, opts ...Option) Uint64 {
	return NewNumber[uint64](table, column, opts...)
}

// ======================== float =======================

// NewFloat32 create new field for float32
func NewFloat32(table, column string, opts ...Option) Float32 {
	return NewNumber[float32](table, column, opts...)
}

// NewFloat64 create new field for float64
func NewFloat64(table, column string, opts ...Option) Float64 {
	return NewNumber[float64](table, column, opts...)
}

// ======================== string =======================

// NewString ...
func NewString(table, column string, opts ...Option) String {
	return newChars[string](expr{col: toColumn(table, column, opts...)})
}

// NewBytes ...
func NewBytes(table, column string, opts ...Option) Bytes {
	return newChars[[]byte](expr{col: toColumn(table, column, opts...)})
}

// ======================== bool =======================

// NewBool ...
func NewBool(table, column string, opts ...Option) Bool {
	return Bool{genericsField[bool]{expr{col: toColumn(table, column, opts...)}}}
}

// ======================== time =======================

// NewTime ...
func NewTime(table, column string, opts ...Option) Time {
	return Time{genericsField[time.Time]{expr{col: toColumn(table, column, opts...)}}}
}

func toColumn(table, column string, opts ...Option) clause.Column {
	col := clause.Column{Table: table, Name: column}
	for _, opt := range opts {
		col = opt(col)
	}
	return banColumnRaw(col)
}

// ======================== boolean operate ========================

// Or return or condition
func Or(exprs ...Expr) Expr {
	return &expr{e: clause.Or(toExpression(exprs...)...)}
}

// And return and condition
func And(exprs ...Expr) Expr {
	return &expr{e: clause.And(toExpression(exprs...)...)}
}

// Not return not condition
func Not(exprs ...Expr) Expr {
	return &expr{e: clause.Not(toExpression(exprs...)...)}
}

func toExpression(conds ...Expr) []clause.Expression {
	exprs := make([]clause.Expression, len(conds))
	for i, cond := range conds {
		exprs[i] = cond.expression()
	}
	return exprs
}

// ======================== subquery method ========================

// ContainsSubQuery return contains subquery
// when len(columns) == 1, equal to columns[0] IN (subquery)
// when len(columns) > 1, equal to (columns[0], columns[1], ...) IN (subquery)
func ContainsSubQuery(columns []Expr, subQuery *gorm.DB) Expr {
	switch len(columns) {
	case 0:
		return expr{e: clause.Expr{}}
	case 1:
		return expr{e: clause.Expr{
			SQL:  "? IN (?)",
			Vars: []interface{}{columns[0].RawExpr(), subQuery},
		}}
	default: // len(columns) > 0
		placeholders := make([]string, len(columns))
		cols := make([]interface{}, len(columns))
		for i, c := range columns {
			placeholders[i], cols[i] = "?", c.RawExpr()
		}
		return expr{e: clause.Expr{
			SQL:  fmt.Sprintf("(%s) IN (?)", strings.Join(placeholders, ",")),
			Vars: append(cols, subQuery),
		}}
	}
}

// AssignSubQuery assign with subquery
func AssignSubQuery(columns []Expr, subQuery *gorm.DB) AssignExpr {
	cols := make([]string, len(columns))
	for i, c := range columns {
		cols[i] = string(c.BuildColumn(subQuery.Statement))
	}

	name := cols[0]
	if len(cols) > 1 {
		name = "(" + strings.Join(cols, ",") + ")"
	}

	return expr{e: clause.Set{{
		Column: clause.Column{Name: name, Raw: true},
		Value:  gorm.Expr("(?)", subQuery),
	}}}
}

// CompareOperator compare operator
type CompareOperator string

const (
	// EqOp =
	EqOp CompareOperator = " = "
	// NeqOp <>
	NeqOp CompareOperator = " <> "
	// GtOp >
	GtOp CompareOperator = " > "
	// GteOp >=
	GteOp CompareOperator = " >= "
	// LtOp <
	LtOp CompareOperator = " < "
	// LteOp <=
	LteOp CompareOperator = " <= "
	// ExistsOp EXISTS
	ExistsOp CompareOperator = "EXISTS "
)

// CompareSubQuery compare with sub query
func CompareSubQuery(op CompareOperator, column Expr, subQuery *gorm.DB) Expr {
	if op == ExistsOp {
		return expr{e: clause.Expr{
			SQL:  fmt.Sprint(op, "(?)"),
			Vars: []interface{}{subQuery},
		}}
	}
	return expr{e: clause.Expr{
		SQL:  fmt.Sprint("?", op, "(?)"),
		Vars: []interface{}{column.RawExpr(), subQuery},
	}}
}

// Value ...
type Value interface {
	expr() clause.Expr

	// implement Condition
	BeCond() interface{}
	CondError() error
}

type val clause.Expr

func (v val) expr() clause.Expr   { return clause.Expr(v) }
func (v val) BeCond() interface{} { return v }
func (val) CondError() error      { return nil }

// Values convert value to expression which implement Value
func Values(value interface{}) Value {
	return val(clause.Expr{
		SQL:                "?",
		Vars:               []interface{}{value},
		WithoutParentheses: true,
	})
}

// ContainsValue return expression which compare with value
func ContainsValue(columns []Expr, value Value) Expr {
	switch len(columns) {
	case 0:
		return expr{e: clause.Expr{}}
	case 1:
		return expr{e: clause.Expr{
			SQL:  "? IN (?)",
			Vars: []interface{}{columns[0].RawExpr(), value.expr()},
		}}
	default: // len(columns) > 0
		vars := make([]string, len(columns))
		queryCols := make([]interface{}, len(columns))
		for i, c := range columns {
			vars[i], queryCols[i] = "?", c.RawExpr()
		}
		return expr{e: clause.Expr{
			SQL:  fmt.Sprintf("(%s) IN (?)", strings.Join(vars, ", ")),
			Vars: append(queryCols, value.expr()),
		}}
	}
}

// EmptyExpr return a empty expression
func EmptyExpr() Expr { return expr{e: clause.Expr{}} }

// AssociationFields all association
var AssociationFields Expr = NewString("", clause.Associations).appendBuildOpts(WithoutQuote)

// Associations ...
var Associations RelationField = NewRelation(clause.Associations, "")

// NewRelation return a new Relation for association
func NewRelation(fieldName string, fieldType string, relations ...Relation) *Relation {
	return &Relation{
		fieldName:      fieldName,
		fieldPath:      fieldName,
		fieldType:      fieldType,
		childRelations: wrapPath(fieldName, relations),
	}
}

// NewRelationWithType return a Relation with specified field type
func NewRelationWithType(relationship RelationshipType, fieldName string, fieldType string, relations ...Relation) *Relation {
	return &Relation{
		relationship:   relationship,
		fieldName:      fieldName,
		fieldType:      fieldType,
		fieldPath:      fieldName,
		childRelations: wrapPath(fieldName, relations),
	}
}

// NewRelationWithModel return a Relation with specified model struct
func NewRelationWithModel(relationship RelationshipType, fieldName string, fieldType string, fieldModel interface{}, relations ...Relation) *Relation {
	return &Relation{
		relationship: relationship,
		fieldName:    fieldName,
		fieldType:    fieldType,
		fieldPath:    fieldName,
		fieldModel:   fieldModel,
	}
}
