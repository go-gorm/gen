package main

import (
	"context"
	"fmt"

	"gorm.io/gen/examples/biz"
	"gorm.io/gen/examples/conf"
	"gorm.io/gen/examples/dal"
	"gorm.io/gen/examples/dal/query"
)

func init() {
	dal.DB = dal.ConnectDB(conf.MySQLDSN).Debug()
}

func main() {
	// start your project here
	fmt.Println("hello world")
	defer fmt.Println("bye~")

	query.SetDefault(dal.DB)
	biz.Query(context.Background())
}
