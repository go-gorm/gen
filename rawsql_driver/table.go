package rawsql_driver

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/parser/test_driver"
	"github.com/pingcap/tidb/parser/types"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

type Table struct {
	ColumnTypes []gorm.ColumnType
	Indexes     []gorm.Index
	Name        string
}

type Parser interface {
	Tables(createSql string) (tables []*Table, err error)
}

type defaultParser struct{}

func (d *defaultParser) Tables(createSql string) (tables []*Table, err error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(createSql, "", "")
	if err != nil {
		return nil, err
	}
	tables = make([]*Table, 0, len(stmtNodes))
	for _, node := range stmtNodes {
		//convert fail will panic
		create := node.(*ast.CreateTableStmt)
		tables = append(tables, &Table{
			Name:        create.Table.Name.String(),
			ColumnTypes: d.getColumnTypes(create),
			Indexes:     d.getIndexes(create),
		})

	}
	return tables, nil
}

func (*defaultParser) getColumnTypes(create *ast.CreateTableStmt) (cols []gorm.ColumnType) {
	if create == nil || len(create.Cols) == 0 {
		return nil
	}
	cols = make([]gorm.ColumnType, 0, len(create.Cols))
	for _, col := range create.Cols {
		ct := &migrator.ColumnType{
			NameValue:        sql.NullString{Valid: true, String: col.Name.String()},
			DataTypeValue:    sql.NullString{Valid: true, String: types.TypeToStr(col.Tp.GetType(), col.Tp.GetCharset())},
			ColumnTypeValue:  sql.NullString{Valid: true, String: col.Tp.String()},
			PrimaryKeyValue:  sql.NullBool{Bool: mysql.HasPriKeyFlag(col.Tp.GetFlag()), Valid: mysql.HasPriKeyFlag(col.Tp.GetFlag())},
			UniqueValue:      sql.NullBool{Bool: mysql.HasUniKeyFlag(col.Tp.GetFlag()), Valid: mysql.HasUniKeyFlag(col.Tp.GetFlag())},
			LengthValue:      sql.NullInt64{Int64: int64(col.Tp.GetFlen()), Valid: col.Tp.IsVarLengthType()},
			DecimalSizeValue: sql.NullInt64{Int64: int64(col.Tp.GetFlen()), Valid: col.Tp.IsDecimalValid()},
			ScaleValue:       sql.NullInt64{Int64: int64(col.Tp.GetDecimal()), Valid: col.Tp.IsDecimalValid()},
			NullableValue:    sql.NullBool{Bool: true, Valid: true},
			SQLColumnType:    &sql.ColumnType{},
			ScanTypeValue:    getType(col.Tp),
		}
		for _, opt := range col.Options {
			if opt.Tp == ast.ColumnOptionNotNull {
				ct.NullableValue.Bool = false
				continue
			}
			if opt.Tp == ast.ColumnOptionComment {
				ct.CommentValue = sql.NullString{String: opt.Expr.(*test_driver.ValueExpr).Datum.GetString(), Valid: true}
				continue
			}
			if opt.Tp == ast.ColumnOptionAutoIncrement {
				ct.AutoIncrementValue = sql.NullBool{Bool: true, Valid: true}
				continue
			}
			if opt.Tp == ast.ColumnOptionDefaultValue {
				if v, ok := opt.Expr.(*test_driver.ValueExpr); ok {
					ct.CommentValue = sql.NullString{Valid: true, String: v.Datum.GetString()}
					continue
				}

				if v2, ok := opt.Expr.(*ast.FuncCallExpr); ok {
					ct.CommentValue = sql.NullString{Valid: true, String: v2.FnName.String()}
				}
			}
		}
		cols = append(cols, ct)
	}
	return cols
}

func (d *defaultParser) getIndexes(create *ast.CreateTableStmt) []gorm.Index {
	if create == nil || len(create.Constraints) == 0 {
		return nil
	}
	indexs := make([]gorm.Index, 0, len(create.Constraints))
	table := create.Table.Name.String()
	for _, cons := range create.Constraints {
		idx := &migrator.Index{TableName: table, NameValue: cons.Name, ColumnList: []string{},
			PrimaryKeyValue: sql.NullBool{Bool: ast.ConstraintPrimaryKey == cons.Tp, Valid: ast.ConstraintPrimaryKey == cons.Tp},
			UniqueValue:     sql.NullBool{Bool: ast.ConstraintUniq == cons.Tp, Valid: ast.ConstraintUniq == cons.Tp},
		}
		for _, col := range cons.Keys {
			idx.ColumnList = append(idx.ColumnList, col.Column.Name.String())
		}
		indexs = append(indexs, idx)
	}
	return indexs
}

var (
	intT    = reflect.TypeOf(int32(0))
	longT   = reflect.TypeOf(int64(0))
	boolT   = reflect.TypeOf(false)
	stringT = reflect.TypeOf("")
	floatT  = reflect.TypeOf(float32(0))
	doubleT = reflect.TypeOf(float64(0))
	timeT   = reflect.TypeOf(time.Time{})
)

func getType(tp *types.FieldType) reflect.Type {
	if tp == nil {
		return nil
	}
	switch tp.GetType() {
	case mysql.TypeTiny, mysql.TypeShort, mysql.TypeLong:
		return intT
	case mysql.TypeFloat:
		return floatT
	case mysql.TypeDouble:
		return doubleT
	case mysql.TypeTimestamp, mysql.TypeLonglong, mysql.TypeInt24:
		return longT
	case mysql.TypeDate, mysql.TypeDatetime, mysql.TypeNewDate:
		return timeT
	default:
		return stringT
	}
}
