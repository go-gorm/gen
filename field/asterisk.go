package field

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Asterisk a type of xxx.*
type Asterisk struct{ asteriskExpr }

// Count count
func (a Asterisk) Count() Asterisk {
	var expr *clause.Expr
	switch {
	case a.e != nil:
		expr = &clause.Expr{
			SQL:  "COUNT(?)",
			Vars: []interface{}{a.e},
		}
	case a.col.Table == "":
		expr = &clause.Expr{SQL: "COUNT(*)"}
	default:
		expr = &clause.Expr{
			SQL:  "COUNT(?.*)",
			Vars: []interface{}{clause.Table{Name: a.col.Table}},
		}
	}
	return Asterisk{asteriskExpr{expr: a.setE(expr)}}
}

// Distinct distinct
func (a Asterisk) Distinct() Asterisk {
	var expr *clause.Expr
	if a.col.Table == "" {
		expr = &clause.Expr{SQL: "DISTINCT *"}
	} else {
		expr = &clause.Expr{
			SQL:  "DISTINCT ?.*",
			Vars: []interface{}{clause.Table{Name: a.col.Table}},
		}
	}
	return Asterisk{asteriskExpr{expr: a.setE(expr)}}
}

type asteriskExpr struct{ expr }

func (e asteriskExpr) BuildWithArgs(*gorm.Statement) (query sql, args []interface{}) {
	// if e.expr has no expression it must be directly calling for "*" or "xxx.*"
	if e.e != nil {
		return "?", []interface{}{e.e}
	}
	if e.col.Table == "" {
		return "*", nil
	}
	return "?.*", []interface{}{clause.Table{Name: e.col.Table}}
}
