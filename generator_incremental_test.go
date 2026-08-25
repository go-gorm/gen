package gen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"gorm.io/gen/internal/generate"
)

type ManifestUserFixture struct {
	ID   uint
	Name string
}

type ManifestOrderFixture struct {
	ID     uint
	UserID uint
}

func newManifestGenerator(t *testing.T, outPath string, incremental, merge bool, mode GenerateMode, models ...interface{}) *Generator {
	t.Helper()
	g := NewGenerator(Config{
		OutPath:     outPath,
		Incremental: incremental,
		MergeQuery:  merge,
		Mode:        mode,
	})
	g.ApplyBasic(models...)
	return g
}

func executePanic(g *Generator) (panicValue interface{}) {
	defer func() {
		panicValue = recover()
	}()
	g.Execute()
	return nil
}

func readManifest(t *testing.T, dir string) *genManifest {
	t.Helper()
	m, _, err := loadManifest(dir)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	return m
}

func assertManifestHashes(t *testing.T, dir string, manifest *genManifest) {
	t.Helper()
	for fileName, wantHash := range manifest.Files {
		content, err := os.ReadFile(filepath.Join(dir, fileName))
		if err != nil {
			t.Fatalf("read generated file %q: %v", fileName, err)
		}
		if gotHash := sha256Hex(content); gotHash != wantHash {
			t.Fatalf("manifest hash for %q = %q, want %q", fileName, gotHash, wantHash)
		}
	}
}

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

func TestIncrementalExecutePreservesUserModifiedGeneratedFile(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "query")
	first := newManifestGenerator(t, outPath, true, false, WithDefaultQuery, ManifestUserFixture{})
	if panicValue := executePanic(first); panicValue != nil {
		t.Fatalf("first Execute() panicked: %v", panicValue)
	}

	manifest := readManifest(t, outPath)
	if len(manifest.Files) != 2 {
		t.Fatalf("manifest files = %v, want table query and gen.go", manifest.Files)
	}
	assertManifestHashes(t, outPath, manifest)

	tableFile := ""
	for fileName := range manifest.Files {
		if fileName != filepath.Base(first.OutFile) {
			tableFile = fileName
			break
		}
	}
	if tableFile == "" {
		t.Fatal("manifest does not contain a table query file")
	}

	modified := []byte("package query\n\n// user-owned modification\n")
	if err := os.WriteFile(filepath.Join(outPath, tableFile), modified, 0640); err != nil {
		t.Fatalf("modify generated file: %v", err)
	}
	second := newManifestGenerator(t, outPath, true, false, WithDefaultQuery, ManifestUserFixture{})
	if panicValue := executePanic(second); panicValue != nil {
		t.Fatalf("second Execute() panicked: %v", panicValue)
	}
	got, err := os.ReadFile(filepath.Join(outPath, tableFile))
	if err != nil {
		t.Fatalf("read modified generated file: %v", err)
	}
	if !reflect.DeepEqual(got, modified) {
		t.Fatalf("incremental generation overwrote a modified file:\n%s", got)
	}
}

func TestMergeQueryExecuteKeepsPreviouslyGeneratedTables(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "query")
	first := newManifestGenerator(t, outPath, false, true, WithDefaultQuery, ManifestUserFixture{})
	if panicValue := executePanic(first); panicValue != nil {
		t.Fatalf("first Execute() panicked: %v", panicValue)
	}
	second := newManifestGenerator(t, outPath, false, true, WithDefaultQuery, ManifestOrderFixture{})
	if panicValue := executePanic(second); panicValue != nil {
		t.Fatalf("second Execute() panicked: %v", panicValue)
	}

	manifest := readManifest(t, outPath)
	if len(manifest.Tables) != 2 {
		t.Fatalf("manifest tables = %#v, want two tables", manifest.Tables)
	}
	for _, name := range []string{"ManifestUserFixture", "ManifestOrderFixture"} {
		if _, ok := manifest.Tables[name]; !ok {
			t.Fatalf("manifest is missing %q: %#v", name, manifest.Tables)
		}
	}
	assertManifestHashes(t, outPath, manifest)

	rootQuery, err := os.ReadFile(filepath.Join(outPath, "gen.go"))
	if err != nil {
		t.Fatalf("read merged gen.go: %v", err)
	}
	for _, queryName := range []string{"manifestUserFixture", "manifestOrderFixture"} {
		if !strings.Contains(string(rootQuery), queryName) {
			t.Fatalf("merged gen.go does not contain %q:\n%s", queryName, rootQuery)
		}
	}
}

