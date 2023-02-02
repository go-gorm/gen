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

	"gorm.io/gen/tests/diy_method"
)

const (
	generateDirPrefix = ".gen/"
	expectDirPrefix   = ".expect/"
)

var _ = os.Setenv("GORM_DIALECT", "mysql")

var generateCase = map[string]func(dir string) *gen.Generator{
	generateDirPrefix + "dal_1": func(dir string) *gen.Generator {
		g := gen.NewGenerator(gen.Config{
			OutPath: dir + "/query",
			Mode:    gen.WithDefaultQuery,
		})
		g.UseDB(DB)
		g.ApplyBasic(g.GenerateAllTable()...)
		return g
	},
	generateDirPrefix + "dal_2": func(dir string) *gen.Generator {
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
		return g
	},
	generateDirPrefix + "dal_3": func(dir string) *gen.Generator {
		g := gen.NewGenerator(gen.Config{
			OutPath: dir + "/query",
			Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,

			WithUnitTest: true,

			FieldNullable:     true,
			FieldCoverable:    true,
			FieldWithIndexTag: true,
		})
		g.UseDB(DB)
		g.WithJSONTagNameStrategy(func(c string) string { return "-" })
		g.ApplyBasic(g.GenerateAllTable()...)
		return g
	},
	generateDirPrefix + "dal_4": func(dir string) *gen.Generator {
		g := gen.NewGenerator(gen.Config{
			OutPath: dir + "/query",
			Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,

			WithUnitTest: true,

			FieldNullable:     true,
			FieldCoverable:    true,
			FieldWithIndexTag: true,
		})
		g.UseDB(DB)
		g.WithJSONTagNameStrategy(func(c string) string { return "-" })
		g.ApplyBasic(g.GenerateAllTable()...)
		g.ApplyInterface(func(testIF diy_method.TestIF, testFor diy_method.TestFor, method diy_method.InsertMethod, selectMethod diy_method.SelectMethod) {
		}, g.GenerateModel("users"))
		return g
	},
}

func TestGenerate(t *testing.T) {
	for dir := range generateCase {
		t.Run("TestGenerate_"+dir, func(dir string) func(t *testing.T) {
			return func(t *testing.T) {
				t.Parallel()
				if err := matchGeneratedFile(dir); err != nil {
					t.Errorf("generated file is unexpected: %s", err)
				}
			}
		}(dir))
	}
}

func matchGeneratedFile(dir string) error {
	_ = os.Remove(dir + "/query/gen_test.db")

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

func TestGenerate_expect(t *testing.T) {
	if os.Getenv("GEN_EXPECT") == "" {
		t.SkipNow()
	}
	g := gen.NewGenerator(gen.Config{
		OutPath: expectDirPrefix + "dal_test" + "/query",
		Mode:    gen.WithDefaultQuery,
	})
	g.UseDB(DB)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
