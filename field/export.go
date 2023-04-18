package field

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// Star a symbol of "*"
	Star = NewAsterisk("")
	// ALL same with Star
	ALL = Star
)

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
	return Field{GenericsField: GenericsField[driver.Valuer]{expr{col: toColumn(table, column, opts...)}}}
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

// NewInt create new Int
func NewInt(table, column string, opts ...Option) Int {
	return Int{GenericsInt[int]{GenericsField: GenericsField[int]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewInt8 create new Int8
func NewInt8(table, column string, opts ...Option) Int8 {
	return Int8{GenericsInt[int8]{GenericsField: GenericsField[int8]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewInt16 ...
func NewInt16(table, column string, opts ...Option) Int16 {
	return Int16{GenericsInt[int16]{GenericsField: GenericsField[int16]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewInt32 ...
func NewInt32(table, column string, opts ...Option) Int32 {
	return Int32{GenericsInt[int32]{GenericsField: GenericsField[int32]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewInt64 ...
func NewInt64(table, column string, opts ...Option) Int64 {
	return Int64{GenericsInt[int64]{GenericsField: GenericsField[int64]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewUint ...
func NewUint(table, column string, opts ...Option) Uint {
	return Uint{GenericsInt[uint]{GenericsField: GenericsField[uint]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewUint8 ...
func NewUint8(table, column string, opts ...Option) Uint8 {
	return Uint8{GenericsInt[uint8]{GenericsField: GenericsField[uint8]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewUint16 ...
func NewUint16(table, column string, opts ...Option) Uint16 {
	return Uint16{GenericsInt[uint16]{GenericsField: GenericsField[uint16]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewUint32 ...
func NewUint32(table, column string, opts ...Option) Uint32 {
	return Uint32{GenericsInt[uint32]{GenericsField: GenericsField[uint32]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewUint64 ...
func NewUint64(table, column string, opts ...Option) Uint64 {
	return Uint64{GenericsInt[uint64]{GenericsField: GenericsField[uint64]{expr{col: toColumn(table, column, opts...)}}}}
}

// ======================== float =======================

// NewFloat32 ...
func NewFloat32(table, column string, opts ...Option) Float32 {
	return Float32{GenericsInt[float32]{GenericsField: GenericsField[float32]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewFloat64 ...
func NewFloat64(table, column string, opts ...Option) Float64 {
	return Float64{GenericsInt[float64]{GenericsField: GenericsField[float64]{expr{col: toColumn(table, column, opts...)}}}}
}

// ======================== string =======================

// NewString ...
func NewString(table, column string, opts ...Option) String {
	return String{GenericsString[string]{GenericsField: GenericsField[string]{expr{col: toColumn(table, column, opts...)}}}}
}

// NewBytes ...
func NewBytes(table, column string, opts ...Option) Bytes {
	return Bytes{GenericsString[[]byte]{GenericsField: GenericsField[[]byte]{expr{col: toColumn(table, column, opts...)}}}}
}

// ======================== bool =======================

// NewBool ...
func NewBool(table, column string, opts ...Option) Bool {
	return Bool{GenericsField[bool]{expr{col: toColumn(table, column, opts...)}}}
}

// ======================== time =======================

// NewTime ...
func NewTime(table, column string, opts ...Option) Time {
	return Time{GenericsField[time.Time]{expr{col: toColumn(table, column, opts...)}}}
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
