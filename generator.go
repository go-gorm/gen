package gen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/internal/check"
	"gorm.io/gen/internal/parser"
	tmpl "gorm.io/gen/internal/template"
)

// TODO implement some unit tests

// T generic type
type T interface{}

// M map[string]interface{}
type M map[string]interface{}

// RowsAffected execute affected raws
type RowsAffected int64

// NewGenerator create a new generator
func NewGenerator(cfg Config) *Generator {
	err := cfg.Revise()
	if err != nil {
		panic(fmt.Errorf("create generator fail: %w", err))
	}

	return &Generator{
		Config: cfg,
		Data:   make(map[string]*genInfo),
	}
}

type GenerateMode uint

const (
	// WithDefaultQuery create default query in generated code
	WithDefaultQuery GenerateMode = 1 << iota

	// WithoutContext generate code without context constrain
	WithoutContext
)

// Config generator's basic configuration
type Config struct {
	db *gorm.DB //nolint

	OutPath       string
	OutFile       string
	ModelPkgPath  string // generated model code's package name
	FieldNullable bool

	Mode GenerateMode // generate mode

	queryPkgName string // generated query code's package name
	dbNameOpts   []check.SchemaNameOpt
}

// WithDbNameOpts set get database name function
func (cfg *Config) WithDbNameOpts(opts ...check.SchemaNameOpt) {
	if cfg.dbNameOpts == nil {
		cfg.dbNameOpts = opts
	} else {
		cfg.dbNameOpts = append(cfg.dbNameOpts, opts...)
	}
}

func (cfg *Config) Revise() (err error) {
	if cfg.ModelPkgPath == "" {
		cfg.ModelPkgPath = check.DefaultModelPkg
	}

	cfg.OutPath, err = filepath.Abs(cfg.OutPath)
	if err != nil {
		return fmt.Errorf("outpath is invalid: %w", err)
	}

	if cfg.db == nil {
		cfg.db, _ = gorm.Open(tests.DummyDialector{})
	}

	return nil
}

func (cfg *Config) judgeMode(mode GenerateMode) bool { return cfg.Mode&mode != 0 }

// genInfo info about generated code
type genInfo struct {
	*check.BaseStruct
	Interfaces []*check.InterfaceMethod
}

//
func (i *genInfo) appendMethods(methods []*check.InterfaceMethod) error {
	for _, newMethod := range methods {
		if i.methodInGenInfo(newMethod) {
			continue
		}
		i.Interfaces = append(i.Interfaces, newMethod)
	}
	return nil
}

func (i *genInfo) methodInGenInfo(m *check.InterfaceMethod) bool {
	for _, method := range i.Interfaces {
		if method.IsRepeatFromSameInterface(m) {
			return true
		}
	}
	return false
}

// Generator code generator
type Generator struct {
	Config

	Data map[string]*genInfo
}

// UseDB set db connection
func (g *Generator) UseDB(db *gorm.DB) {
	if db != nil {
		g.db = db
	}
}

/*
** The feature of mapping table from database server to Golang struct
** Provided by @qqxhb
 */

// GenerateModel catch table info from db, return a BaseStruct
func (g *Generator) GenerateModel(tableName string, opts ...check.MemberOpt) *check.BaseStruct {
	return g.GenerateModelAs(tableName, g.db.Config.NamingStrategy.SchemaName(tableName), opts...)
}

// GenerateModel catch table info from db, return a BaseStruct
func (g *Generator) GenerateModelAs(tableName string, modelName string, fieldOpts ...check.MemberOpt) *check.BaseStruct {
	s, err := check.GenBaseStructs(g.db, g.Config.ModelPkgPath, tableName, modelName, g.dbNameOpts, fieldOpts, g.Config.FieldNullable)
	if err != nil {
		g.db.Logger.Error(context.Background(), "generate struct from table fail: %s", err)
		panic("generate struct fail")
	}

	g.successInfo(fmt.Sprintf("got %d columns from table <%s>", len(s.Members), s.TableName))
	return s
}

// ApplyBasic specify models which will implement basic method
func (g *Generator) ApplyBasic(models ...interface{}) {
	g.ApplyInterface(func() {}, models...)
}

