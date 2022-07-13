package gen

import "strings"

var importPkgs = importPkgS{}.Add(
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
	"",
)

var unitTestImportPkgs = importPkgS{}.Add(
	"context",
	"fmt",
	"strconv",
	"testing",
	"",
	"gorm.io/driver/sqlite",
	"gorm.io/gorm",
	"",
)

type importPkgS struct {
	paths []string
}

func (ip importPkgS) Add(paths ...string) *importPkgS {
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p != "" && p[len(p)-1] != '"' {
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
			ip.paths = append(ip.paths, p)
		}
	}
	ip.paths = append(ip.paths, "")
	return &ip
}

func (ip importPkgS) Output() []string {
	return ip.paths
}
