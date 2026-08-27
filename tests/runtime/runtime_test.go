package runtime_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gen/field"
	fixture "gorm.io/gen/tests/fixture/model"
	"gorm.io/gen/tests/fixture/query"
	"gorm.io/gorm"
)

var _ func([]*string) field.AssignExpr = query.Q.User.Photos.Value

func newSQLiteDB(t *testing.T) (*gorm.DB, *query.Query) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "runtime.db")
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		t.Fatalf("open SQLite: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close SQLite: %v", err)
		}
	})
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("enable SQLite foreign keys: %v", err)
	}
	if err := db.AutoMigrate(&fixture.Company{}, &fixture.User{}, &fixture.Order{}); err != nil {
		t.Fatalf("migrate SQLite fixture: %v", err)
	}
	seedRuntimeData(t, db)
	return db, query.Use(db)
}

func seedRuntimeData(t *testing.T, db *gorm.DB) {
	t.Helper()
	companies := []fixture.Company{{ID: 1, Name: "Acme"}, {ID: 2, Name: "Globex"}}
	users := []fixture.User{
		{ID: 1, Name: "alice", Age: 30, Active: true, Role: "admin", CompanyID: 1},
		{ID: 2, Name: "bob", Age: 20, Active: false, Role: "viewer", CompanyID: 1},
		{ID: 3, Name: "carol", Age: 40, Active: true, Role: "viewer", CompanyID: 2},
	}
	orders := []fixture.Order{
		{ID: 1, UserID: 1, Amount: 100, Status: "paid"},
		{ID: 2, UserID: 1, Amount: 200, Status: "pending"},
		{ID: 3, UserID: 3, Amount: 300, Status: "paid"},
	}
	for _, value := range []interface{}{&companies, &users, &orders} {
		if err := db.Create(value).Error; err != nil {
			t.Fatalf("seed %T: %v", value, err)
		}
	}
}

func TestDIYMethodsExecuteAllControlStructures(t *testing.T) {
	_, q := newSQLiteDB(t)
	u := q.User.WithContext(context.Background())

	byName, err := u.FindByNameBranch("alice", false)
	assertUsers(t, byName, err, "alice")
	byActive, err := u.FindByNameBranch("", false)
	assertUsers(t, byActive, err, "bob")

	byNames, err := u.FindByNames([]string{"alice", "carol"})
	assertUsers(t, byNames, err, "alice", "carol")
	allUsers, err := u.FindByNames(nil)
	assertUsers(t, allUsers, err, "alice", "bob", "carol")

	byAttributes, err := u.FindByAttributes(map[string]string{"name": "alice", "role": "admin"})
	assertUsers(t, byAttributes, err, "alice")

	rows, err := u.UpdateOptional(2, "bobby", 25)
	if err != nil || rows != 1 {
		t.Fatalf("UpdateOptional() = (%d, %v), want (1, nil)", rows, err)
	}
	updated, err := u.Where(q.User.ID.Eq(2)).First()
	if err != nil || updated.Name != "bobby" || updated.Age != 25 {
		t.Fatalf("updated user = %#v, err=%v", updated, err)
	}
	rows, err = u.UpdateOptional(2, "", 0)
	if err != nil || rows != 1 {
		t.Fatalf("empty UpdateOptional() = (%d, %v), want (1, nil)", rows, err)
	}

	trimmed, err := u.FindWithTrim([]string{"alice", "carol"})
	assertUsers(t, trimmed, err, "alice", "carol")
	emptyTrim, err := u.FindWithTrim(nil)
	assertUsers(t, emptyTrim, err)

	rows, err = u.InsertUser("dave", 28, true, "editor", 2)
	if err != nil || rows != 1 {
		t.Fatalf("InsertUser() = (%d, %v), want (1, nil)", rows, err)
	}
	inserted, err := u.Where(q.User.Name.Eq("dave")).First()
	if err != nil || inserted.Role != "editor" || inserted.CompanyID != 2 {
		t.Fatalf("inserted user = %#v, err=%v", inserted, err)
	}
}

