package gen

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm/clause"
	"gorm.io/hints"
	"gorm.io/plugin/dbresolver"
)

func checkConds(conds []clause.Expression) error {
	for _, cond := range conds {
		if err := CheckClause(cond); err != nil {
			return err
		}
	}
	return nil
}

var banClauses = map[string]bool{
	// "INSERT": true,
	"VALUES": true,
	// "ON CONFLICT": true,
	"SELECT":   true,
	"FROM":     true,
	"WHERE":    true,
	"GROUP BY": true,
	"ORDER BY": true,
	"LIMIT":    true,
	// "FOR":      true,
	"UPDATE": true,
	"SET":    true,
	"DELETE": true,
}

// CheckClause check security of Expression
func CheckClause(cond clause.Expression) error {
	switch cond := cond.(type) {
	case hints.Hints, hints.IndexHint, dbresolver.Operation:
		return nil
	case clause.OnConflict:
		return checkOnConflict(cond)
	case clause.Locking:
		return checkLocking(cond)
	case clause.Insert:
		return checkInsert(cond)
	case clause.Interface:
		if banClauses[cond.Name()] {
			return fmt.Errorf("clause %s is banned", cond.Name())
		}
		return nil
	}
	return fmt.Errorf("unknown clause %v", cond)
}

func checkOnConflict(c clause.OnConflict) error {
	for _, item := range c.DoUpdates {
		switch item.Value.(type) {
		case clause.Expr, *clause.Expr:
			return errors.New("OnConflict clause assignment with gorm.Expr is banned for security reasons for now")
		}
	}
	return nil
}

func checkLocking(c clause.Locking) error {
	if strength := strings.ToUpper(strings.TrimSpace(c.Strength)); strength != "UPDATE" && strength != "SHARE" {
		return errors.New("Locking clause's Strength only allow assignments of UPDATE/SHARE")
	}
	if c.Table.Raw {
		return errors.New("Locking clause's Table cannot be set Raw==true")
	}
	if options := strings.ToUpper(strings.TrimSpace(c.Options)); options != "" && options != "NOWAIT" && options != "SKIP LOCKED" {
		return errors.New("Locking clause's Options only allow assignments of NOWAIT/SKIP LOCKED for now")
	}
	return nil
}

// checkInsert check if clause.Insert is safe
// https://dev.mysql.com/doc/refman/8.0/en/sql-statements.html#insert
func checkInsert(c clause.Insert) error {
	if c.Table.Raw == true {
		return errors.New("Table Raw cannot be true")
	}

	if c.Modifier == "" {
		return nil
	}

	var priority, ignore string
	if modifiers := strings.SplitN(strings.ToUpper(strings.TrimSpace(c.Modifier)), " ", 2); len(modifiers) == 2 {
		priority, ignore = strings.TrimSpace(modifiers[0]), strings.TrimSpace(modifiers[1])
	} else {
		ignore = strings.TrimSpace(modifiers[0])
	}
	if priority != "" && !in(priority, "LOW_PRIORITY", "DELAYED", "HIGH_PRIORITY") {
		return errors.New("invalid priority value")
	}
	if ignore != "" && ignore != "IGNORE" {
		return errors.New("invalid modifiers value, should be IGNORE")
	}
	return nil
}

func in(s string, v ...string) bool {
	for _, vv := range v {
		if vv == s {
			return true
		}
	}
	return false
}
