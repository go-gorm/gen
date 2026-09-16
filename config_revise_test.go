package gen

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfigReviseEmptyOutPath is a regression test for an empty OutPath: it used
// to be passed to filepath.Abs first, which resolves "" to the current working
// directory and returns a nil error, so the "./query/" default was never applied
// and the query package name was derived from the working directory.
func TestConfigReviseEmptyOutPath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{}
	if err := cfg.Revise(); err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(wd, "query")
	if cfg.OutPath != want {
		t.Errorf("OutPath = %q, want %q", cfg.OutPath, want)
	}
	if cfg.queryPkgName != "query" {
		t.Errorf("queryPkgName = %q, want %q", cfg.queryPkgName, "query")
	}
	if cfg.OutFile != filepath.Join(want, "gen.go") {
		t.Errorf("OutFile = %q, want %q", cfg.OutFile, filepath.Join(want, "gen.go"))
	}
}

// TestConfigReviseExplicitOutPath ensures an explicitly configured OutPath is
// still honoured and merely made absolute.
func TestConfigReviseExplicitOutPath(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{OutPath: dir}
	if err := cfg.Revise(); err != nil {
		t.Fatal(err)
	}

	if cfg.OutPath != dir {
		t.Errorf("OutPath = %q, want %q", cfg.OutPath, dir)
	}
	if cfg.queryPkgName != filepath.Base(dir) {
		t.Errorf("queryPkgName = %q, want %q", cfg.queryPkgName, filepath.Base(dir))
	}
}
