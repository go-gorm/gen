package field

import "gorm.io/gorm/clause"

var (
	UnixTimestamp = Uint64{expr{e: clause.Expr{SQL: "UNIX_TIMESTAMP()"}}}
)
