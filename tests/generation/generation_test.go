package generation_test

import (
	"path/filepath"
	"testing"

	"gorm.io/gen/tests/internal/golden"
	"gorm.io/gen/tests/internal/runtimefixture"
)

func TestRuntimeFixtureMatchesGenerator(t *testing.T) {
	gotDir := filepath.Join(t.TempDir(), "query")
	runtimefixture.Generate(gotDir)

	wantDir := filepath.Join("..", "fixture", "query")
	if err := golden.CompareDirs(wantDir, gotDir); err != nil {
		t.Fatalf("runtime query fixture is stale: %v", err)
	}
}
