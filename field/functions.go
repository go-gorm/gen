package field

import "gorm.io/gorm/clause"

// Func sql functions
var Func = &function{
	UnixTimestamp: Uint64{expr{e: clause.Expr{SQL: "UNIX_TIMESTAMP()"}}},
}

type function struct {
	UnixTimestamp Uint64
}
