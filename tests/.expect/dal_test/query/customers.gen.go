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

	"gorm.io/gen/tests/.expect/dal_test/model"
)

func newCustomer(db *gorm.DB, opts ...gen.DOOption) customer {
	_customer := customer{}

	_customer.customerDo.UseDB(db, opts...)
	_customer.customerDo.UseModel(&model.Customer{})

	tableName := _customer.customerDo.TableName()
	_customer.ALL = field.NewAsterisk(tableName)
	_customer.ID = field.NewInt64(tableName, "id")
	_customer.CreatedAt = field.NewTime(tableName, "created_at")
	_customer.UpdatedAt = field.NewTime(tableName, "updated_at")
	_customer.DeletedAt = field.NewField(tableName, "deleted_at")
	_customer.BankID = field.NewInt64(tableName, "bank_id")

	_customer.fillFieldMap()

	return _customer
}

type customer struct {
	customerDo customerDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	BankID    field.Int64

	fieldMap map[string]field.Expr
}

func (c customer) Table(newTableName string) *customer {
	c.customerDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c customer) As(alias string) *customer {
	c.customerDo.DO = *(c.customerDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *customer) updateTableName(table string) *customer {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewInt64(table, "id")
	c.CreatedAt = field.NewTime(table, "created_at")
	c.UpdatedAt = field.NewTime(table, "updated_at")
	c.DeletedAt = field.NewField(table, "deleted_at")
	c.BankID = field.NewInt64(table, "bank_id")

	c.fillFieldMap()

	return c
}

func (c *customer) WithContext(ctx context.Context) *customerDo { return c.customerDo.WithContext(ctx) }

func (c customer) TableName() string { return c.customerDo.TableName() }

func (c customer) Alias() string { return c.customerDo.Alias() }

func (c customer) Columns(cols ...field.Expr) gen.Columns { return c.customerDo.Columns(cols...) }

func (c *customer) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *customer) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 5)
	c.fieldMap["id"] = c.ID
	c.fieldMap["created_at"] = c.CreatedAt
	c.fieldMap["updated_at"] = c.UpdatedAt
	c.fieldMap["deleted_at"] = c.DeletedAt
	c.fieldMap["bank_id"] = c.BankID
}

func (c customer) clone(db *gorm.DB) customer {
	c.customerDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c customer) replaceDB(db *gorm.DB) customer {
	c.customerDo.ReplaceDB(db)
	return c
}

type customerDo struct{ gen.DO }

func (c customerDo) Debug() *customerDo {
	return c.withDO(c.DO.Debug())
}

func (c customerDo) WithContext(ctx context.Context) *customerDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c customerDo) ReadDB() *customerDo {
	return c.Clauses(dbresolver.Read)
}

func (c customerDo) WriteDB() *customerDo {
	return c.Clauses(dbresolver.Write)
}

func (c customerDo) Session(config *gorm.Session) *customerDo {
	return c.withDO(c.DO.Session(config))
}

func (c customerDo) Clauses(conds ...clause.Expression) *customerDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c customerDo) Returning(value interface{}, columns ...string) *customerDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c customerDo) Not(conds ...gen.Condition) *customerDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c customerDo) Or(conds ...gen.Condition) *customerDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c customerDo) Select(conds ...field.Expr) *customerDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c customerDo) Where(conds ...gen.Condition) *customerDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c customerDo) Order(conds ...field.Expr) *customerDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c customerDo) Distinct(cols ...field.Expr) *customerDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c customerDo) Omit(cols ...field.Expr) *customerDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c customerDo) Join(table schema.Tabler, on ...field.Expr) *customerDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c customerDo) LeftJoin(table schema.Tabler, on ...field.Expr) *customerDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c customerDo) RightJoin(table schema.Tabler, on ...field.Expr) *customerDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c customerDo) Group(cols ...field.Expr) *customerDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c customerDo) Having(conds ...gen.Condition) *customerDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c customerDo) Limit(limit int) *customerDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c customerDo) Offset(offset int) *customerDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c customerDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *customerDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c customerDo) Unscoped() *customerDo {
	return c.withDO(c.DO.Unscoped())
}

func (c customerDo) Create(values ...*model.Customer) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c customerDo) CreateInBatches(values []*model.Customer, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c customerDo) Save(values ...*model.Customer) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c customerDo) First() (*model.Customer, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) Take() (*model.Customer, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) Last() (*model.Customer, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) Find() ([]*model.Customer, error) {
	result, err := c.DO.Find()
	return result.([]*model.Customer), err
}

func (c customerDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Customer, err error) {
	buf := make([]*model.Customer, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c customerDo) FindInBatches(result *[]*model.Customer, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c customerDo) Attrs(attrs ...field.AssignExpr) *customerDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c customerDo) Assign(attrs ...field.AssignExpr) *customerDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c customerDo) Joins(fields ...field.RelationField) *customerDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c customerDo) Preload(fields ...field.RelationField) *customerDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c customerDo) FirstOrInit() (*model.Customer, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) FirstOrCreate() (*model.Customer, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Customer), nil
	}
}

func (c customerDo) FindByPage(offset int, limit int) (result []*model.Customer, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c customerDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c customerDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c customerDo) Delete(models ...*model.Customer) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *customerDo) withDO(do gen.Dao) *customerDo {
	c.DO = *do.(*gen.DO)
	return c
}
