// Package ifaces provides importable interface fixtures for parser tests.
package ifaces

// InsertMethod is a test interface for GetInterfacePath.
type InsertMethod interface {
	AddUser(name string, age int) error
}

// TestIF is a test interface for GetInterfacePath.
type TestIF interface {
	FindByID(id int) int
}