func TestGeneratedQueryPaginationAndAggregates(t *testing.T) {
	_, q := newSQLiteDB(t)
	u := q.User.WithContext(context.Background())

	page, err := u.Order(q.User.ID).Limit(2).Offset(1).Find()
	if err != nil || len(page) != 2 || page[0].ID != 2 || page[1].ID != 3 {
		t.Fatalf("paginated users = %#v, err=%v", page, err)
	}

	count, err := u.Distinct(q.User.Role).Count()
	if err != nil || count != 2 {
		t.Fatalf("distinct role count = %d, err=%v", count, err)
	}
	var roles []string
	if err := u.Distinct(q.User.Role).Order(q.User.Role).Pluck(q.User.Role, &roles); err != nil {
		t.Fatalf("pluck roles: %v", err)
	}
	if !reflect.DeepEqual(roles, []string{"admin", "viewer"}) {
		t.Fatalf("roles = %v", roles)
	}

	var projection []struct {
		Name string
		Age  int
	}
	if err := u.Select(q.User.Name, q.User.Age).Where(q.User.Active.Is(true)).Order(q.User.ID).Scan(&projection); err != nil {
		t.Fatalf("scan projection: %v", err)
	}
	if !reflect.DeepEqual(projection, []struct {
		Name string
		Age  int
	}{{Name: "alice", Age: 30}, {Name: "carol", Age: 40}}) {
		t.Fatalf("projection = %#v", projection)
	}

	var grouped []struct {
		Role  string
		Total int64
	}
	countExpr := q.User.ID.Count()
	if err := u.Select(q.User.Role, countExpr.As("total")).Group(q.User.Role).Having(countExpr.Gt(1)).Scan(&grouped); err != nil {
		t.Fatalf("group/having: %v", err)
	}
	if len(grouped) != 1 || grouped[0].Role != "viewer" || grouped[0].Total != 2 {
		t.Fatalf("grouped result = %#v", grouped)
	}
}

func TestGeneratedQueryBatchingScopesAndWrites(t *testing.T) {
	_, q := newSQLiteDB(t)
	u := q.User.WithContext(context.Background())

	var batchResult []*fixture.User
	var batches []int
	var visited []uint
	err := u.Order(q.User.ID).FindInBatches(&batchResult, 2, func(_ gen.Dao, batch int) error {
		batches = append(batches, batch)
		for _, user := range batchResult {
			visited = append(visited, user.ID)
		}
		return nil
	})
	if err != nil || !reflect.DeepEqual(batches, []int{1, 2}) || !reflect.DeepEqual(visited, []uint{1, 2, 3}) {
		t.Fatalf("FindInBatches() batches=%v visited=%v err=%v", batches, visited, err)
	}

	stopErr := errors.New("stop batches")
	batches = nil
	err = u.Order(q.User.ID).FindInBatches(&batchResult, 1, func(_ gen.Dao, batch int) error {
		batches = append(batches, batch)
		return stopErr
	})
	if !errors.Is(err, stopErr) || !reflect.DeepEqual(batches, []int{1}) {
		t.Fatalf("callback stop: batches=%v err=%v", batches, err)
	}

	adultScope := func(dao gen.Dao) gen.Dao { return dao.Where(q.User.Age.Gte(30)) }
	adults, err := u.Scopes(adultScope).Order(q.User.ID).Find()
	assertUsers(t, adults, err, "alice", "carol")
	all, err := u.Find()
	assertUsers(t, all, err, "alice", "bob", "carol")

	info, err := u.Where(q.User.ID.Eq(1)).Update(q.User.Age, 31)
	if err != nil || info.RowsAffected != 1 {
		t.Fatalf("Update() info=%#v err=%v", info, err)
	}
	info, err = u.Where(q.User.ID.Eq(2)).Updates(map[string]interface{}{"name": "robert", "role": "editor"})
	if err != nil || info.RowsAffected != 1 {
		t.Fatalf("Updates() info=%#v err=%v", info, err)
	}
	info, err = u.Where(q.User.ID.Eq(3)).UpdateColumn(q.User.Age, 41)
	if err != nil || info.RowsAffected != 1 {
		t.Fatalf("UpdateColumn() info=%#v err=%v", info, err)
	}

	info, err = u.Where(q.User.ID.Eq(2)).Delete()
	if err != nil || info.RowsAffected != 1 {
		t.Fatalf("Delete() info=%#v err=%v", info, err)
	}
	remaining, err := u.Order(q.User.ID).Find()
	assertUsers(t, remaining, err, "alice", "carol")
}

