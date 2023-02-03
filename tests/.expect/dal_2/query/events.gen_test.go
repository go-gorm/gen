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
	"gorm.io/gen/tests/.gen/dal_2/model"
	"gorm.io/gorm/clause"
)

func init() {
	InitializeDB()
	err := db.AutoMigrate(&model.Event{})
	if err != nil {
		fmt.Printf("Error: AutoMigrate(&model.Event{}) fail: %s", err)
	}
}

func Test_eventQuery(t *testing.T) {
	event := newEvent(db)
	event = *event.As(event.TableName())
	_do := event.WithContext(context.Background()).Debug()

	primaryKey := field.NewString(event.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <events> fail:", err)
		return
	}

	_, ok := event.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from event success")
	}

	err = _do.Create(&model.Event{})
	if err != nil {
		t.Error("create item in table <events> fail:", err)
	}

	err = _do.Save(&model.Event{})
	if err != nil {
		t.Error("create item in table <events> fail:", err)
	}

	err = _do.CreateInBatches([]*model.Event{{}, {}}, 10)
	if err != nil {
		t.Error("create item in table <events> fail:", err)
	}

	_, err = _do.Select(event.ALL).Take()
	if err != nil {
		t.Error("Take() on table <events> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <events> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <events> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <events> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*model.Event{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <events> fail:", err)
	}

	_, err = _do.Select(event.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <events> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <events> fail:", err)
	}

	_, err = _do.Select(event.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <events> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <events> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <events> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <events> fail:", err)
	}

	_, err = _do.ScanByPage(&model.Event{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <events> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <events> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <events> fail:", err)
	}

	var _a _another
	var _aPK = field.NewString(_a.TableName(), clause.PrimaryKey)

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <events> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <events> fail:", err)
	}

	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <events> fail:", err)
	}
}
