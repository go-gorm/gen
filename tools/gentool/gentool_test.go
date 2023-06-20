package main

import (
	"fmt"
	"testing"
)

var params *CmdParams

func init() {
	testing.Init() // fix err=> flag provided but not defined: -c

	params = argParse()
}

// go test -v -run  PgConfig --args -c "./pg-gen.yml"
func TestPgConfig(t *testing.T) {
	// flag.Set("c", "pg-gen.yml")
	fmt.Println(params)

	// output
	// &{postgresql://postgres:root@localhost:5432/postgres postgres [user corp] false /tmp/db  true  true true false true}
}
