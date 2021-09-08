package field

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ Expr = new(Field)

// Expr a query expression about field
type Expr interface {
	As(alias string) Expr
	ColumnName() sql
	BuildColumn(*gorm.Statement, ...BuildOpt) sql
	Build(*gorm.Statement) sql
	RawExpr() expression

	ConditionTag()

	expression() clause.Expression
}

type expression interface{}

type sql string

func (e sql) String() string { return string(e) }

type expr struct {
	Col clause.Column

	e clause.Expression
}

func (expr) ConditionTag() {}

func (e expr) expression() clause.Expression {
	if e.e == nil {
		return clause.NamedExpr{SQL: "?", Vars: []interface{}{e.Col}}
	}
	return e.e
}

func (e expr) ColumnName() sql { return sql(e.Col.Name) }

type BuildOpt func(clause.Column) interface{}

var (
	// WithTable build column with table
	WithTable BuildOpt = func(col clause.Column) interface{} { return clause.Column{Table: col.Table, Name: col.Name} }

	// WithAll build column with table and alias
	WithAll BuildOpt = func(col clause.Column) interface{} { return col }
)

func (e expr) BuildColumn(stmt *gorm.Statement, opts ...BuildOpt) sql {
	var col interface{} = e.Col.Name
	for _, opt := range opts {
		col = opt(e.Col)
	}
	return sql(stmt.Quote(col))
}

func (e expr) Build(stmt *gorm.Statement) sql {
	if e.e == nil {
		return sql(e.BuildColumn(stmt, WithAll))
	}
	newStmt := &gorm.Statement{DB: stmt.DB, Table: stmt.Table, Schema: stmt.Schema}
	e.e.Build(newStmt)
	return sql(newStmt.SQL.String())
}

func (e expr) RawExpr() expression {
	if e.e == nil {
		return e.Col
	}
	return e.e
}

func (e expr) setE(expression clause.Expression) expr {
	e.e = expression
	return e
}

// ======================== basic function ========================
func (e expr) IsNull() Expr {
	return e.setE(clause.Expr{SQL: "? IS NULL", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) IsNotNull() Expr {
	return e.setE(clause.Expr{SQL: "? IS NOT NULL", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) Count() Int {
	return Int{e.setE(clause.Expr{SQL: "COUNT(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Distinct() Int {
	return Int{e.setE(clause.Expr{SQL: "DISTINCT ?", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Length() Int {
	return Int{e.setE(clause.Expr{SQL: "LENGTH(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Max() Float64 {
	return Float64{e.setE(clause.Expr{SQL: "MAX(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Min() Float64 {
	return Float64{e.setE(clause.Expr{SQL: "MIN(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Avg() Float64 {
	return Float64{e.setE(clause.Expr{SQL: "AVG(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) Sum() Float64 {
	return Float64{e.setE(clause.Expr{SQL: "SUM(?)", Vars: []interface{}{e.RawExpr()}})}
}

func (e expr) WithTable(table string) Expr {
	e.Col.Table = table
	return e
}

// ======================== comparison between columns ========================
func (e expr) EqCol(col Expr) Expr {
	return e.setE(clause.Expr{SQL: "? = ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) NeqCol(col Expr) Expr {
	return e.setE(clause.Expr{SQL: "? <> ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) GtCol(col Expr) Expr {
	return e.setE(clause.Expr{SQL: "? > ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) GteCol(col Expr) Expr {
	return e.setE(clause.Expr{SQL: "? >= ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) LtCol(col Expr) Expr {
	return e.setE(clause.Expr{SQL: "? < ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

func (e expr) LteCol(col Expr) Expr {
	return e.setE(clause.Expr{SQL: "? <= ?", Vars: []interface{}{e.RawExpr(), col.RawExpr()}})
}

// ======================== keyword ========================
func (e expr) As(alias string) Expr {
	if e.e != nil {
		return e.setE(clause.Expr{SQL: "? AS ?", Vars: []interface{}{e.e, clause.Column{Name: alias}}})
	}
	e.Col.Alias = alias
	return e
}

func (e expr) Desc() Expr {
	return e.setE(clause.Expr{SQL: "? DESC", Vars: []interface{}{e.RawExpr()}})
}

// ======================== general experssion ========================
func (e expr) between(values []interface{}) expr {
	return e.setE(clause.Expr{SQL: "? BETWEEN ? AND ?", Vars: append([]interface{}{e.RawExpr()}, values...)})
}

func (e expr) add(value interface{}) expr {
	switch v := value.(type) {
	case time.Duration:
		return e.setE(clause.Expr{SQL: "DATE_ADD(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.RawExpr(), v.Microseconds()}})
	default:
		return e.setE(clause.Expr{SQL: "?+?", Vars: []interface{}{e.RawExpr(), value}})
	}
}

func (e expr) sub(value interface{}) expr {
	switch v := value.(type) {
	case time.Duration:
		return e.setE(clause.Expr{SQL: "DATE_SUB(?, INTERVAL ? MICROSECOND)", Vars: []interface{}{e.RawExpr(), v.Microseconds()}})
	default:
		return e.setE(clause.Expr{SQL: "?-?", Vars: []interface{}{e.RawExpr(), value}})
	}
}

func (e expr) mul(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?*?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)*?", Vars: []interface{}{e.e, value}})
}

func (e expr) div(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?/?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)/?", Vars: []interface{}{e.e, value}})
}

func (e expr) mod(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?%?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)%?", Vars: []interface{}{e.e, value}})
}

func (e expr) floorDiv(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "? DIV ?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?) DIV ?", Vars: []interface{}{e.e, value}})
}

func (e expr) floor() expr {
	return e.setE(clause.Expr{SQL: "FLOOR(?)", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) rightShift(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?>>?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)>>?", Vars: []interface{}{e.e, value}})
}

func (e expr) leftShift(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?<<?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)<<?", Vars: []interface{}{e.e, value}})
}

func (e expr) bitXor(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?^?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)^?", Vars: []interface{}{e.e, value}})
}

func (e expr) bitAnd(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?&?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)&?", Vars: []interface{}{e.e, value}})
}

func (e expr) bitOr(value interface{}) expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "?|?", Vars: []interface{}{e.Col, value}})
	}
	return e.setE(clause.Expr{SQL: "(?)|?", Vars: []interface{}{e.e, value}})
}

func (e expr) bitFlip() expr {
	if e.isPure() {
		return e.setE(clause.Expr{SQL: "~?", Vars: []interface{}{e.Col}})
	}
	return e.setE(clause.Expr{SQL: "~(?)", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) regexp(value interface{}) expr {
	return e.setE(clause.Expr{SQL: "? REGEXP ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) not() expr {
	return e.setE(clause.Expr{SQL: "NOT ?", Vars: []interface{}{e.RawExpr()}})
}

func (e expr) is(value interface{}) expr {
	if value, ok := value.(bool); ok {
		return e.setE(clause.Expr{SQL: fmt.Sprintf("? IS %t", value), Vars: []interface{}{e.RawExpr()}})
	}
	return e.setE(clause.Expr{SQL: "? IS ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) and(value interface{}) expr {
	return e.setE(clause.Expr{SQL: "? AND ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) or(value interface{}) expr {
	return e.setE(clause.Expr{SQL: "? OR ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) xor(value interface{}) expr {
	return e.setE(clause.Expr{SQL: "? XOR ?", Vars: []interface{}{e.RawExpr(), value}})
}

func (e expr) isPure() bool {
	return e.e == nil
}
