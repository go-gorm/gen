# Repository Guidelines

## Project Structure & Module Organization

Public generator APIs live in root-level `.go` files. Query expressions are under `field/`, clause helpers under `helper/`, and implementation details under `internal/` (`generate`, `parser`, `template`, `model`, and `diagnostic`). Unit tests sit beside their packages.

`examples/`, `tests/`, and `tools/gentool/` are independent Go modules with their own `go.mod` files. `examples/dal/` contains checked-in generated output. `tests/.expect/` stores golden code-generation fixtures; update these only when output changes intentionally.

## Build, Test, and Development Commands

- `go test ./...` runs the root module’s unit tests, matching the main CI job.
- `go test -race ./...` adds race detection for concurrency-sensitive changes.
- `golangci-lint run --timeout 5m` runs the pull-request lint configuration.
- `./examples/generate.sh` regenerates the selected example; edit `TARGET_DIR` in the script to choose a scenario.
- `docker compose -f tests/docker-compose.yml up -d` starts databases needed by integration tests.
- `GORM_DIALECT=mysql ./tests/test.sh` runs integration tests with race detection. Set `GEN_DSN` for non-default database settings.

Run commands from the repository root unless noted. For nested-module changes, also run `go test ./...` inside that module.

## Coding Style & Naming Conventions

Use standard Go formatting and tabs; run `goimports` on edited Go files. Package names are short, lowercase nouns. Exported identifiers use `PascalCase`, local identifiers use `camelCase`, and tests use `*_test.go`. Preserve Go 1.18 compatibility unless maintainers change it.

The lint configuration also enables `bodyclose`, `revive`, and `unparam`. Do not hand-edit generated `.gen.go` files unless the change specifically tests generated output.

## Testing Guidelines

Add focused table-driven tests near the code under test. Generator changes should cover behavior and emitted output; update applicable golden fixtures and verify integration generation. There is no stated coverage threshold, but bug fixes should include regression tests. Database tests must clean up created state.

## Commit & Pull Request Guidelines

Recent history generally uses concise, imperative Conventional Commit prefixes such as `feat:`, `fix:`, `test:`, `docs:`, `refactor:`, and scoped forms like `feat(gentool):`. Keep commits focused.

Pull requests should explain the problem and solution, link related issues, list verification commands, and call out generated-file or compatibility changes. Include screenshots only for documentation or UI changes where they clarify the result.

## Security & Configuration

Never commit real DSNs, passwords, or production schemas. Supply database credentials through environment variables and use disposable local containers for integration testing.
