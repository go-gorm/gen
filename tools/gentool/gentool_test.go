package main

import (
	"bytes"
	"flag"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCmdParamsReviseDefaultsAndTables(t *testing.T) {
	params := (&CmdParams{Tables: []string{" users ", "", " orders"}}).revise()
	if params.DB != string(dbMySQL) || params.OutPath != defaultQueryPath {
		t.Fatalf("defaults = db:%q out:%q", params.DB, params.OutPath)
	}
	if !reflect.DeepEqual(params.Tables, []string{"users", "orders"}) {
		t.Fatalf("tables = %v", params.Tables)
	}
	if (*CmdParams)(nil).revise() != nil {
		t.Fatal("nil revise should remain nil")
	}
}

func TestParseArgsUsesIndependentFlagSet(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.SetOutput(new(bytes.Buffer))
	params, err := parseArgs(fs, []string{
		"-db", "sqlite",
		"-dsn", "fixture.db",
		"-tables", " users, orders ",
		"-outPath", "generated/query",
		"-onlyModel",
		"-withDefaultQuery",
		"-withQueryInterface",
		"-useAny",
	})
	if err != nil {
		t.Fatalf("parseArgs(): %v", err)
	}
	params.revise()
	if params.DB != "sqlite" || params.DSN != "fixture.db" || params.OutPath != "generated/query" {
		t.Fatalf("parsed params = %#v", params)
	}
	if !reflect.DeepEqual(params.Tables, []string{"users", "orders"}) {
		t.Fatalf("tables = %v", params.Tables)
	}
	if !params.OnlyModel || !params.WithDefaultQuery || !params.WithQueryInterface || !params.UseAny {
		t.Fatalf("boolean flags = %#v", params)
	}

}

func TestParseArgsConfigFileIsExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gen.yml")
	content := []byte("version: 1\ndatabase:\n  db: sqlite\n  dsn: yaml.db\n  outPath: yaml/query\n  tables: [users]\n")
	if err := os.WriteFile(path, content, 0640); err != nil {
		t.Fatalf("write config: %v", err)
	}
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.SetOutput(new(bytes.Buffer))
	params, err := parseArgs(fs, []string{"-c", path, "-dsn", "ignored.db", "-outPath", "ignored/query"})
	if err != nil {
		t.Fatalf("parseArgs(): %v", err)
	}
	if params.DSN != "yaml.db" || params.OutPath != "yaml/query" {
		t.Fatalf("config file was not exclusive: %#v", params)
	}
}

func TestLoadConfigErrors(t *testing.T) {
	if _, err := loadConfig(filepath.Join(t.TempDir(), "missing.yml")); err == nil {
		t.Fatal("loadConfig() accepted a missing file")
	}
	for name, content := range map[string]string{
		"invalid.yml":     "database: [",
		"no-database.yml": "version: 1\n",
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			if err := os.WriteFile(path, []byte(content), 0640); err != nil {
				t.Fatalf("write config: %v", err)
			}
			if _, err := loadConfig(path); err == nil {
				t.Fatalf("loadConfig() accepted %s", name)
			}
		})
	}
}

func TestDialectorFor(t *testing.T) {
	tests := []struct {
		dbType DBType
		dsn    string
		name   string
	}{
		{dbMySQL, "user:pass@tcp(localhost:3306)/db", "mysql"},
		{dbPostgres, "host=localhost user=test dbname=test", "postgres"},
		{dbSQLite, "fixture.db", "sqlite"},
		{dbSQLServer, "sqlserver://user:pass@localhost:1433?database=test", "sqlserver"},
		{dbClickHouse, "clickhouse://localhost:9000/test", "clickhouse"},
	}
	for _, tt := range tests {
		t.Run(string(tt.dbType), func(t *testing.T) {
			dialector, err := dialectorFor(tt.dbType, tt.dsn)
			if err != nil {
				t.Fatalf("dialectorFor(): %v", err)
			}
			if got := dialector.Name(); got != tt.name {
				t.Fatalf("dialector name = %q, want %q", got, tt.name)
			}
		})
	}
	if _, err := dialectorFor(dbSQLite, ""); err == nil {
		t.Fatal("dialectorFor() accepted an empty DSN")
	}
	if _, err := dialectorFor(DBType("oracle"), "dsn"); err == nil || !strings.Contains(err.Error(), "unknown db") {
		t.Fatalf("unknown db error = %v", err)
	}
}

func TestRunGeneratesParsableSQLiteFixture(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "schema.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("open SQLite fixture: %v", err)
	}
	if err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)").Error; err != nil {
		t.Fatalf("create SQLite fixture: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close SQLite fixture: %v", err)
	}

	outPath := filepath.Join(tmp, "query")
	if err := run([]string{
		"-db", "sqlite",
		"-dsn", dbPath,
		"-tables", "users",
		"-outPath", outPath,
		"-withDefaultQuery",
	}); err != nil {
		t.Fatalf("run(): %v", err)
	}

	var generated []string
	err = filepath.WalkDir(tmp, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() && strings.HasSuffix(path, ".go") {
			generated = append(generated, path)
			if _, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("parse generated Go: %v", err)
	}
	if len(generated) < 2 {
		t.Fatalf("generated files = %v, want model and query output", generated)
	}
}

func TestRunReturnsConfigurationErrors(t *testing.T) {
	if err := run([]string{"-db", "unknown", "-dsn", "value"}); err == nil {
		t.Fatal("run() accepted an unknown database")
	}
	if err := run([]string{"-db", "sqlite"}); err == nil {
		t.Fatal("run() accepted an empty DSN")
	}
}
