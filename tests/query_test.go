package tests_test

import (
	"testing"

	"gorm.io/gen/tests/.expect/dal_test/model"
	"gorm.io/gen/tests/.expect/dal_test/query"
)

func TestQuery_Find(t *testing.T) {
	useOnce.Do(CRUDInit)

	u := query.User

	err := u.WithContext(ctx).Create(&model.User{ID: 100})
	if err != nil {
		t.Errorf("create model fail: %s", err)
	}

	user, err := u.WithContext(ctx).Where(u.ID.Eq(100)).Take()
	if err != nil {
		t.Errorf("take model fail: %s", err)
	}
	if user.ID != 100 {
		t.Errorf("take model fail: %+v", user)
	}
	t.Logf("got model: %+v", user)
}
