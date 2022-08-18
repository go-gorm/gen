package field

import (
	"gorm.io/gorm/clause"
)

// Asterisk a type of xxx.*
type Asterisk struct{ expr }

// Count count
func (a Asterisk) Count() Int {
	return Int{a.setE(clause.Expr{
		SQL:  "COUNT(?.*)",
		Vars: []interface{}{clause.Column{Table: a.col.Table}},
	})}
}

func (a Asterisk) Distinct() Int {
	return Int{a.setE(clause.Expr{
		SQL:  "DISTINCT ?",
		Vars: []interface{}{clause.Column{Table: a.col.Table}},
	})}
}
