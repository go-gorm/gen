package gen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gen/helper"
	"gorm.io/gen/internal/check"
	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
	tmpl "gorm.io/gen/internal/template"
	"gorm.io/gen/internal/utils/pools"
)

// T generic type
type T interface{}

// M map[string]interface{}
type M map[string]interface{}

// RowsAffected execute affected raws
type RowsAffected int64

var concurrent = runtime.NumCPU()

func init() { runtime.GOMAXPROCS(runtime.NumCPU()) }

// NewGenerator create a new generator
func NewGenerator(cfg Config) *Generator {
	err := cfg.Revise()
	if err != nil {
		panic(fmt.Errorf("create generator fail: %w", err))
	}

	return &Generator{
		Config:    cfg,
		Data:      make(map[string]*genInfo),
		modelData: make(map[string]*check.BaseStruct),
	}
}

// genInfo info about generated code
type genInfo struct {
	*check.BaseStruct
	Interfaces []*check.InterfaceMethod
}

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

	Data      map[string]*genInfo          //gen query data
	modelData map[string]*check.BaseStruct //gen model data
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
func (g *Generator) GenerateModel(tableName string, opts ...FieldOpt) *check.BaseStruct {
	return g.GenerateModelAs(tableName, g.db.Config.NamingStrategy.SchemaName(tableName), opts...)
}

// GenerateModel catch table info from db, return a BaseStruct
func (g *Generator) GenerateModelAs(tableName string, modelName string, fieldOpts ...FieldOpt) *check.BaseStruct {
	modelFieldOpts := make([]model.FieldOpt, len(fieldOpts))
	for i, opt := range fieldOpts {
		modelFieldOpts[i] = opt
	}
	s, err := check.GenBaseStruct(g.db, model.Conf{
		ModelPkg:       g.Config.ModelPkgPath,
		TablePrefix:    g.getTablePrefix(),
		TableName:      tableName,
		ModelName:      modelName,
		ImportPkgPaths: g.importPkgPaths,
		SchemaNameOpts: g.dbNameOpts,
		TableNameNS:    g.tableNameNS,
		ModelNameNS:    g.modelNameNS,
		FileNameNS:     g.fileNameNS,
		FieldConf: model.FieldConf{
			DataTypeMap: g.dataTypeMap,

			FieldSignable:     g.FieldSignable,
			FieldNullable:     g.FieldNullable,
			FieldCoverable:    g.FieldCoverable,
			FieldWithIndexTag: g.FieldWithIndexTag,
			FieldWithTypeTag:  g.FieldWithTypeTag,

			FieldJSONTagNS: g.fieldJSONTagNS,
			FieldNewTagNS:  g.fieldNewTagNS,

			FieldOpts: modelFieldOpts,
		},
	})
	if err != nil {
		g.db.Logger.Error(context.Background(), "generate struct from table fail: %s", err)
		panic("generate struct fail")
	}
	g.modelData[s.StructName] = s

	g.successInfo(fmt.Sprintf("got %d columns from table <%s>", len(s.Fields), s.TableName))
	return s
}

func (g *Generator) getTablePrefix() string {
	if ns, ok := g.db.NamingStrategy.(schema.NamingStrategy); ok {
		return ns.TablePrefix
	}
	return ""
}

// GenerateAllTable generate all tables in db
func (g *Generator) GenerateAllTable(opts ...FieldOpt) (tableModels []interface{}) {
	tableList, err := g.db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get all tables fail: %w", err))
	}

	g.successInfo(fmt.Sprintf("find %d table from db: %s", len(tableList), tableList))

	tableModels = make([]interface{}, len(tableList))
	for i, tableName := range tableList {
		tableModels[i] = g.GenerateModel(tableName, opts...)
	}
	return tableModels
}

