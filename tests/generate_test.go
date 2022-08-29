package tests_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"gorm.io/gen"
)

var _ = os.Setenv("GORM_DIALECT", "mysql")

const (
	generateDirPrefix = ".gen/"
	expectDirPrefix   = ".expect/"
)

func TestGenerate_all(t *testing.T) {
	dir := generateDirPrefix + "dal_1"

	g := gen.NewGenerator(gen.Config{
		OutPath: dir + "/query",
		Mode:    gen.WithDefaultQuery,
	})

	g.UseDB(DB)

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()

	_ = os.Remove(dir + "/query/gen_test.db")

	err := matchGeneratedFile(dir)
	if err != nil {
		t.Errorf("generated file is unexpected: %s", err)
	}
}

func TestGenerate_ultimate(t *testing.T) {
	dir := generateDirPrefix + "dal_2"

	g := gen.NewGenerator(gen.Config{
		OutPath: dir + "/query",
		Mode:    gen.WithDefaultQuery,

		WithUnitTest: true,

		FieldNullable:     true,
		FieldCoverable:    true,
		FieldWithIndexTag: true,
	})

	g.UseDB(DB)

	g.WithJSONTagNameStrategy(func(c string) string { return "-" })

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()

	_ = os.Remove(dir + "/query/gen_test.db")

	err := matchGeneratedFile(dir)
	if err != nil {
		t.Errorf("generated file is unexpected: %s", err)
	}
}

func matchGeneratedFile(dir string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	expectDir := expectDirPrefix + strings.TrimPrefix(dir, generateDirPrefix)
	diffResult, err := exec.CommandContext(ctx, "diff", "-r", expectDir, dir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("diff %s and %s got: %w\n%s", expectDir, dir, err, diffResult)
	}
	if len(diffResult) != 0 {
		return fmt.Errorf("unexpected content: %s", diffResult)
	}
	return nil
}
