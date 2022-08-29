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
	err := db.AutoMigrate(&model.Player{})
	if err != nil {
		fmt.Printf("Error: AutoMigrate(&model.Player{}) fail: %s", err)
	}
}

func Test_playerQuery(t *testing.T) {
	player := newPlayer(db)
	player = *player.As(player.TableName())
	_do := player.WithContext(context.Background()).Debug()

	primaryKey := field.NewString(player.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <players> fail:", err)
		return
	}

	_, ok := player.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from player success")
	}

	err = _do.Create(&model.Player{})
	if err != nil {
		t.Error("create item in table <players> fail:", err)
	}

	err = _do.Save(&model.Player{})
	if err != nil {
		t.Error("create item in table <players> fail:", err)
	}

	err = _do.CreateInBatches([]*model.Player{{}, {}}, 10)
	if err != nil {
		t.Error("create item in table <players> fail:", err)
	}

	_, err = _do.Select(player.ALL).Take()
	if err != nil {
		t.Error("Take() on table <players> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <players> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <players> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <players> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*model.Player{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <players> fail:", err)
	}

	_, err = _do.Select(player.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <players> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <players> fail:", err)
	}

	_, err = _do.Select(player.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <players> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <players> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <players> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <players> fail:", err)
	}

	_, err = _do.ScanByPage(&model.Player{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <players> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <players> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <players> fail:", err)
	}

	var _a _another
	var _aPK = field.NewString(_a.TableName(), clause.PrimaryKey)

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <players> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <players> fail:", err)
	}

	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <players> fail:", err)
	}
}
