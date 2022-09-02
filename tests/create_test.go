package tests_test

import (
	"testing"

	"gorm.io/gen/tests/.expect/dal_test/model"
	"gorm.io/gen/tests/.expect/dal_test/query"
)

func TestCreate(t *testing.T) {
	useOnce.Do(CRUDInit)

	u := query.User

	err := u.WithContext(ctx).Create(&model.User{ID: 1})
	if err != nil {
		t.Errorf("create model fail: %s", err)
	}
}
