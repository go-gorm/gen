package gen

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	logger.Config{
		SlowThreshold:             time.Second,   // 慢 SQL 阈值
		LogLevel:                  logger.Silent, // 日志级别
		IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  true,          // 禁用彩色打印
	},
)

const (
	DbMySQL     string = "mysql"
	DbPostgres  string = "postgres"
	DbSQLite    string = "sqlite"
	DbSQLServer string = "sqlserver"
)

var (
	App = &Conf{
		DSN:           "host=192.168.32.130 user=postgres password=postgres dbname=status port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		DB:            DbPostgres,
		OutPath:       "./app/dao",
		WithUnitTest:  true,
		ModelPkgName:  "",
		FieldNullable: false,
		//FieldWithIndexTag: false,
		//FieldWithTypeTag:  false,
		FieldSignable: false,
	}
)

const (
	DefaultOutPath = "./dao/query"
)

type Conf struct {
	*Generator    `yaml:"-"`
	Reset         bool     `yaml:"reset"`         //是否重置数据库
	DSN           string   `yaml:"dsn"`           // consult[https://gorm.io/docs/connecting_to_the_database.html]"
	DB            string   `yaml:"db"`            // input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
	Tables        []string `yaml:"tables"`        // enter the required data table or leave it blank
	OnlyModel     bool     `yaml:"onlyModel"`     // only generate model
	OutPath       string   `yaml:"outPath"`       // specify a directory for output
	OutFile       string   `yaml:"outFile"`       // query code file name, default: gen.go
	WithUnitTest  bool     `yaml:"withUnitTest"`  // generate unit test for query code
	ModelPkgName  string   `yaml:"modelPkgName"`  // generated model code's package name
	FieldNullable bool     `yaml:"fieldNullable"` // generate with pointer when field is nullable
	//FieldWithIndexTag bool     `yaml:"fieldWithIndexTag"` // generate field with gorm index tag
	//FieldWithTypeTag  bool     `yaml:"fieldWithTypeTag"`  // generate field with gorm column type tag
	FieldSignable bool `yaml:"fieldSignable"` // detect integer field's unsigned type, adjust generated data type
}

func (r *Conf) Connect() (err error) {
	if r.DSN == "" {
		return fmt.Errorf("dsn cannot be empty")
	}
	switch r.DB {
	case DbMySQL:
		if r.db, err = gorm.Open(mysql.Open(r.DSN), &gorm.Config{Logger: newLogger}); err != nil {
			return
		}
		return nil
	case DbPostgres:
		if r.db, err = gorm.Open(postgres.Open(r.DSN), &gorm.Config{Logger: newLogger}); err != nil {
			return
		}
		return nil
	case DbSQLite:
		if r.db, err = gorm.Open(sqlite.Open(r.DSN), &gorm.Config{Logger: newLogger}); err != nil {
			return
		}
		return nil
	case DbSQLServer:
		if r.db, err = gorm.Open(sqlserver.Open(r.DSN), &gorm.Config{Logger: newLogger}); err != nil {
			return
		}
		return nil
	default:
		return fmt.Errorf("unknow db %q (support mysql || postgres || sqlite || sqlserver for now)", r.DB)
	}
}
func (r *Conf) Delete() {
	if r.Reset {
		var (
			err    error
			db     = r.db
			tables = make([]string, 0)
		)
		if tables, err = db.Migrator().GetTables(); err == nil {
			if err = db.Migrator().DropTable(gconv.SliceAny(tables)...); err != nil {
				log.Println("删除失败", err)
			}
		}
	}
}

func (r *Conf) GenModels() (models []interface{}, err error) {
	var (
		tablesList []string
	)
	if tablesList, err = r.db.Migrator().GetTables(); err != nil {
		return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
	}
	models = make([]interface{}, len(tablesList))
	for i, tableName := range tablesList {
		if opt := r.GetModel(tableName); opt != nil {
			models[i] = r.GenerateModel(tableName, WithMethod(opt))
		} else {
			models[i] = r.GenerateModel(tableName)
		}
	}
	return models, nil
}

func Parse() {
	if !gfile.Exists("database.yml") {
		if encode, err := gyaml.Encode(App); err != nil {
			log.Fatalf("encode fail: %v", err)
		} else {
			if err = gfile.PutBytes("database.yml", encode); err != nil {
				log.Fatalf("write fail: %v", err)
			}
		}
		log.Fatalln("create database.yml success")
	} else {
		if decode := gfile.GetBytes("database.yml"); decode == nil {
			log.Fatalf("read fail")
		} else {
			if err := gyaml.DecodeTo(decode, App); err != nil {
				log.Fatalf("decode fail: %v", err)
			}
		}
	}
	App.Generator = NewGenerator(Config{
		Mode:              WithDefaultQuery | WithoutContext,
		OutPath:           App.OutPath,
		OutFile:           App.OutFile,
		ModelPkgPath:      App.ModelPkgName,
		WithUnitTest:      App.WithUnitTest,
		FieldNullable:     App.FieldNullable,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		FieldSignable:     App.FieldSignable,
	})
	App.Generator.Schema = Schema{
		Schema:    make(map[string]*schema.Schema),
		Model:     make(map[string]any),
		Generator: App.Generator,
	}
}
