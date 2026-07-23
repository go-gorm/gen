package main

import "gorm.io/gen"

type MainMethod interface {
	// FindLocal
	//
	// select * from @@table where id=@id
	FindLocal(id int) gen.T
}

func main() {}