// GenerateModelFrom generate model from object
func (g *Generator) GenerateModelFrom(obj helper.Object) *check.BaseStruct {
	s, err := check.GenBaseStructFromObject(obj, model.Conf{
		ModelPkg:       g.Config.ModelPkgPath,
		ImportPkgPaths: g.importPkgPaths,
		TableNameNS:    g.tableNameNS,
		ModelNameNS:    g.modelNameNS,
		FileNameNS:     g.fileNameNS,
	})
	if err != nil {
		panic(fmt.Errorf("generate struct from object fail: %w", err))
	}
	g.modelData[s.StructName] = s

	g.successInfo(fmt.Sprintf("parse object %s", obj.StructName()))
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
			interfaceStruct.ReviseFieldName()
		}

		data, err := g.pushBaseStruct(interfaceStruct)
		if err != nil {
			g.db.Logger.Error(context.Background(), "gen struct fail: %v", err)
			panic("gen struct fail")
		}

		functions, err := check.BuildDiyMethod(readInterface, interfaceStruct, data.Interfaces)
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
	g.successInfo("Start generating code.")

	if err := g.generateModelFile(); err != nil {
		g.db.Logger.Error(context.Background(), "generate model struct fail: %s", err)
		panic("generate model struct fail")
	}

	if err := g.generateQueryFile(); err != nil {
		g.db.Logger.Error(context.Background(), "generate query code fail: %s", err)
		panic("generate query code fail")
	}

	g.successInfo("Generate code done.")
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
	if len(g.Data) == 0 {
		return nil
	}

	if err := os.MkdirAll(g.OutPath, os.ModePerm); err != nil {
		return fmt.Errorf("make dir outpath(%s) fail: %s", g.OutPath, err)
	}

	errChan := make(chan error)
	pool := pools.NewPool(concurrent)
	// generate query code for all struct
	for _, info := range g.Data {
		pool.Wait()
		go func(info *genInfo) {
			defer pool.Done()
			err = g.generateSingleQueryFile(info)
			if err != nil {
				errChan <- err
			}

			if g.WithUnitTest {
				err = g.generateQueryUnitTestFile(info)
				if err != nil { // do not panic
					g.db.Logger.Error(context.Background(), "generate unit test fail: %s", err)
				}
			}
		}(info)
	}
	select {
	case err = <-errChan:
		return err
	case <-pool.AsyncWaitAll():
	}

	// generate query file
	var buf bytes.Buffer
	err = render(tmpl.Header, &buf, map[string]interface{}{
		"Package":        g.queryPkgName,
		"StructPkgPath":  "",
		"ImportPkgPaths": g.importPkgPaths,
	})
	if err != nil {
		return err
	}

	if g.judgeMode(WithDefaultQuery) {
		err = render(tmpl.DefaultQuery, &buf, g)
		if err != nil {
			return err
		}
	}
	err = render(tmpl.QueryMethod, &buf, g)
	if err != nil {
		return err
	}

	err = g.output(g.OutFile, buf.Bytes())
	if err != nil {
		return err
	}
	g.successInfo("generate query file: " + g.OutFile)

	// generate query unit test file
	if g.WithUnitTest {
		buf.Reset()

		err = render(tmpl.UnitTestHeader, &buf, map[string]interface{}{
			"Package":        g.queryPkgName,
			"StructPkgPath":  "",
			"ImportPkgPaths": g.importPkgPaths,
		})
		if err != nil {
			g.db.Logger.Error(context.Background(), "generate query unit test fail: %s", err)
			return nil
		}
		err = render(tmpl.DIYMethod_TEST_Basic, &buf, nil)
		if err != nil {
			return err
		}
		err = render(tmpl.QueryMethod_TEST, &buf, g)
		if err != nil {
			g.db.Logger.Error(context.Background(), "generate query unit test fail: %s", err)
			return nil
		}
		fileName := strings.TrimSuffix(g.OutFile, ".go") + "_test.go"
		err = g.output(fileName, buf.Bytes())
		if err != nil {
			g.db.Logger.Error(context.Background(), "generate query unit test fail: %s", err)
			return nil
		}
		g.successInfo("generate unit test file: " + fileName)
	}

	return nil
}

// generateSingleQueryFile generate query code and save to file
func (g *Generator) generateSingleQueryFile(data *genInfo) (err error) {
	var buf bytes.Buffer

	structPkgPath := data.StructInfo.PkgPath
	if structPkgPath == "" {
		structPkgPath = g.modelPkgPath
	}
	err = render(tmpl.Header, &buf, map[string]interface{}{
		"Package":        g.queryPkgName,
		"StructPkgPath":  structPkgPath,
		"ImportPkgPaths": data.ImportPkgPaths,
	})
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

	defer g.successInfo(fmt.Sprintf("generate query file: %s/%s.gen.go", g.OutPath, data.FileName))
	return g.output(fmt.Sprintf("%s/%s.gen.go", g.OutPath, data.FileName), buf.Bytes())
}

