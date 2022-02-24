package field

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Star a symbol of "*"
var Star = NewString("", "*")
var ALL = Star

type FieldOption func(clause.Column) clause.Column

var (
	banColumnRaw FieldOption = func(col clause.Column) clause.Column {
		col.Raw = false
		return col
	}
)

// ======================== generic field =======================

func NewField(table, column string, opts ...FieldOption) Field {
	return Field{expr: expr{col: toColumn(table, column, opts...)}}
}

// ======================== integer =======================

func NewInt(table, column string, opts ...FieldOption) Int {
	return Int{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewInt8(table, column string, opts ...FieldOption) Int8 {
	return Int8{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewInt16(table, column string, opts ...FieldOption) Int16 {
	return Int16{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewInt32(table, column string, opts ...FieldOption) Int32 {
	return Int32{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewInt64(table, column string, opts ...FieldOption) Int64 {
	return Int64{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewUint(table, column string, opts ...FieldOption) Uint {
	return Uint{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewUint8(table, column string, opts ...FieldOption) Uint8 {
	return Uint8{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewUint16(table, column string, opts ...FieldOption) Uint16 {
	return Uint16{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewUint32(table, column string, opts ...FieldOption) Uint32 {
	return Uint32{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewUint64(table, column string, opts ...FieldOption) Uint64 {
	return Uint64{expr: expr{col: toColumn(table, column, opts...)}}
}

// ======================== float =======================

func NewFloat32(table, column string, opts ...FieldOption) Float32 {
	return Float32{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewFloat64(table, column string, opts ...FieldOption) Float64 {
	return Float64{expr: expr{col: toColumn(table, column, opts...)}}
}

// ======================== string =======================

func NewString(table, column string, opts ...FieldOption) String {
	return String{expr: expr{col: toColumn(table, column, opts...)}}
}

func NewBytes(table, column string, opts ...FieldOption) Bytes {
	return Bytes{expr: expr{col: toColumn(table, column, opts...)}}
}

// ======================== bool =======================

func NewBool(table, column string, opts ...FieldOption) Bool {
	return Bool{expr: expr{col: toColumn(table, column, opts...)}}
}

// ======================== time =======================

func NewTime(table, column string, opts ...FieldOption) Time {
	return Time{expr: expr{col: toColumn(table, column, opts...)}}
}

func toColumn(table, column string, opts ...FieldOption) clause.Column {
	col := clause.Column{Table: table, Name: column}
	for _, opt := range opts {
		col = opt(col)
	}
	return banColumnRaw(col)
}

// ======================== boolean operate ========================
func Or(exprs ...Expr) Expr {
	return &expr{e: clause.Or(toExpression(exprs...)...)}
}

func And(exprs ...Expr) Expr {
	return &expr{e: clause.And(toExpression(exprs...)...)}
}

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

type CompareOperate string

const (
	EqOp     CompareOperate = " = "
	NeqOp    CompareOperate = " <> "
	GtOp     CompareOperate = " > "
	GteOp    CompareOperate = " >= "
	LtOp     CompareOperate = " < "
	LteOp    CompareOperate = " <= "
	ExistsOp CompareOperate = "EXISTS "
)

func CompareSubQuery(op CompareOperate, column Expr, subQuery *gorm.DB) Expr {
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

func Values(value interface{}) Value {
	return val(clause.Expr{
		SQL:                "?",
		Vars:               []interface{}{value},
		WithoutParentheses: true,
	})
}

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

func EmptyExpr() Expr { return expr{e: clause.Expr{}} }

var AssociationFields Expr = NewString("", clause.Associations)
var Associations RelationField = NewRelation(clause.Associations, "")

func NewRelation(fieldName string, fieldType string, relations ...Relation) *Relation {
	return &Relation{
		fieldName:      fieldName,
		fieldPath:      fieldName,
		fieldType:      fieldType,
		childRelations: wrapPath(fieldName, relations),
	}
}

func NewRelationWithType(relationship RelationshipType, fieldName string, fieldType string, relations ...Relation) *Relation {
	return &Relation{
		relationship:   relationship,
		fieldName:      fieldName,
		fieldType:      fieldType,
		fieldPath:      fieldName,
		childRelations: wrapPath(fieldName, relations),
	}
}

func NewRelationWithModel(relationship RelationshipType, fieldName string, fieldType string, fieldModel interface{}, relations ...Relation) *Relation {
	return &Relation{
		relationship: relationship,
		fieldName:    fieldName,
		fieldType:    fieldType,
		fieldPath:    fieldName,
		fieldModel:   fieldModel,
	}
}
