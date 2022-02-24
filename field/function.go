package field

import (
	"gorm.io/gorm/clause"
)

// Func sql functions
var Func = new(function)

type function struct{}

func (f *function) UnixTimestamp(date ...string) Uint64 {
	if len(date) > 0 {
		return Uint64{expr{e: clause.Expr{SQL: "UNIX_TIMESTAMP(?)", Vars: []interface{}{date[0]}}}}
	}
	return Uint64{expr{e: clause.Expr{SQL: "UNIX_TIMESTAMP()"}}}
}
