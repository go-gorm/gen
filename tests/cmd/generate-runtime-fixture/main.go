package main

import (
	"flag"
	"path/filepath"

	"gorm.io/gen/tests/internal/runtimefixture"
)

func main() {
	outPath := flag.String("out", filepath.Join("fixture", "query"), "generated query output directory")
	flag.Parse()
	runtimefixture.Generate(*outPath)
}
