package rawsql_driver

import (
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/pingcap/tidb/parser/test_driver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Config struct {
	DriverName string   //mysql
	FilePath   []string //create table sql file or file path
	SQL        []string //create table sql content
	Parser
}

type Dialector struct {
	*Config
	tables map[string]*Table
}

func New(config Config) gorm.Dialector {
	return &Dialector{Config: &config, tables: map[string]*Table{}}
}

func (dialector Dialector) Name() string {
	return dialector.DriverName
}

func (dialector Dialector) Initialize(db *gorm.DB) error {
	if dialector.DriverName == "" {
		dialector.DriverName = "mysql"
	}
	if dialector.SQL == nil {
		dialector.SQL = make([]string, 0)
	}
	if dialector.tables == nil {
		dialector.tables = make(map[string]*Table)
	}
	if dialector.Parser == nil {
		dialector.Parser = &defaultParser{}
	}
	if err := dialector.fileTOSQL(); err != nil {
		return err
	}

	if err := dialector.sqlTOTable(); err != nil {
		return err
	}

	return nil
}

func (dialector Dialector) sqlTOTable() error {
	for _, sql := range dialector.SQL {
		tables, err := dialector.Parser.Tables(sql)
		if err != nil {
			return err
		}
		for _, table := range tables {
			dialector.tables[table.Name] = table
		}
	}
	return nil
}

func (dialector Dialector) fileTOSQL() error {
	for _, f := range dialector.FilePath {
		if f == "" {
			continue
		}
		v, err := os.Stat(f)
		if err != nil {
			return err
		}
		if v.IsDir() {
			err = dialector.readFiles(f)
		} else {
			err = dialector.readFile(f)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (dialector Dialector) readFiles(folder string) (err error) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		fn := filepath.Join(folder, file.Name())
		if file.IsDir() {
			err = dialector.readFiles(fn)
		} else {
			err = dialector.readFile(fn)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (dialector Dialector) readFile(fileName string) error {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	dialector.SQL = append(dialector.SQL, string(content))
	return nil
}

func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{
		Migrator: migrator.Migrator{
			Config: migrator.Config{
				DB:        db,
				Dialector: dialector,
			},
		},
		Dialector: dialector,
	}
}

func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	return ""
}
func (dialector Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "DEFAULT"}
}
func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {}
func (dialector Dialector) QuoteTo(writer clause.Writer, str string)                            {}
func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	return ""
}
