package tests_test

import (
	"context"
	"sync"

	"gorm.io/gen/tests/.expect/dal_test/query"
	relquery "gorm.io/gen/tests/.expect/dal_test_relation/query"
)

var useOnce sync.Once
var ctx = context.Background()

func CRUDInit() {
	query.Use(DB)
	query.SetDefault(DB)
	relquery.Use(DB)
	relquery.SetDefault(DB)
}
