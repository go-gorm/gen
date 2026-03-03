# GORM Gen

Friendly & Safer GORM powered by Code Generation.

[![Release](https://img.shields.io/github/v/release/go-gorm/gen)](https://github.com/go-gorm/gen/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-gorm/gen)](https://goreportcard.com/report/github.com/go-gorm/gen)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![OpenIssue](https://img.shields.io/github/issues/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aopen+is%3Aissue)
[![ClosedIssue](https://img.shields.io/github/issues-closed/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aissue+is%3Aclosed)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/go-gorm/gen)](https://www.tickgit.com/browse?repo=github.com/go-gorm/gen)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/gorm.io/gen?tab=doc)

## Overview

- Generate idiomatic, reusable DAO APIs from database schema and/or interface-based SQL templates
- Type-safe query DSL (fields, conditions, assignments) with strong static typing (including generics mode)
- Database-to-struct follows GORM conventions (tags, nullable/default/unsigned/index/type, etc.)
- Built on top of GORM: use the same plugins, dialectors, and ecosystem you already have

## Documentation

- Gen Guides (official site): https://gorm.io/gen/index.html
- GORM Guides: https://gorm.io/docs

## Quick Start

Install as a library:

```bash
go get gorm.io/gen@latest
```

### 1) Generate code

Create a generator entry (recommended: `cmd/gen/main.go`):

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("dsn"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/dal/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface, // enable query.SetDefault(db)
	})

	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
```

Run generation:

```bash
go run ./cmd/gen
```

### 2) Use generated code

If you enabled `WithDefaultQuery`, initialize once at startup:

```go
package main

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"your/module/internal/dal/query"
)

func main() {
	db, err := gorm.Open(mysql.Open("dsn"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	query.SetDefault(db)

	u := query.User
	_, _ = u.WithContext(context.Background()).
		Where(u.Age.Gt(18)).
		Find()
}
```

`query.SetDefault(db)` only needs to run once (e.g. during service startup). After that, use `query.<Table>` directly.

If you don’t want a global default, skip `WithDefaultQuery` and use:

```go
q := query.Use(db)
_, _ = q.User.WithContext(ctx).Where(q.User.Age.Gt(18)).Find()
```

More runnable examples: [examples](./examples)

## Common Setups

Gen has one generator entry point (`gen.NewGenerator(gen.Config{...})`). The main knobs you typically use are:

- What to generate: DB schema → models/query; plus optional interface-SQL methods
- How the query API looks: `Config.Mode` flags (`WithDefaultQuery`, `WithoutContext`, `WithQueryInterface`, `WithGeneric`)

### Setup A: DB schema → model + query (recommended baseline)

```go
g := gen.NewGenerator(gen.Config{
	OutPath: "internal/dal/query",
	Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,

	FieldNullable:     true,
	FieldCoverable:    true,
	FieldWithIndexTag: true,
})
g.UseDB(db)
g.ApplyBasic(g.GenerateAllTable()...)
g.Execute()
```

### Setup B: Interface SQL templates → reusable typed methods

Define an interface with SQL comments/templates:

```go
package dal

import "gorm.io/gen"

type UserMethods interface {
	// FindByID
	//
	// SELECT * FROM users WHERE id=@id
	FindByID(id int) gen.T

	// FindByOptionalName
	//
	// SELECT * FROM users
	// {{where}}
	// {{if name != ""}}
	// name=@name
	// {{end}}
	// {{end}}
	FindByOptionalName(name string) []gen.T
}
```

Bind the interface to a generated model/table:

```go
g.ApplyInterface(func(m UserMethods) {}, g.GenerateModel("users"))
```

Template syntax reference exists in the test corpus: [method.go](./tests/diy_method/method.go).

### Setup C: Generics query API

Generics is not a separate “workflow”; it changes the generated query API surface for stronger typing.

```go
g := gen.NewGenerator(gen.Config{
	OutPath: "internal/dal/query",
	Mode:    gen.WithDefaultQuery | gen.WithGeneric,
})
```

## Recommended Project Layout

Keep code generation isolated and reproducible (works well with `go:generate` and CI).

```
your-repo/
  internal/
    dal/
      model/        # generated model structs (optional)
      query/        # generated query (+ DIY methods)
  cmd/
    gen/
      main.go       # generator entry (checked in)
```

## CLI Tool

If you prefer a CLI workflow, use GenTool:

- [GenTool README](./tools/gentool/README.md)
- [GenTool README (ZH-CN)](./tools/gentool/README.ZH_CN.md)

## Maintainers

[@riverchu](https://github.com/riverchu) [@iDer](https://github.com/idersec) [@qqxhb](https://github.com/qqxhb) [@dino-ma](https://github.com/dino-ma)

[@jinzhu](https://github.com/jinzhu)

## Contributing

[You can help to deliver a better GORM/Gen, check out things you can do](https://gorm.io/contribute.html)

## License

Released under the [MIT License](https://github.com/go-gorm/gen/blob/master/License)
