package gen

import "strings"

var (
	importList = new(importPkgS).Add(
		"context",
		"database/sql",
		"strings",
		"",
		"gorm.io/gorm",
		"gorm.io/gorm/schema",
		"gorm.io/gorm/clause",
		"",
		"gorm.io/gen",
		"gorm.io/gen/field",
		"gorm.io/gen/helper",
		"",
		"gorm.io/plugin/dbresolver",
	)
	unitTestImportList = new(importPkgS).Add(
		"context",
		"fmt",
		"strconv",
		"testing",
		"",
		"gorm.io/driver/sqlite",
		"gorm.io/gorm",
	)
)

type importPkgS struct {
	paths []string
}

func (ip importPkgS) Add(paths ...string) *importPkgS {
	purePaths := make([]string, 0, len(paths)+1)
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			purePaths = append(purePaths, p)
			continue
		}

		if p[len(p)-1] != '"' {
			p = `"` + p + `"`
		}

		var exists bool
		for _, existsP := range ip.paths {
			if p == existsP {
				exists = true
				break
			}
		}
		if !exists {
			purePaths = append(purePaths, p)
		}
	}
	purePaths = append(purePaths, "")

	ip.paths = append(ip.paths, purePaths...)

	return &ip
}

func (ip importPkgS) Paths() []string { return ip.paths }
