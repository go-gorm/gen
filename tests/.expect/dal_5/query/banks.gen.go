// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/gen/tests/.gen/dal_5/model"
)

func newBank(db *gorm.DB, opts ...gen.DOOption) bank {
	_bank := bank{}

	_bank.bankDo.UseDB(db, opts...)
	_bank.bankDo.UseModel(&model.Bank{})

	tableName := _bank.bankDo.TableName()
	_bank.ALL = field.NewAsterisk(tableName)
	_bank.ID = field.NewInt64(tableName, "id")
	_bank.Name = field.NewString(tableName, "name")
	_bank.Address = field.NewString(tableName, "address")
	_bank.Scale = field.NewInt64(tableName, "scale")

	_bank.fillFieldMap()

	return _bank
}

type bank struct {
	bankDo bankDo

	ALL     field.Asterisk
	ID      field.Int64
	Name    field.String
	Address field.String
	Scale   field.Int64

	fieldMap map[string]field.Expr
}

func (b bank) Table(newTableName string) *bank {
	b.bankDo.UseTable(newTableName)
	return b.updateTableName(newTableName)
}

func (b bank) As(alias string) *bank {
	b.bankDo.DO = *(b.bankDo.As(alias).(*gen.DO))
	return b.updateTableName(alias)
}

func (b *bank) updateTableName(table string) *bank {
	b.ALL = field.NewAsterisk(table)
	b.ID = field.NewInt64(table, "id")
	b.Name = field.NewString(table, "name")
	b.Address = field.NewString(table, "address")
	b.Scale = field.NewInt64(table, "scale")

	b.fillFieldMap()

	return b
}

func (b *bank) WithContext(ctx context.Context) IBankDo { return b.bankDo.WithContext(ctx) }

func (b bank) TableName() string { return b.bankDo.TableName() }

func (b bank) Alias() string { return b.bankDo.Alias() }

func (b *bank) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := b.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (b *bank) fillFieldMap() {
	b.fieldMap = make(map[string]field.Expr, 4)
	b.fieldMap["id"] = b.ID
	b.fieldMap["name"] = b.Name
	b.fieldMap["address"] = b.Address
	b.fieldMap["scale"] = b.Scale
}

func (b bank) clone(db *gorm.DB) bank {
	b.bankDo.ReplaceConnPool(db.Statement.ConnPool)
	return b
}

func (b bank) replaceDB(db *gorm.DB) bank {
	b.bankDo.ReplaceDB(db)
	return b
}

type bankDo struct{ gen.DO }

type IBankDo interface {
	WithContext(ctx context.Context) IBankDo
}

func (b bankDo) WithContext(ctx context.Context) IBankDo {
	return b.withDO(b.DO.WithContext(ctx))
}

func (b *bankDo) withDO(do gen.Dao) *bankDo {
	b.DO = *do.(*gen.DO)
	return b
}
