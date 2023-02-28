package gen

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/os/gfile"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	g             *Generator
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

func (r *Conf) Connect() (*gorm.DB, error) {
	if r.DSN == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}
	switch r.DB {
	case DbMySQL:
		return gorm.Open(mysql.Open(r.DSN), &gorm.Config{Logger: newLogger})
	case DbPostgres:
		return gorm.Open(postgres.Open(r.DSN), &gorm.Config{Logger: newLogger})
	case DbSQLite:
		return gorm.Open(sqlite.Open(r.DSN), &gorm.Config{Logger: newLogger})
	case DbSQLServer:
		return gorm.Open(sqlserver.Open(r.DSN), &gorm.Config{Logger: newLogger})
	default:
		return nil, fmt.Errorf("unknow db %q (support mysql || postgres || sqlite || sqlserver for now)", r.DB)
	}
}
func (r *Conf) Delete() {
	if r.Reset {
		var (
			err    error
			db     = r.g.db
			tables = make([]string, 0)
		)
		if tables, err = db.Migrator().GetTables(); err == nil {
			if err = db.Migrator().DropTable(tables); err != nil {
				log.Println("删除失败", err)
			}
		}
	}
}

func (r *Conf) GenModels() (models []interface{}, err error) {
	var (
		g = NewGenerator(Config{
			Mode:              WithDefaultQuery | WithoutContext,
			OutPath:           r.OutPath,
			OutFile:           r.OutFile,
			ModelPkgPath:      r.ModelPkgName,
			WithUnitTest:      r.WithUnitTest,
			FieldNullable:     r.FieldNullable,
			FieldWithIndexTag: true,
			FieldWithTypeTag:  true,
			FieldSignable:     r.FieldSignable,
		})
		tablesList []string
	)
	if g.db, err = r.Connect(); err != nil {
		return nil, fmt.Errorf("connect to database fail: %w", err)
	}
	r.Delete()
	if tablesList, err = r.g.db.Migrator().GetTables(); err != nil {
		return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
	}
	models = make([]interface{}, len(tablesList))
	for i, tableName := range tablesList {
		if opt := g.GetModel(tableName); opt != nil {
			models[i] = g.GenerateModel(tableName, WithMethod(opt))
		} else {
			models[i] = g.GenerateModel(tableName)
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
}