// ApplyInterface specifies method interfaces on structures, implment codes will be generated after calling g.Execute()
// eg: g.ApplyInterface(func(model.Method){}, model.User{}, model.Company{})
func (g *Generator) ApplyInterface(fc interface{}, models ...interface{}) {
	structs, err := check.CheckStructs(g.db, models...)
	if err != nil {
		g.db.Logger.Error(context.Background(), "check struct fail: %v", err)
		panic("check struct fail")
	}
	g.apply(fc, structs)
}

func (g *Generator) apply(fc interface{}, structs []*check.BaseStruct) {
	readInterface := new(parser.InterfaceSet)
	interfacePaths, err := parser.GetInterfacePath(fc)
	if err != nil {
		g.db.Logger.Error(context.Background(), "get interface name or file fail: %s", err)
		panic("check interface fail")
	}

	err = readInterface.ParseFile(interfacePaths, check.GetStructNames(structs))
	if err != nil {
		g.db.Logger.Error(context.Background(), "parser interface file fail: %s", err)
		panic("parser interface file fail")
	}

	for _, interfaceStruct := range structs {
		if g.judgeMode(WithoutContext) {
			interfaceStruct.ReviseMemberName()
		}

		data, err := g.pushBaseStruct(interfaceStruct)
		if err != nil {
			g.db.Logger.Error(context.Background(), "gen struct fail: %v", err)
			panic("gen struct fail")
		}

		functions, err := check.CheckInterface(readInterface, interfaceStruct, data.Interfaces)
		if err != nil {
			g.db.Logger.Error(context.Background(), "check interface fail: %v", err)
			panic("check interface fail")
		}
		err = data.appendMethods(functions)
		if err != nil {
			g.db.Logger.Error(context.Background(), "check interface fail: %v", err)
			panic("check interface fail")
		}
	}
}

// Execute generate code to output path
func (g *Generator) Execute() {
	var err error

	if g.OutPath == "" {
		g.OutPath = "./query/"
	}
	if g.OutFile == "" {
		g.OutFile = g.OutPath + "/gen.go"
	}
	if err := os.MkdirAll(g.OutPath, os.ModePerm); err != nil {
		g.db.Logger.Error(context.Background(), "create outpath(%s) fail: %s", g.OutPath, err)
		panic("create outpath fail")
	}
	g.queryPkgName = filepath.Base(g.OutPath)

	err = g.generateBaseStruct()
	if err != nil {
		g.db.Logger.Error(context.Background(), "generate basic struct from table fail: %s", err)
		panic("generate basic struct from table fail")
	}
	g.deleteHistoryGeneratedFile()
	err = g.generateQueryFile()
	if err != nil {
		g.db.Logger.Error(context.Background(), "generate query code: %s", err)
		panic("generate query code fail")
	}

	g.successInfo(
		"Successfully generate query file："+g.OutFile,
		"Successfully generate code",
	)
}

// successInfo logger
func (g *Generator) successInfo(logInfos ...string) {
	for _, l := range logInfos {
		g.db.Logger.Info(context.Background(), l)
		log.Println(l)
	}
}

// generateQueryFile generate query code and save to file
func (g *Generator) generateQueryFile() (err error) {
	var buf bytes.Buffer

	err = render(tmpl.HeaderTmpl, &buf, g.queryPkgName)
	if err != nil {
		return err
	}

	if g.judgeMode(WithDefaultQuery) {
		err = render(tmpl.DefaultQueryTmpl, &buf, g)
		if err != nil {
			return err
		}
	}

	err = render(tmpl.QueryTmpl, &buf, g)
	if err != nil {
		return err
	}

	err = g.output(g.OutFile, buf.Bytes())
	if err != nil {
		return err
	}

	for _, info := range g.Data {
		err = g.generateSubQuery(info)
		if err != nil {
			return err
		}
	}
	return nil
}

