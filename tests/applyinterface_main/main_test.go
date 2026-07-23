package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/gen"
	testhelper "gorm.io/gen/tests/applyinterface_main/helper"
	"gorm.io/gen/tests/applyinterface_main/model"
)

func TestApplyInterfaceMainPackageThroughExternalHelper(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "query")
	g := gen.NewGenerator(gen.Config{OutPath: outPath})

	testhelper.Apply(g, func(MainMethod) {}, model.MainUser{})
	g.Execute()

	file := filepath.Join(outPath, "main_users.gen.go")
	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read generated file %s: %v", file, err)
	}
	if !strings.Contains(string(got), "func (m mainUserDo) FindLocal") {
		t.Fatalf("generated file is missing FindLocal method:\n%s", got)
	}
	if strings.Contains(string(got), "FindWrongPackage") {
		t.Fatalf("generated file used helper package interface instead of package main interface:\n%s", got)
	}
}