func TestGeneratedQueryUpdatesRawSerializerValue(t *testing.T) {
	_, q := newSQLiteDB(t)
	u := q.User.WithContext(context.Background())
	one, two := "1", "2"

	info, err := u.Where(q.User.ID.Eq(1)).UpdateSimple(q.User.Photos.Value([]*string{&one, &two}))
	if err != nil || info.RowsAffected != 1 {
		t.Fatalf("UpdateSimple() info=%#v err=%v", info, err)
	}
	updated, err := u.Where(q.User.ID.Eq(1)).First()
	if err != nil || !reflect.DeepEqual(updated.Photos, []*string{&one, &two}) {
		t.Fatalf("serialized photos = %#v, err=%v", updated.Photos, err)
	}
	matched, err := u.Where(q.User.Photos.Eq([]*string{&one, &two})).First()
	if err != nil || matched.ID != 1 {
		t.Fatalf("serialized Eq() result = %#v, err=%v", matched, err)
	}
	matched, err = u.Where(q.User.Photos.In([]*string{}, []*string{&one, &two})).First()
	if err != nil || matched.ID != 1 {
		t.Fatalf("serialized In() result = %#v, err=%v", matched, err)
	}

	if _, err = u.Where(q.User.ID.Eq(1)).UpdateSimple(q.User.Photos.Value([]*string{})); err != nil {
		t.Fatalf("update empty serializer value: %v", err)
	}
	updated, err = u.Where(q.User.ID.Eq(1)).First()
	if err != nil || updated.Photos == nil || len(updated.Photos) != 0 {
		t.Fatalf("empty serialized photos = %#v, err=%v", updated.Photos, err)
	}

	if _, err = u.Where(q.User.ID.Eq(1)).UpdateSimple(q.User.Photos.Value([]*string(nil))); err != nil {
		t.Fatalf("update nil serializer value: %v", err)
	}
	updated, err = u.Where(q.User.ID.Eq(1)).First()
	if err != nil || updated.Photos != nil {
		t.Fatalf("nil serialized photos = %#v, err=%v", updated.Photos, err)
	}
}

func TestGeneratedQueryWriteSafetyRelationsTransactionsAndErrors(t *testing.T) {
	_, q := newSQLiteDB(t)
	u := q.User.WithContext(context.Background())

	if _, err := u.Update(q.User.Age, 99); !errors.Is(err, gorm.ErrMissingWhereClause) {
		t.Fatalf("global Update() error = %v, want ErrMissingWhereClause", err)
	}
	if _, err := u.Delete(); !errors.Is(err, gorm.ErrMissingWhereClause) {
		t.Fatalf("global Delete() error = %v, want ErrMissingWhereClause", err)
	}

	withRelations, err := u.Preload(q.User.Company, q.User.Orders).Order(q.User.ID).Find()
	if err != nil {
		t.Fatalf("preload relations: %v", err)
	}
	if len(withRelations) != 3 || withRelations[0].Company.Name != "Acme" || len(withRelations[0].Orders) != 2 {
		t.Fatalf("preloaded users = %#v", withRelations)
	}
	joined, err := u.Joins(q.User.Company).Where(q.User.CompanyID.Eq(1)).Order(q.User.ID).Find()
	if err != nil || len(joined) != 2 || joined[0].Company.Name != "Acme" {
		t.Fatalf("joined users = %#v, err=%v", joined, err)
	}

	if err := q.Transaction(func(tx *query.Query) error {
		return tx.User.WithContext(context.Background()).Create(&fixture.User{Name: "committed", CompanyID: 1})
	}); err != nil {
		t.Fatalf("commit transaction: %v", err)
	}
	if count, err := u.Where(q.User.Name.Eq("committed")).Count(); err != nil || count != 1 {
		t.Fatalf("committed count = %d, err=%v", count, err)
	}

	rollbackErr := errors.New("rollback")
	err = q.Transaction(func(tx *query.Query) error {
		if err := tx.User.WithContext(context.Background()).Create(&fixture.User{Name: "rolled-back", CompanyID: 1}); err != nil {
			return err
		}
		return rollbackErr
	})
	if !errors.Is(err, rollbackErr) {
		t.Fatalf("rollback transaction error = %v", err)
	}
	if count, err := u.Where(q.User.Name.Eq("rolled-back")).Count(); err != nil || count != 0 {
		t.Fatalf("rolled-back count = %d, err=%v", count, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := q.User.WithContext(ctx).Find(); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled query error = %v, want context.Canceled", err)
	}
}

func assertUsers(t *testing.T, users []*fixture.User, err error, wantNames ...string) {
	t.Helper()
	if err != nil {
		t.Fatalf("query users: %v", err)
	}
	got := make([]string, len(users))
	for i, user := range users {
		got[i] = user.Name
	}
	sort.Strings(got)
	sort.Strings(wantNames)
	if len(got) != len(wantNames) {
		t.Fatalf("user names = %v, want %v", got, wantNames)
	}
	for index := range got {
		if got[index] != wantNames[index] {
			t.Fatalf("user names = %v, want %v", got, wantNames)
		}
	}
}
