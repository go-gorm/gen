package gen

import (
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils/tests"
	"gorm.io/plugin/dbresolver"
)

func TestDOClausesWithClauseChecker(t *testing.T) {
	base, err := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	dry1 := base.Session(&gorm.Session{DryRun: true, NewDB: true})
	var d1 DO
	d1.UseDB(dry1)
	dao1 := d1.Clauses(dbresolver.Use("db_1")).(*DO)
	if dao1.db.Error == nil {
		t.Fatalf("expected error, got nil")
	}

	dry2 := base.Session(&gorm.Session{DryRun: true, NewDB: true})
	var d2 DO
	d2.UseDB(dry2, WithClauseChecker(func(clause.Expression) error { return ErrClauseNotHandled }))
	dao2 := d2.Clauses(dbresolver.Use("db_1")).(*DO)
	if dao2.db.Error == nil {
		t.Fatalf("expected error, got nil")
	}

	dry3 := base.Session(&gorm.Session{DryRun: true, NewDB: true})
	var d3 DO
	d3.UseDB(dry3, WithClauseChecker(func(clause.Expression) error { return nil }))
	dao3 := d3.Clauses(dbresolver.Use("db_1")).(*DO)
	if dao3.db.Error != nil {
		t.Fatalf("unexpected error: %v", dao3.db.Error)
	}
}
