package gen

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"gorm.io/gen/internal/generate"
)

func TestOutputWithManifest_IncrementalSkipDoesNotOverwrite(t *testing.T) {
	tmp := t.TempDir()
	g := NewGenerator(Config{OutPath: tmp})
	g.Incremental = true

	m := &genManifest{Version: 1, Files: map[string]string{}}
	var mu sync.Mutex

	fileName := filepath.Join(tmp, "x.go")
	content := []byte("package p\n\nfunc A() {}\n")

	if err := g.outputWithManifest(fileName, content, m, filepath.Base(fileName), &mu); err != nil {
		t.Fatalf("first output: %v", err)
	}

	if err := os.WriteFile(fileName, []byte("package p\n\nfunc B() {}\n"), 0640); err != nil {
		t.Fatalf("tamper file: %v", err)
	}

	if err := g.outputWithManifest(fileName, content, m, filepath.Base(fileName), &mu); err != nil {
		t.Fatalf("second output: %v", err)
	}

	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.Contains(string(b), "func B") {
		t.Fatalf("expected file to keep tampered content, got:\n%s", string(b))
	}
}

func TestBuildMergedQueryData_MergeKeepsPreviousTables(t *testing.T) {
	tmp := t.TempDir()
	g := NewGenerator(Config{OutPath: tmp})
	g.MergeQuery = true

	g.Data = map[string]*genInfo{
		"UserB": {QueryStructMeta: &generate.QueryStructMeta{ModelStructName: "UserB", QueryStructName: "userB", FileName: "user_b"}},
	}

	if err := os.WriteFile(filepath.Join(tmp, "user_a.gen.go"), []byte("package query\n"), 0640); err != nil {
		t.Fatalf("write user_a.gen.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "user_b.gen.go"), []byte("package query\n"), 0640); err != nil {
		t.Fatalf("write user_b.gen.go: %v", err)
	}

	manifest := &genManifest{
		Version: 1,
		Tables: map[string]genManifestTable{
			"UserA": {ModelStructName: "UserA", QueryStructName: "userA", FileName: "user_a"},
		},
		Files: map[string]string{},
	}

	mergedTables, dataForGenGo := g.buildMergedQueryData(manifest)

	if _, ok := mergedTables["UserA"]; !ok {
		t.Fatalf("expected merged tables to contain UserA")
	}
	if _, ok := mergedTables["UserB"]; !ok {
		t.Fatalf("expected merged tables to contain UserB")
	}
	if _, ok := dataForGenGo["UserA"]; !ok {
		t.Fatalf("expected merged data to contain placeholder UserA")
	}
	if _, ok := dataForGenGo["UserB"]; !ok {
		t.Fatalf("expected merged data to contain current UserB")
	}
}

