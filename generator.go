package gen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
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

// NewGenerator create a new generator
func NewGenerator(cfg Config) *Generator {
	if cfg.ModelPkgName == "" {
		cfg.ModelPkgName = check.ModelPkg
	}
	if cfg.db == nil {
		cfg.db, _ = gorm.Open(tests.DummyDialector{})
	}
	return &Generator{
		Config: cfg,
		Data:   make(map[string]*genInfo),
	}
}

// Config generator's basic configuration
type Config struct {
	db *gorm.DB //nolint

	OutPath      string
	OutFile      string
	ModelPkgName string // generated model code's package name

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
	g.db = db
}

var (
	// FieldIgnore ignore some columns by name
	FieldIgnore = func(columnNames ...string) check.MemberOpt {
		return func(m *check.Member) *check.Member {
			for _, name := range columnNames {
				if m.Name == name {
					return nil
				}
			}
			return m
		}
	}
	// FieldRename specify field name in generated struct
	FieldRename = func(columnName string, newName string) check.MemberOpt {
		return func(m *check.Member) *check.Member {
			if m.Name == columnName {
				m.Name = newName
			}
			return m
		}
	}
	// FieldSetType specify field type in generated struct
	FieldSetType = func(columnName string, newType string) check.MemberOpt {
		return func(m *check.Member) *check.Member {
			if m.Name == columnName {
				m.Type = newType
				m.ModelType = newType
			}
			return m
		}
	}

	// FieldTag specify json tag and gorm tag
	FieldTag = func(columnName string, gormTag, jsonTag string) check.MemberOpt {
		return func(m *check.Member) *check.Member {
			if m.Name == columnName {
				m.GORMTag, m.JSONTag = gormTag, jsonTag
			}
			return m
		}
	}
	// FieldJSONTag specify json tag
	FieldJSONTag = func(columnName string, jsonTag string) check.MemberOpt {
		return func(m *check.Member) *check.Member {
			if m.Name == columnName {
				m.JSONTag = jsonTag
			}
			return m
		}
	}
	// FieldGORMTag specify gorm tag
	FieldGORMTag = func(columnName string, gormTag string) check.MemberOpt {
		return func(m *check.Member) *check.Member {
			if m.Name == columnName {
				m.GORMTag = gormTag
			}
			return m
		}
	}
	// FieldNewTag add new tag
	FieldNewTag = func(columnName string, newTag string) check.MemberOpt {
		return func(m *check.Member) *check.Member {
			if m.Name == columnName {
				m.NewTag += " " + newTag
			}
			return m
		}
	}
)

/*
** The feature of mapping table from database server to Golang struct
** Provided by @qqxhb
 */

// GenerateModel catch table info from db, return a BaseStruct
func (g *Generator) GenerateModel(tableName string, opts ...check.MemberOpt) *check.BaseStruct {
	return g.GenerateModelAs(tableName, g.db.Config.NamingStrategy.SchemaName(tableName), opts...)
}

// GenerateModel catch table info from db, return a BaseStruct
func (g *Generator) GenerateModelAs(tableName string, modelName string, opts ...check.MemberOpt) *check.BaseStruct {
	colNameOpts := make([]check.MemberOpt, len(opts))
	for i, opt := range opts {
		opt := opt
		colNameOpts[i] = opt
	}

	s, err := check.GenBaseStructs(g.db, g.Config.ModelPkgName, tableName, modelName, g.dbNameOpts, colNameOpts)
	if err != nil {
		g.db.Logger.Error(context.Background(), "generated struct from table has error: %s", err)
		panic("panic with generated struct error")
	}
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
		g.db.Logger.Error(context.Background(), "check struct error: %v", err)
		panic("panic with check struct error")
	}
	g.apply(fc, structs)
}

func (g *Generator) apply(fc interface{}, structs []*check.BaseStruct) {
	readInterface := new(parser.InterfaceSet)
	interfacePaths, err := parser.GetInterfacePath(fc)
	if err != nil {
		g.db.Logger.Error(context.Background(), "can not get interface name or file: %s", err)
		panic("panic with check interface error")
	}

	err = readInterface.ParseFile(interfacePaths, check.GetNames(structs))
	if err != nil {
		g.db.Logger.Error(context.Background(), "can not parser interface file: %s", err)
		panic("panic with parser interface file error")
	}

	for _, interfaceStruct := range structs {
		data, err := g.pushBaseStruct(interfaceStruct)
		if err != nil {
			g.db.Logger.Error(context.Background(), "gen struct error: %v", err)
			panic("panic with gen struct error")

		}

		functions, err := check.CheckInterface(readInterface, interfaceStruct, data.Interfaces)
		if err != nil {
			g.db.Logger.Error(context.Background(), "check interface error: %v", err)
			panic("panic with check interface error")
		}
		err = data.appendMethods(functions)
		if err != nil {
			g.db.Logger.Error(context.Background(), "check interface error: %v", err)
			panic("panic with check interface error")
		}
	}
}

