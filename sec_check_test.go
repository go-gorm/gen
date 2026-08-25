package gen

import (
	"errors"
	"testing"

	"gorm.io/gorm/clause"
	"gorm.io/hints"
	"gorm.io/plugin/dbresolver"
)

type policyClause string

func (c policyClause) Name() string             { return string(c) }
func (policyClause) Build(clause.Builder)       {}
func (policyClause) MergeClause(*clause.Clause) {}

func TestCheckClauseDefaultPolicy(t *testing.T) {
	for _, name := range []string{
		"VALUES", "SELECT", "FROM", "WHERE", "GROUP BY",
		"ORDER BY", "LIMIT", "UPDATE", "SET", "DELETE",
	} {
		t.Run("reject_"+name, func(t *testing.T) {
			if err := CheckClause(policyClause(name)); err == nil {
				t.Fatalf("expected %s to be rejected", name)
			}
		})
	}

	allowed := []clause.Expression{
		policyClause("FOR"),
		clause.Insert{},
		clause.OnConflict{},
		clause.Locking{Strength: clause.LockingStrengthUpdate},
		hints.New("MAX_EXECUTION_TIME(1000)"),
		hints.IndexHint{},
		dbresolver.Read,
	}
	for _, expression := range allowed {
		if err := CheckClause(expression); err != nil {
			t.Fatalf("expected %T to be allowed: %v", expression, err)
		}
	}

	if err := CheckClause(clause.Expr{SQL: "unsafe"}); err == nil {
		t.Fatal("expected an unknown expression to be rejected")
	}
}

func TestCheckCondsStopsAtFirstRejectedClause(t *testing.T) {
	conds := []clause.Expression{
		policyClause("FOR"),
		policyClause("WHERE"),
		policyClause("SELECT"),
	}
	if err := checkConds(conds); err == nil {
		t.Fatal("expected a rejected condition")
	}
	if err := checkConds(nil); err != nil {
		t.Fatalf("empty conditions should be accepted: %v", err)
	}
}

func TestCheckCondsWithCheckerContract(t *testing.T) {
	t.Run("handled", func(t *testing.T) {
		calls := 0
		err := checkCondsWithChecker([]clause.Expression{policyClause("WHERE")}, func(clause.Expression) error {
			calls++
			return nil
		})
		if err != nil {
			t.Fatalf("custom checker should handle the clause: %v", err)
		}
		if calls != 1 {
			t.Fatalf("checker calls: got %d want 1", calls)
		}
	})

	t.Run("fallback", func(t *testing.T) {
		err := checkCondsWithChecker([]clause.Expression{policyClause("WHERE")}, func(clause.Expression) error {
			return ErrClauseNotHandled
		})
		if err == nil {
			t.Fatal("default checker should reject WHERE after fallback")
		}
	})

	t.Run("custom_error", func(t *testing.T) {
		want := errors.New("custom policy error")
		err := checkCondsWithChecker([]clause.Expression{policyClause("FOR")}, func(clause.Expression) error {
			return want
		})
		if !errors.Is(err, want) {
			t.Fatalf("got %v want %v", err, want)
		}
	})

	if err := checkCondsWithChecker(nil, nil); err != nil {
		t.Fatalf("empty conditions should be accepted: %v", err)
	}
}

func TestCheckInsert(t *testing.T) {
	tests := []struct {
		name    string
		insert  clause.Insert
		wantErr bool
	}{
		{name: "empty"},
		{name: "ignore", insert: clause.Insert{Modifier: "IGNORE"}},
		{name: "priority and ignore", insert: clause.Insert{Modifier: " low_priority   ignore "}},
		{name: "delayed and ignore", insert: clause.Insert{Modifier: "DELAYED IGNORE"}},
		{name: "high priority and ignore", insert: clause.Insert{Modifier: "HIGH_PRIORITY IGNORE"}},
		{name: "raw table", insert: clause.Insert{Table: clause.Table{Name: "users", Raw: true}}, wantErr: true},
		{name: "priority only currently rejected", insert: clause.Insert{Modifier: "LOW_PRIORITY"}, wantErr: true},
		{name: "invalid priority", insert: clause.Insert{Modifier: "FAST IGNORE"}, wantErr: true},
		{name: "invalid modifier", insert: clause.Insert{Modifier: "LOW_PRIORITY REPLACE"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkInsert(tt.insert)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckLocking(t *testing.T) {
	tests := []struct {
		name    string
		locking clause.Locking
		wantErr bool
	}{
		{name: "update", locking: clause.Locking{Strength: " update "}},
		{name: "share nowait", locking: clause.Locking{Strength: "SHARE", Options: "nowait"}},
		{name: "skip locked", locking: clause.Locking{Strength: "UPDATE", Options: " skip locked "}},
		{name: "empty strength", locking: clause.Locking{}, wantErr: true},
		{name: "invalid strength", locking: clause.Locking{Strength: "EXCLUSIVE"}, wantErr: true},
		{name: "raw table", locking: clause.Locking{Strength: "UPDATE", Table: clause.Table{Raw: true}}, wantErr: true},
		{name: "invalid option", locking: clause.Locking{Strength: "UPDATE", Options: "WAIT"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkLocking(tt.locking)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkLocking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckOnConflict(t *testing.T) {
	tests := []struct {
		name       string
		assignment interface{}
		wantErr    bool
	}{
		{name: "plain value", assignment: "alice"},
		{name: "expression", assignment: clause.Expr{SQL: "excluded.name"}, wantErr: true},
		{name: "expression pointer", assignment: &clause.Expr{SQL: "excluded.name"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkOnConflict(clause.OnConflict{DoUpdates: clause.Set{{
				Column: clause.Column{Name: "name"},
				Value:  tt.assignment,
			}}})
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkOnConflict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIn(t *testing.T) {
	if !in("b", "a", "b", "c") {
		t.Fatal("expected value to be found")
	}
	if in("x", "a", "b", "c") {
		t.Fatal("unexpected value found")
	}
}