func TestMergeQueryModeMismatchDoesNotChangeExistingFiles(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "query")
	first := newManifestGenerator(t, outPath, false, true, WithDefaultQuery, ManifestUserFixture{})
	if panicValue := executePanic(first); panicValue != nil {
		t.Fatalf("first Execute() panicked: %v", panicValue)
	}
	beforeGen, err := os.ReadFile(filepath.Join(outPath, "gen.go"))
	if err != nil {
		t.Fatalf("read gen.go: %v", err)
	}
	beforeManifest, err := os.ReadFile(filepath.Join(outPath, manifestFileName))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	second := newManifestGenerator(t, outPath, false, true, WithQueryInterface, ManifestOrderFixture{})
	if panicValue := executePanic(second); panicValue == nil {
		t.Fatal("Execute() did not panic for a merge mode mismatch")
	}
	afterGen, err := os.ReadFile(filepath.Join(outPath, "gen.go"))
	if err != nil {
		t.Fatalf("read gen.go after mismatch: %v", err)
	}
	afterManifest, err := os.ReadFile(filepath.Join(outPath, manifestFileName))
	if err != nil {
		t.Fatalf("read manifest after mismatch: %v", err)
	}
	if !reflect.DeepEqual(afterGen, beforeGen) || !reflect.DeepEqual(afterManifest, beforeManifest) {
		t.Fatal("merge mode mismatch changed existing output")
	}
}

func TestLoadManifestNormalizesMissingFieldsAndRejectsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, manifestFileName)
	if err := os.WriteFile(path, []byte(`{"mode":1}`), 0640); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	manifest := readManifest(t, dir)
	if manifest.Version != 1 || manifest.Tables == nil || manifest.Files == nil {
		t.Fatalf("manifest was not normalized: %#v", manifest)
	}

	invalid := []byte(`{"version":`)
	if err := os.WriteFile(path, invalid, 0640); err != nil {
		t.Fatalf("write invalid manifest: %v", err)
	}
	if _, _, err := loadManifest(dir); err == nil {
		t.Fatal("loadManifest() accepted invalid JSON")
	}

	missingDir := filepath.Join(t.TempDir(), "query")
	manifest, gotPath, err := loadManifest(missingDir)
	if err != nil {
		t.Fatalf("load missing manifest: %v", err)
	}
	if manifest.Version != 1 || gotPath != filepath.Join(missingDir, manifestFileName) {
		t.Fatalf("unexpected missing manifest result: manifest=%#v path=%q", manifest, gotPath)
	}

	if _, err := json.Marshal(manifest); err != nil {
		t.Fatalf("normalized manifest cannot be marshaled: %v", err)
	}
}

func TestExecuteRejectsCorruptManifestBeforeWritingGeneratedFiles(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "query")
	if err := os.MkdirAll(outPath, 0750); err != nil {
		t.Fatalf("make output directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outPath, manifestFileName), []byte("{"), 0640); err != nil {
		t.Fatalf("write invalid manifest: %v", err)
	}
	g := newManifestGenerator(t, outPath, true, false, WithDefaultQuery, ManifestUserFixture{})
	if panicValue := executePanic(g); panicValue == nil {
		t.Fatal("Execute() did not panic for a corrupt manifest")
	}
	entries, err := os.ReadDir(outPath)
	if err != nil {
		t.Fatalf("read output directory: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != manifestFileName {
		t.Fatalf("corrupt manifest left partial output: %v", entries)
	}
}
