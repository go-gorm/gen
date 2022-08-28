package tests_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"gorm.io/gen"
)

var _ = os.Setenv("GORM_DIALECT", "mysql")

func TestGenerate_all(t *testing.T) {
	dir := "dal_1"

	g := gen.NewGenerator(gen.Config{
		OutPath: dir + "/query",
		Mode:    gen.WithDefaultQuery,
	})

	g.UseDB(DB)

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()

	err := matchGeneratedFile(dir)
	if err != nil {
		t.Errorf("generated file is unexpected: %s", err)
	}
}

func TestGenerate_ultimate(t *testing.T) {
	dir := "dal_2"

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

	err := matchGeneratedFile(dir)
	if err != nil {
		t.Errorf("generated file is unexpected: %s", err)
	}
}

func matchGeneratedFile(dir string) error {
	// walkFunc := func(path string, info os.FileInfo, _ error) error {
	// 	// skip dir
	// 	if info.IsDir() {
	// 		return nil
	// 	}

	// 	generatePath := strings.TrimPrefix(path, "exepct/")
	// 	fmt.Println("expected: ", path)
	// 	fmt.Println("generated: ", generatePath)

	// 	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	// 	defer cancel()

	// 	diffResult, err := exec.CommandContext(ctx, "diff", path, generatePath).CombinedOutput()
	// 	if err != nil {
	// 		errs = append(errs, fmt.Errorf("diff %s and %s fail: %w", path, generatePath, err))
	// 	}
	// 	if len(diffResult) != 0 {
	// 		errs = append(errs, fmt.Errorf("unexpected content: %s", diffResult))
	// 	}
	// 	return nil
	// }

	// filepath.Walk("./expect/"+dir, walkFunc)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	diffResult, err := exec.CommandContext(ctx, "diff", "-r", ".expect/"+dir, dir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("diff %s and %s fail: %w", ".expect/"+dir, dir, err)
	}
	if len(diffResult) != 0 {
		return fmt.Errorf("unexpected content: %s", diffResult)
	}
	return nil
}
