package ifaces

type InsertMethod interface {
	AddUser(name string, age int) error
}

type TestIF interface {
	FindByID(id int) int
}
