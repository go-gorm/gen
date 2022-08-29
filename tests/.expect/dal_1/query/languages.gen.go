// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"gorm.io/gen/tests/.gen/dal_1/model"
)

func newLanguage(db *gorm.DB) language {
	_language := language{}

	_language.languageDo.UseDB(db)
	_language.languageDo.UseModel(&model.Language{})

	tableName := _language.languageDo.TableName()
	_language.ALL = field.NewAsterisk(tableName)
	_language.ID = field.NewInt64(tableName, "id")
	_language.CreatedAt = field.NewTime(tableName, "created_at")
	_language.UpdatedAt = field.NewTime(tableName, "updated_at")
	_language.DeletedAt = field.NewField(tableName, "deleted_at")
	_language.Name = field.NewString(tableName, "name")

	_language.fillFieldMap()

	return _language
}

type language struct {
	languageDo languageDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	Name      field.String

	fieldMap map[string]field.Expr
}

func (l language) Table(newTableName string) *language {
	l.languageDo.UseTable(newTableName)
	return l.updateTableName(newTableName)
}

func (l language) As(alias string) *language {
	l.languageDo.DO = *(l.languageDo.As(alias).(*gen.DO))
	return l.updateTableName(alias)
}

func (l *language) updateTableName(table string) *language {
	l.ALL = field.NewAsterisk(table)
	l.ID = field.NewInt64(table, "id")
	l.CreatedAt = field.NewTime(table, "created_at")
	l.UpdatedAt = field.NewTime(table, "updated_at")
	l.DeletedAt = field.NewField(table, "deleted_at")
	l.Name = field.NewString(table, "name")

	l.fillFieldMap()

	return l
}

func (l *language) WithContext(ctx context.Context) *languageDo { return l.languageDo.WithContext(ctx) }

func (l language) TableName() string { return l.languageDo.TableName() }

func (l language) Alias() string { return l.languageDo.Alias() }

func (l *language) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := l.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (l *language) fillFieldMap() {
	l.fieldMap = make(map[string]field.Expr, 5)
	l.fieldMap["id"] = l.ID
	l.fieldMap["created_at"] = l.CreatedAt
	l.fieldMap["updated_at"] = l.UpdatedAt
	l.fieldMap["deleted_at"] = l.DeletedAt
	l.fieldMap["name"] = l.Name
}

func (l language) clone(db *gorm.DB) language {
	l.languageDo.ReplaceDB(db)
	return l
}

type languageDo struct{ gen.DO }

func (l languageDo) Debug() *languageDo {
	return l.withDO(l.DO.Debug())
}

func (l languageDo) WithContext(ctx context.Context) *languageDo {
	return l.withDO(l.DO.WithContext(ctx))
}

func (l languageDo) ReadDB() *languageDo {
	return l.Clauses(dbresolver.Read)
}

func (l languageDo) WriteDB() *languageDo {
	return l.Clauses(dbresolver.Write)
}

func (l languageDo) Clauses(conds ...clause.Expression) *languageDo {
	return l.withDO(l.DO.Clauses(conds...))
}

func (l languageDo) Returning(value interface{}, columns ...string) *languageDo {
	return l.withDO(l.DO.Returning(value, columns...))
}

func (l languageDo) Not(conds ...gen.Condition) *languageDo {
	return l.withDO(l.DO.Not(conds...))
}

func (l languageDo) Or(conds ...gen.Condition) *languageDo {
	return l.withDO(l.DO.Or(conds...))
}

func (l languageDo) Select(conds ...field.Expr) *languageDo {
	return l.withDO(l.DO.Select(conds...))
}

func (l languageDo) Where(conds ...gen.Condition) *languageDo {
	return l.withDO(l.DO.Where(conds...))
}