// Execute generate code to output path
func (g *Generator) Execute() {
	var err error
	if g.OutPath == "" {
		g.OutPath = "./query"
	}
	if g.OutFile == "" {
		g.OutFile = g.OutPath + "/gorm_generated.go"
	}
	if _, err := os.Stat(g.OutPath); err != nil {
		if err := os.Mkdir(g.OutPath, os.ModePerm); err != nil {
			g.db.Logger.Error(context.Background(), "mkdir failed: %s", err)
			panic("panic with mkdir query dir error")
		}
	}
	g.queryPkgName = filepath.Base(g.OutPath)

	err = g.generatedBaseStruct()
	if err != nil {
		g.db.Logger.Error(context.Background(), "generate basic struct from table fail: %v", err)
		panic("panic with generate basic struct from table error")
	}
	err = g.generatedToOutFile()
	if err != nil {
		g.db.Logger.Error(context.Background(), "generate query to file: %v", err)
		panic("panic with generate query to file error")
	}

	g.successInfo(
		"Gorm generate query object file successful!",
		"Generated path："+g.OutPath,
		"Generated file："+g.OutFile,
	)
}

// successInfo logger
func (g *Generator) successInfo(logInfos ...string) {
	for _, l := range logInfos {
		g.db.Logger.Info(context.Background(), l)
		log.Println(l)
	}
}

// generatedToOutFile save generate code to file
func (g *Generator) generatedToOutFile() (err error) {
	var buf bytes.Buffer

	err = render(tmpl.HeaderTmpl, &buf, g.queryPkgName)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(g.Data))
	for key := range g.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		data := g.Data[key]
		err = render(tmpl.BaseStruct, &buf, data.BaseStruct)
		if err != nil {
			return err
		}

		for _, method := range data.Interfaces {
			err = render(tmpl.FuncTmpl, &buf, method)
			if err != nil {
				return err
			}
		}

		err = render(tmpl.BaseGormFunc, &buf, data.BaseStruct)
		if err != nil {
			return err
		}
	}

	err = render(tmpl.UseTmpl, &buf, g)
	if err != nil {
		return err
	}

	result, err := imports.Process(g.OutFile, buf.Bytes(), nil)
	if err != nil {
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		line := strings.Split(buf.String(), "\n")
		for i := -3; i < 3; i++ {
			fmt.Println(i+errLine, line[i+errLine])
		}
		return fmt.Errorf("can not format query file: %w", err)
	}
	return outputFile(g.OutFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, result)
}

// generatedBaseStruct generate basic structures
func (g *Generator) generatedBaseStruct() (err error) {
	var outPath string
	outPath, err = filepath.Abs(g.OutPath)
	if err != nil {
		return err
	}
	pkg := g.ModelPkgName
	if pkg == "" {
		pkg = check.ModelPkg
	}
	outPath = fmt.Sprint(filepath.Dir(outPath), "/", pkg, "/")

	_, err = os.Stat(outPath)
	created := err == nil

	for _, data := range g.Data {
		if data.BaseStruct == nil || !data.BaseStruct.GenBaseStruct {
			continue
		}

		if !created {
			if err := os.Mkdir(outPath, os.ModePerm); err != nil {
				g.db.Logger.Error(context.Background(), "mkdir failed: %s", err)
				panic("panic with mkdir base struct dir error")
			}
			created = true
		}

		var buf bytes.Buffer
		err = render(tmpl.ModelTemplate, &buf, data.BaseStruct)
		if err != nil {
			return err
		}
		modelFile := fmt.Sprint(outPath, data.BaseStruct.TableName, ".go")
		result, err := imports.Process(modelFile, buf.Bytes(), nil)
		if err != nil {
			for i, line := range strings.Split(buf.String(), "\n") {
				fmt.Println(i, line)
			}
			return fmt.Errorf("can not format struct file: %w", err)
		}
		err = outputFile(modelFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, result)
		if err != nil {
			return nil
		}
		g.successInfo(fmt.Sprintf("Generate struct [%s.%s] from table [%s]", data.StructInfo.Package, data.StructInfo.Type, data.TableName))
	}
	return nil
}

func (g *Generator) pushBaseStruct(base *check.BaseStruct) (*genInfo, error) {
	structName := base.StructName
	if g.Data[structName] == nil {
		g.Data[structName] = &genInfo{BaseStruct: base}
	}
	if g.Data[structName].Source != base.Source {
		return nil, fmt.Errorf("can not generate struct with the same name from different source:%s.%s and %s.%s",
			base.StructInfo.Package, base.StructName, g.Data[structName].StructInfo.Package, g.Data[structName].StructName)
	}
	return g.Data[structName], nil
}

func outputFile(filename string, flag int, data []byte) error {
	out, err := os.OpenFile(filename, flag, 0640)
	if err != nil {
		return fmt.Errorf("can not open out file: %w", err)
	}
	return output(out, data)
}

func output(wr io.WriteCloser, data []byte) (err error) {
	defer func() {
		if e := wr.Close(); e != nil {
			err = fmt.Errorf("can not close: %w", e)
		}
	}()

	if _, err = wr.Write(data); err != nil {
		return fmt.Errorf("can not write: %w", err)
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
