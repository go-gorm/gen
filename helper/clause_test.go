package helper

import (
	"errors"
	"strings"
	"testing"

	"gorm.io/gorm/clause"
)

func TestIfClause(t *testing.T) {
	tests := []struct {
		name  string
		conds []Cond
		want  string
	}{
		{name: "nil", want: " "},
		{name: "true", conds: []Cond{{Cond: true, Result: " name = ? "}}, want: " name = ?"},
		{name: "false", conds: []Cond{{Cond: false, Result: "name = ?"}}, want: " "},
		{
			name: "mixed",
			conds: []Cond{
				{Cond: true, Result: " name = ? "},
				{Cond: false, Result: "ignored"},
				{Cond: true, Result: " age > ? "},
			},
			want: " name = ?  age > ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IfClause(tt.conds); got != tt.want {
				t.Fatalf("IfClause() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWhereClause(t *testing.T) {
	tests := []struct {
		name  string
		conds []string
		want  string
	}{
		{name: "nil"},
		{name: "empty values", conds: []string{"", "  "}},
		{name: "single", conds: []string{"name = ?"}, want: " WHERE name = ?"},
		{name: "mixed prefixes", conds: []string{" AND name = ? ", "or age > ?"}, want: " WHERE name = ? or age > ?"},
		{name: "xor", conds: []string{"xor enabled = ?"}, want: " WHERE enabled = ?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WhereClause(tt.conds); got != tt.want {
				t.Fatalf("WhereClause() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetClause(t *testing.T) {
	tests := []struct {
		name  string
		conds []string
		want  string
	}{
		{name: "nil"},
		{name: "empty values", conds: []string{"", ","}},
		{name: "single", conds: []string{"name = ?,"}, want: " SET name = ?"},
		{name: "multiple", conds: []string{" name = ?, ", ", age = ?"}, want: " SET name = ?,age = ?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetClause(tt.conds); got != tt.want {
				t.Fatalf("SetClause() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTrimHelpers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "and", input: " AND name = ? AND ", want: "name = ? "},
		{name: "or", input: " or name = ? OR ", want: "name = ? "},
		{name: "xor", input: " XoR enabled = ? xor ", want: "enabled = ? "},
		{name: "comma", input: ",name = ? ,", want: "name = ? "},
		{name: "ordinary word", input: "order_no = ?", want: "order_no = ?"},
		{name: "empty", input: "   ", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimAll(tt.input); got != tt.want {
				t.Fatalf("trimAll(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestJoinBuilders(t *testing.T) {
	tests := []struct {
		name string
		join func(*strings.Builder, strings.Builder)
		in   string
		want string
	}{
		{name: "where empty", join: JoinWhereBuilder, want: ""},
		{name: "where", join: JoinWhereBuilder, in: "and name = ? or ", want: "WHERE name = ?  "},
		{name: "set empty", join: JoinSetBuilder, want: ""},
		{name: "set", join: JoinSetBuilder, in: "name = ?,", want: "SET name = ? "},
		{name: "trim empty", join: JoinTrimAllBuilder, want: " "},
		{name: "trim", join: JoinTrimAllBuilder, in: "name = ? or", want: "name = ?  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var src strings.Builder
			var value strings.Builder
			value.WriteString(tt.in)
			tt.join(&src, value)
			if got := src.String(); got != tt.want {
				t.Fatalf("builder output = %q, want %q", got, tt.want)
			}
		})
	}
}

type testClauseBuilder struct {
	strings.Builder
	vars []interface{}
	err  error
}

func (b *testClauseBuilder) WriteQuoted(value interface{}) {
	switch value := value.(type) {
	case string:
		b.WriteByte('`')
		b.WriteString(value)
		b.WriteByte('`')
	case clause.Column:
		b.WriteByte('`')
		b.WriteString(value.Name)
		b.WriteByte('`')
	}
}

func (b *testClauseBuilder) AddVar(writer clause.Writer, values ...interface{}) {
	for idx, value := range values {
		if idx > 0 {
			_ = writer.WriteByte(',')
		}
		_, _ = writer.WriteString("?")
		b.vars = append(b.vars, value)
	}
}

func (b *testClauseBuilder) AddError(err error) error {
	b.err = errors.Join(b.err, err)
	return b.err
}

func TestJoinTblExprBuild(t *testing.T) {
	NewJoinTblExpr(clause.Join{}, nil).Build(nil)

	t.Run("on", func(t *testing.T) {
		builder := new(testClauseBuilder)
		join := NewJoinTblExpr(clause.Join{
			Type: clause.LeftJoin,
			ON: clause.Where{Exprs: []clause.Expression{
				clause.Expr{SQL: "users.id = orders.user_id"},
			}},
		}, clause.Expr{SQL: "users"})
		join.Build(builder)
		if got, want := builder.String(), "LEFT JOIN users ON users.id = orders.user_id"; got != want {
			t.Fatalf("Build() = %q, want %q", got, want)
		}
	})

	t.Run("using", func(t *testing.T) {
		builder := new(testClauseBuilder)
		join := NewJoinTblExpr(clause.Join{
			Type:  clause.InnerJoin,
			Using: []string{"id", "tenant_id"},
		}, clause.Expr{SQL: "users"})
		join.Build(builder)
		if got, want := builder.String(), "INNER JOIN users USING (`id`,`tenant_id`)"; got != want {
			t.Fatalf("Build() = %q, want %q", got, want)
		}
	})
}
