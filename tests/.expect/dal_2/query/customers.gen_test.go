// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"fmt"
	"testing"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/tests/dal_2/model"
	"gorm.io/gorm/clause"
)

func init() {
	InitializeDB()
	err := db.AutoMigrate(&model.Customer{})
	if err != nil {
		fmt.Printf("Error: AutoMigrate(&model.Customer{}) fail: %s", err)
	}
}

func Test_customerQuery(t *testing.T) {
	customer := newCustomer(db)
	customer = *customer.As(customer.TableName())
	_do := customer.WithContext(context.Background()).Debug()

	primaryKey := field.NewString(customer.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <customers> fail:", err)
		return
	}

	_, ok := customer.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from customer success")
	}

	err = _do.Create(&model.Customer{})
	if err != nil {
		t.Error("create item in table <customers> fail:", err)
	}

	err = _do.Save(&model.Customer{})
	if err != nil {
		t.Error("create item in table <customers> fail:", err)
	}

	err = _do.CreateInBatches([]*model.Customer{{}, {}}, 10)
	if err != nil {
		t.Error("create item in table <customers> fail:", err)
	}

	_, err = _do.Select(customer.ALL).Take()
	if err != nil {
		t.Error("Take() on table <customers> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <customers> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <customers> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <customers> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*model.Customer{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <customers> fail:", err)
	}

	_, err = _do.Select(customer.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <customers> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <customers> fail:", err)
	}

	_, err = _do.Select(customer.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <customers> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <customers> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <customers> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <customers> fail:", err)
	}

	_, err = _do.ScanByPage(&model.Customer{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <customers> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <customers> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <customers> fail:", err)
	}

	var _a _another
	var _aPK = field.NewString(_a.TableName(), clause.PrimaryKey)

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <customers> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <customers> fail:", err)
	}

	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <customers> fail:", err)
	}
}
