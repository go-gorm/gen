package gen

import (
	"fmt"

	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

func checkConds(conds []clause.Expression) error {
	for _, cond := range conds {
		if err := checkClause(cond); err != nil {
			return err
		}
	}
	return nil
}

var banClauses = map[string]bool{
	"INSERT": true,
	"VALUES": true,
	// "ON CONFLICT": true,
	"SELECT":   true,
	"FROM":     true,
	"WHERE":    true,
	"GROUP BY": true,
	"ORDER BY": true,
	"LIMIT":    true,
	"FOR":      true,
	"UPDATE":   true,
	"SET":      true,
	"DELETE":   true,
}

func checkClause(cond clause.Expression) error {
	switch cond := cond.(type) {
	case hints.Hints, hints.IndexHint:
		return nil
	case clause.OnConflict:
		return checkOnConflict(cond)
	case clause.Interface:
		if banClauses[cond.Name()] {
			return fmt.Errorf("clause %s is banned", cond.Name())
		}
		return nil
	}
	return fmt.Errorf("unknown clause %v", cond)
}

func checkOnConflict(cond clause.OnConflict) error {
	for _, item := range cond.DoUpdates {
		switch item.Value.(type) {
		case clause.Expr, *clause.Expr:
			return fmt.Errorf("OnConflict clause assignment with gorm.Expr is banned for security reasons for now")
		}
	}
	return nil
}
