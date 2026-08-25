// Package gen generates type-safe data-access code from database schemas and Go models.
//
// A Generator owns the generation configuration and model metadata. Applications
// register models or tables, optionally bind custom interfaces, and call Execute
// to write query objects and model code.
package gen