// generateQueryUnitTestFile generate unit test file for query
func (g *Generator) generateQueryUnitTestFile(data *genInfo) (err error) {
	var buf bytes.Buffer

	structPkgPath := data.StructInfo.PkgPath
	if structPkgPath == "" {
		structPkgPath = g.modelPkgPath
	}
	err = render(tmpl.UnitTestHeader, &buf, map[string]interface{}{
		"Package":        g.queryPkgName,
		"StructPkgPath":  structPkgPath,
		"ImportPkgPaths": data.ImportPkgPaths,
	})
	if err != nil {
		return err
	}

	err = render(tmpl.CRUDMethod_TEST, &buf, data.BaseStruct)
	if err != nil {
		return err
	}

	for _, method := range data.Interfaces {
		err = render(tmpl.DIYMethod_TEST, &buf, method)
		if err != nil {
			return err
		}
	}

	defer g.successInfo(fmt.Sprintf("generate unit test file: %s/%s.gen_test.go", g.OutPath, data.FileName))
	return g.output(fmt.Sprintf("%s/%s.gen_test.go", g.OutPath, data.FileName), buf.Bytes())
}

// generateModelFile generate model structures and save to file
func (g *Generator) generateModelFile() error {
	if len(g.modelData) == 0 {
		return nil
	}

	modelOutPath, err := g.getModelOutputPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(modelOutPath, os.ModePerm); err != nil {
		return fmt.Errorf("create model pkg path(%s) fail: %s", modelOutPath, err)
	}

	errChan := make(chan error)
	pool := pools.NewPool(concurrent)
	for _, data := range g.modelData {
		if data == nil || !data.GenBaseStruct {
			continue
		}
		pool.Wait()
		go func(data *check.BaseStruct) {
			defer pool.Done()
			var buf bytes.Buffer
			err = render(tmpl.Model, &buf, data)
			if err != nil {
				errChan <- err
			}

			modelFile := modelOutPath + data.FileName + ".gen.go"
			err = g.output(modelFile, buf.Bytes())
			if err != nil {
				errChan <- err
			}

			g.successInfo(fmt.Sprintf("generate model file(table <%s> -> {%s.%s}): %s", data.TableName, data.StructInfo.Package, data.StructInfo.Type, modelFile))
		}(data)
	}
	select {
	case err = <-errChan:
		return err
	case <-pool.AsyncWaitAll():
		g.fillModelPkgPath(modelOutPath)
	}
	return nil
}

func (g *Generator) getModelOutputPath() (outPath string, err error) {
	if strings.Contains(g.ModelPkgPath, "/") {
		outPath, err = filepath.Abs(g.ModelPkgPath)
		if err != nil {
			return "", fmt.Errorf("cannot parse model pkg path: %w", err)
		}
	} else {
		outPath = filepath.Dir(g.OutPath) + "/" + g.ModelPkgPath
	}
	return outPath + "/", nil
}

func (g *Generator) fillModelPkgPath(filePath string) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName,
		Dir:  filePath,
	})
	if err != nil {
		g.db.Logger.Warn(context.Background(), "parse model pkg path fail: %s", err)
		return
	}
	if len(pkgs) == 0 {
		g.db.Logger.Warn(context.Background(), "parse model pkg path fail: got 0 packages")
		return
	}
	g.Config.modelPkgPath = pkgs[0].PkgPath
}

// output format and output
func (g *Generator) output(fileName string, content []byte) error {
	result, err := imports.Process(fileName, content, nil)
	if err != nil {
		line := strings.Split(string(content), "\n")
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		startLine, endLine := errLine-5, errLine+5
		fmt.Println("Format fail:", errLine, err)

		for i := startLine; i <= endLine; i++ {
			fmt.Println(i, line[i])
		}
		return fmt.Errorf("cannot format file: %w", err)
	}
	return ioutil.WriteFile(fileName, result, 0640)
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

func render(tmpl string, wr io.Writer, data interface{}) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}
