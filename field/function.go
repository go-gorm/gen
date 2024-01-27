package field

import (
	"strings"

	"gorm.io/gorm/clause"
)

// Func sql functions
var Func = new(function)

type function struct{}

// UnixTimestamp same as UNIX_TIMESTAMP([date])
func (f *function) UnixTimestamp(date ...string) Uint64 {
	if len(date) > 0 {
		return Uint64{GenericsInt: GenericsInt[uint64]{GenericsField: GenericsField[uint64]{expr{e: clause.Expr{SQL: "UNIX_TIMESTAMP(?)", Vars: []interface{}{date[0]}}}}}}
	}
	return Uint64{GenericsInt: GenericsInt[uint64]{GenericsField: GenericsField[uint64]{expr{e: clause.Expr{SQL: "UNIX_TIMESTAMP()"}}}}}
}

// FromUnixTime FROM_UNIXTIME(unix_timestamp[,format])
func (f *function) FromUnixTime(date uint64, format string) String {
	if strings.TrimSpace(format) != "" {
		return String{GenericsString: GenericsString[string]{GenericsField: GenericsField[string]{expr{e: clause.Expr{SQL: "FROM_UNIXTIME(?, ?)", Vars: []interface{}{date, format}}}}}}
	}
	return String{GenericsString: GenericsString[string]{GenericsField: GenericsField[string]{expr{e: clause.Expr{SQL: "FROM_UNIXTIME(?)", Vars: []interface{}{date}}}}}}
}

func (f *function) Rand() String {
	return String{GenericsString: GenericsString[string]{GenericsField: GenericsField[string]{expr{e: clause.Expr{SQL: "RAND()"}}}}}

}
