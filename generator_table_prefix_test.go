package gen

import (
	"fmt"
	"strings"
	"testing"

	"gorm.io/gen/internal/generate"
)

// newPrefixTestGenerator builds a Generator suitable for exercising
// registerModel directly (white-box); no database introspection is involved.
func newPrefixTestGenerator() *Generator {
	g := &Generator{
		Config:            Config{},
		Data:              make(map[string]*genInfo),
		models:            make(map[string]*generate.QueryStructMeta),
		modelTableSources: make(map[string]string),
		logger:            testLogger{},
	}
	g.db = db
	return g
}

type testLogger struct{}

func (testLogger) Println(v ...any) {}

func prefixTestMeta(tableName, modelName string) *generate.QueryStructMeta {
	return &generate.QueryStructMeta{
		TableName:       tableName,
		ModelStructName: modelName,
	}
}

// TestRegisterModelTablePrefixCollision covers the audit finding M1: with
// TablePrefix "t_" the source tables "users" (prefix applied) and "t_users"
// (prefix skipped) used to normalize to the same generated table name, so the
// generated DO silently queried the wrong physical table at runtime.
// registerModel must reject that collision instead of emitting broken code.
func TestRegisterModelTablePrefixCollision(t *testing.T) {
	t.Run("distinct sources with the same generated table name are rejected", func(t *testing.T) {
		g := newPrefixTestGenerator()
		g.registerModel("users", prefixTestMeta("t_users", "User"))

		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("registerModel() did not panic on table name collision")
			}
			msg := fmt.Sprint(r)
			for _, want := range []string{`"users"`, `"t_users"`, "t_users"} {
				if !strings.Contains(msg, want) {
					t.Fatalf("panic message %q does not mention %q", msg, want)
				}
			}
		}()
		g.registerModel("t_users", prefixTestMeta("t_users", "TUser"))
	})

	t.Run("object and table sources colliding are rejected", func(t *testing.T) {
		g := newPrefixTestGenerator()
		g.registerModel("t_users", prefixTestMeta("t_users", "UserView")) // object-sourced model

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("registerModel() did not panic on object/table collision")
			}
		}()
		g.registerModel("users", prefixTestMeta("t_users", "User")) // table-sourced, prefix applied
	})

	t.Run("the same source table may back multiple models", func(t *testing.T) {
		g := newPrefixTestGenerator()
		g.registerModel("users", prefixTestMeta("t_users", "User"))
		g.registerModel("users", prefixTestMeta("t_users", "ActiveUser")) // must not panic

		if len(g.models) != 2 {
			t.Fatalf("registered models = %d, want 2", len(g.models))
		}
		if got := g.modelTableSources["t_users"]; got != "users" {
			t.Fatalf("modelTableSources[t_users] = %q, want users", got)
		}
	})

	t.Run("distinct sources with distinct table names are accepted", func(t *testing.T) {
		g := newPrefixTestGenerator()
		g.registerModel("users", prefixTestMeta("t_users", "User"))
		g.registerModel("orders", prefixTestMeta("t_orders", "Order")) // must not panic

		if len(g.models) != 2 {
			t.Fatalf("registered models = %d, want 2", len(g.models))
		}
	})

	t.Run("nil meta is ignored", func(t *testing.T) {
		g := newPrefixTestGenerator()
		g.registerModel("users", nil) // must not panic
		if len(g.models) != 0 || len(g.modelTableSources) != 0 {
			t.Fatalf("unexpected state after nil meta: models=%d sources=%d", len(g.models), len(g.modelTableSources))
		}
	})
}
