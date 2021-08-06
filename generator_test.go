package gen

import (
	"testing"
)

func TestConfig(t *testing.T) {
	_ = &Config{
		db: nil,

		OutPath: "path",
		OutFile: "",

		ModelPkgName: "models",

		queryPkgName: "query",
	}
}