// generateSubQuery generate query code and save to file
func (g *Generator) generateSubQuery(data *genInfo) (err error) {
	var buf bytes.Buffer

	err = render(tmpl.HeaderTmpl, &buf, g.queryPkgName)
	if err != nil {
		return err
	}

	structTmpl := tmpl.BaseStructWithContext
	if g.judgeMode(WithoutContext) {
		structTmpl = tmpl.BaseStruct
	}

	err = render(structTmpl, &buf, data.BaseStruct)
	if err != nil {
		return err
	}

	for _, method := range data.Interfaces {
		err = render(tmpl.DIYMethod, &buf, method)
		if err != nil {
			return err
		}
	}

	err = render(tmpl.CRUDMethod, &buf, data.BaseStruct)
	if err != nil {
		return err
	}
	return g.output(fmt.Sprintf("%s/%s.gen.go", g.OutPath, strings.ToLower(data.TableName)), buf.Bytes())
}

// remove history GEN generated file
func (g *Generator) deleteHistoryGeneratedFile() {
	historyFile := g.OutPath + "/gorm_generated.go"
	if _, err := os.Stat(g.OutPath); err == nil {
		_ = os.Remove(historyFile)
	}
}

// generateBaseStruct generate basic structures and save to file
func (g *Generator) generateBaseStruct() (err error) {
	var outPath string
	outPath, err = filepath.Abs(g.OutPath)
	if err != nil {
		return err
	}
	path := filepath.Clean(g.ModelPkgPath)
	if path == "" {
		path = check.DefaultModelPkg
	}
	if strings.Contains(path, "/") {
		outPath, err = filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("cannot parse model pkg path: %w", err)
		}
		outPath += "/"
	} else {
		outPath = fmt.Sprint(filepath.Dir(outPath), "/", path, "/")
	}

	created := false
	mkdir := func() {
		if created {
			return
		}
		if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
			g.db.Logger.Error(context.Background(), "create model pkg path(%s) fail: %s", outPath, err)
			panic("create model pkg path fail")
		}
	}

	for _, data := range g.Data {
		if data.BaseStruct == nil || !data.BaseStruct.GenBaseStruct {
			continue
		}

		mkdir()

		var buf bytes.Buffer
		err = render(tmpl.ModelTemplate, &buf, data.BaseStruct)
		if err != nil {
			return err
		}
		modelFile := fmt.Sprint(outPath, data.BaseStruct.TableName, ".gen.go")
		err = g.output(modelFile, buf.Bytes())
		if err != nil {
			return err
		}

		g.successInfo(fmt.Sprintf("Generate struct [%s.%s] from table [%s]", data.StructInfo.Package, data.StructInfo.Type, data.TableName))
		g.successInfo(fmt.Sprintf("Successfully generate struct file: %s", modelFile))
	}
	return nil
}

// output format and output
func (g *Generator) output(fileName string, content []byte) error {
	result, err := imports.Process(fileName, content, nil)
	if err != nil {
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		startLine, endLine := 0, errLine+3
		fmt.Println("Format fail:")
		line := strings.Split(string(content), "\n")
		for i := startLine; i <= endLine; i++ {
			fmt.Println(i+errLine, line[i+errLine])
		}
		return fmt.Errorf("cannot format struct file: %w", err)
	}
	return outputFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, result)
}

func (g *Generator) pushBaseStruct(base *check.BaseStruct) (*genInfo, error) {
	structName := base.StructName
	if g.Data[structName] == nil {
		g.Data[structName] = &genInfo{BaseStruct: base}
	}
	if g.Data[structName].Source != base.Source {
		return nil, fmt.Errorf("cannot generate struct with the same name from different source:%s.%s and %s.%s",
			base.StructInfo.Package, base.StructName, g.Data[structName].StructInfo.Package, g.Data[structName].StructName)
	}
	return g.Data[structName], nil
}

func outputFile(filename string, flag int, data []byte) error {
	out, err := os.OpenFile(filename, flag, 0640)
	if err != nil {
		return fmt.Errorf("open out file fail: %w", err)
	}
	return output(out, data)
}

func output(wr io.WriteCloser, data []byte) (err error) {
	defer func() {
		if e := wr.Close(); e != nil {
			err = fmt.Errorf("close file fail: %w", e)
		}
	}()

	if _, err = wr.Write(data); err != nil {
		return fmt.Errorf("write file fail: %w", err)
	}
	return nil
}

func render(tmpl string, wr io.Writer, data interface{}) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}