func (l languageDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *languageDo {
	return l.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (l languageDo) Order(conds ...field.Expr) *languageDo {
	return l.withDO(l.DO.Order(conds...))
}

func (l languageDo) Distinct(cols ...field.Expr) *languageDo {
	return l.withDO(l.DO.Distinct(cols...))
}

func (l languageDo) Omit(cols ...field.Expr) *languageDo {
	return l.withDO(l.DO.Omit(cols...))
}

func (l languageDo) Join(table schema.Tabler, on ...field.Expr) *languageDo {
	return l.withDO(l.DO.Join(table, on...))
}

func (l languageDo) LeftJoin(table schema.Tabler, on ...field.Expr) *languageDo {
	return l.withDO(l.DO.LeftJoin(table, on...))
}

func (l languageDo) RightJoin(table schema.Tabler, on ...field.Expr) *languageDo {
	return l.withDO(l.DO.RightJoin(table, on...))
}

func (l languageDo) Group(cols ...field.Expr) *languageDo {
	return l.withDO(l.DO.Group(cols...))
}

func (l languageDo) Having(conds ...gen.Condition) *languageDo {
	return l.withDO(l.DO.Having(conds...))
}

func (l languageDo) Limit(limit int) *languageDo {
	return l.withDO(l.DO.Limit(limit))
}

func (l languageDo) Offset(offset int) *languageDo {
	return l.withDO(l.DO.Offset(offset))
}

func (l languageDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *languageDo {
	return l.withDO(l.DO.Scopes(funcs...))
}

func (l languageDo) Unscoped() *languageDo {
	return l.withDO(l.DO.Unscoped())
}

func (l languageDo) Create(values ...*model.Language) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Create(values)
}

func (l languageDo) CreateInBatches(values []*model.Language, batchSize int) error {
	return l.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (l languageDo) Save(values ...*model.Language) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Save(values)
}

func (l languageDo) First() (*model.Language, error) {
	if result, err := l.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Language), nil
	}
}

func (l languageDo) Take() (*model.Language, error) {
	if result, err := l.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Language), nil
	}
}

func (l languageDo) Last() (*model.Language, error) {
	if result, err := l.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Language), nil
	}
}

func (l languageDo) Find() ([]*model.Language, error) {
	result, err := l.DO.Find()
	return result.([]*model.Language), err
}

func (l languageDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Language, err error) {
	buf := make([]*model.Language, 0, batchSize)
	err = l.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (l languageDo) FindInBatches(result *[]*model.Language, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return l.DO.FindInBatches(result, batchSize, fc)
}

func (l languageDo) Attrs(attrs ...field.AssignExpr) *languageDo {
	return l.withDO(l.DO.Attrs(attrs...))
}

func (l languageDo) Assign(attrs ...field.AssignExpr) *languageDo {
	return l.withDO(l.DO.Assign(attrs...))
}

func (l languageDo) Joins(fields ...field.RelationField) *languageDo {
	for _, _f := range fields {
		l = *l.withDO(l.DO.Joins(_f))
	}
	return &l
}

func (l languageDo) Preload(fields ...field.RelationField) *languageDo {
	for _, _f := range fields {
		l = *l.withDO(l.DO.Preload(_f))
	}
	return &l
}

func (l languageDo) FirstOrInit() (*model.Language, error) {
	if result, err := l.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Language), nil
	}
}

func (l languageDo) FirstOrCreate() (*model.Language, error) {
	if result, err := l.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Language), nil
	}
}

func (l languageDo) FindByPage(offset int, limit int) (result []*model.Language, count int64, err error) {
	result, err = l.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = l.Offset(-1).Limit(-1).Count()
	return
}

func (l languageDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = l.Count()
	if err != nil {
		return
	}

	err = l.Offset(offset).Limit(limit).Scan(result)
	return
}

func (l languageDo) Scan(result interface{}) (err error) {
	return l.DO.Scan(result)
}

func (l languageDo) Delete(models ...*model.Language) (result gen.ResultInfo, err error) {
	return l.DO.Delete(models)
}

func (l *languageDo) withDO(do gen.Dao) *languageDo {
	l.DO = *do.(*gen.DO)
	return l
}