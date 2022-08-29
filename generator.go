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
	"gorm.io/gen/internal/generate"
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
	if err := cfg.Revise(); err != nil {
		panic(fmt.Errorf("create generator fail: %w", err))
	}

	return &Generator{
		Config: cfg,
		Data:   make(map[string]*genInfo),
		models: make(map[string]*generate.QueryStructMeta),
	}
}

// genInfo info about generated code
type genInfo struct {
	*generate.QueryStructMeta
	Interfaces []*generate.InterfaceMethod
}

func (i *genInfo) appendMethods(methods []*generate.InterfaceMethod) {
	for _, newMethod := range methods {
		if i.methodInGenInfo(newMethod) {
			continue
		}
		i.Interfaces = append(i.Interfaces, newMethod)
	}
}

func (i *genInfo) methodInGenInfo(m *generate.InterfaceMethod) bool {
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

	Data   map[string]*genInfo                  //gen query data
	models map[string]*generate.QueryStructMeta //gen model data
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
func (g *Generator) GenerateModel(tableName string, opts ...ModelOpt) *generate.QueryStructMeta {
	return g.GenerateModelAs(tableName, g.db.Config.NamingStrategy.SchemaName(tableName), opts...)
}

// GenerateModelAs catch table info from db, return a BaseStruct
func (g *Generator) GenerateModelAs(tableName string, modelName string, opts ...ModelOpt) *generate.QueryStructMeta {
	meta, err := generate.GetQueryStructMeta(g.db, g.genModelConfig(tableName, modelName, opts))
	if err != nil {
		g.db.Logger.Error(context.Background(), "generate struct from table fail: %s", err)
		panic("generate struct fail")
	}
	if meta == nil {
		g.info(fmt.Sprintf("ignore table <%s>", tableName))
		return nil
	}
	g.models[meta.ModelStructName] = meta

	g.info(fmt.Sprintf("got %d columns from table <%s>", len(meta.Fields), meta.TableName))
	return meta
}

// GenerateAllTable generate all tables in db
func (g *Generator) GenerateAllTable(opts ...ModelOpt) (tableModels []interface{}) {
	tableList, err := g.db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get all tables fail: %w", err))
	}

	g.info(fmt.Sprintf("find %d table from db: %s", len(tableList), tableList))

	tableModels = make([]interface{}, len(tableList))
	for i, tableName := range tableList {
		tableModels[i] = g.GenerateModel(tableName, opts...)
	}
	return tableModels
}

// GenerateModelFrom generate model from object
func (g *Generator) GenerateModelFrom(obj helper.Object) *generate.QueryStructMeta {
	s, err := generate.GetQueryStructMetaFromObject(obj, g.genModelObjConfig())
	if err != nil {
		panic(fmt.Errorf("generate struct from object fail: %w", err))
	}
	g.models[s.ModelStructName] = s

	g.info(fmt.Sprintf("parse object %s", obj.StructName()))
	return s
}

func (g *Generator) genModelConfig(tableName string, modelName string, modelOpts []ModelOpt) *model.Config {
	return &model.Config{
		ModelPkg:       g.Config.ModelPkgPath,
		TablePrefix:    g.getTablePrefix(),
		TableName:      tableName,
		ModelName:      modelName,
		ImportPkgPaths: g.importPkgPaths,
		ModelOpts:      modelOpts,
		NameStrategy: model.NameStrategy{
			SchemaNameOpts: g.dbNameOpts,
			TableNameNS:    g.tableNameNS,
			ModelNameNS:    g.modelNameNS,
			FileNameNS:     g.fileNameNS,
		},
		FieldConfig: model.FieldConfig{
			DataTypeMap: g.dataTypeMap,

			FieldSignable:     g.FieldSignable,
			FieldNullable:     g.FieldNullable,
			FieldCoverable:    g.FieldCoverable,
			FieldWithIndexTag: g.FieldWithIndexTag,
			FieldWithTypeTag:  g.FieldWithTypeTag,

			FieldJSONTagNS: g.fieldJSONTagNS,
			FieldNewTagNS:  g.fieldNewTagNS,
		},
	}
}

func (g *Generator) getTablePrefix() string {
	if ns, ok := g.db.NamingStrategy.(schema.NamingStrategy); ok {
		return ns.TablePrefix
	}
	return ""
}

func (g *Generator) genModelObjConfig() *model.Config {
	return &model.Config{
		ModelPkg:       g.Config.ModelPkgPath,
		ImportPkgPaths: g.importPkgPaths,
		NameStrategy: model.NameStrategy{
			TableNameNS: g.tableNameNS,
			ModelNameNS: g.modelNameNS,
			FileNameNS:  g.fileNameNS,
		},
	}
}

// ApplyBasic specify models which will implement basic method
func (g *Generator) ApplyBasic(models ...interface{}) {
	g.ApplyInterface(func() {}, models...)
}

// ApplyInterface specifies method interfaces on structures, implment codes will be generated after calling g.Execute()
// eg: g.ApplyInterface(func(model.Method){}, model.User{}, model.Company{})
func (g *Generator) ApplyInterface(fc interface{}, models ...interface{}) {
	structs, err := generate.ConvertStructs(g.db, models...)
	if err != nil {
		g.db.Logger.Error(context.Background(), "check struct fail: %v", err)
		panic("check struct fail")
	}
	g.apply(fc, structs)
}

func (g *Generator) apply(fc interface{}, structs []*generate.QueryStructMeta) {
	interfacePaths, err := parser.GetInterfacePath(fc)
	if err != nil {
		g.db.Logger.Error(context.Background(), "get interface name or file fail: %s", err)
		panic("check interface fail")
	}

	readInterface := new(parser.InterfaceSet)
	err = readInterface.ParseFile(interfacePaths, generate.GetStructNames(structs))
	if err != nil {
		g.db.Logger.Error(context.Background(), "parser interface file fail: %s", err)
		panic("parser interface file fail")
	}

	for _, interfaceStructMeta := range structs {
		if g.judgeMode(WithoutContext) {
			interfaceStructMeta.ReviseFieldName()
		}

		genInfo, err := g.pushQueryStructMeta(interfaceStructMeta)
		if err != nil {
			g.db.Logger.Error(context.Background(), "gen struct fail: %v", err)
			panic("gen struct fail")
		}

		functions, err := generate.BuildDIYMethod(readInterface, interfaceStructMeta, genInfo.Interfaces)
		if err != nil {
			g.db.Logger.Error(context.Background(), "check interface fail: %v", err)
			panic("check interface fail")
		}
		genInfo.appendMethods(functions)
	}
}

// Execute generate code to output path
func (g *Generator) Execute() {
	g.info("Start generating code.")

	if err := g.generateModelFile(); err != nil {
		g.db.Logger.Error(context.Background(), "generate model struct fail: %s", err)
		panic("generate model struct fail")
	}

	if err := g.generateQueryFile(); err != nil {
		g.db.Logger.Error(context.Background(), "generate query code fail: %s", err)
		panic("generate query code fail")
	}

	g.info("Generate code done.")
}

// info logger
func (g *Generator) info(logInfos ...string) {
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

	if err = os.MkdirAll(g.OutPath, os.ModePerm); err != nil {
		return fmt.Errorf("make dir outpath(%s) fail: %s", g.OutPath, err)
	}

	errChan := make(chan error)
	pool := pools.NewPool(concurrent)
	// generate query code for all struct
	for _, info := range g.Data {
		pool.Wait()
		go func(info *genInfo) {
			defer pool.Done()
			err := g.generateSingleQueryFile(info)
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
		"ImportPkgPaths": importList.Add(g.importPkgPaths...).Paths(),
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
	g.info("generate query file: " + g.OutFile)

	// generate query unit test file
	if g.WithUnitTest {
		buf.Reset()

		err = render(tmpl.Header, &buf, map[string]interface{}{
			"Package":        g.queryPkgName,
			"ImportPkgPaths": unitTestImportList.Add(g.importPkgPaths...).Paths(),
		})
		if err != nil {
			g.db.Logger.Error(context.Background(), "generate query unit test fail: %s", err)
			return nil
		}
		err = render(tmpl.DIYMethodTestBasic, &buf, nil)
		if err != nil {
			return err
		}
		err = render(tmpl.QueryMethodTest, &buf, g)
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
		g.info("generate unit test file: " + fileName)
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
		"ImportPkgPaths": importList.Add(structPkgPath).Add(getImportPkgPaths(data)...).Paths(),
	})
	if err != nil {
		return err
	}

	data.QueryStructMeta = data.QueryStructMeta.IfaceMode(g.judgeMode(WithQueryInterface))

	structTmpl := tmpl.TableQueryStructWithContext
	if g.judgeMode(WithoutContext) {
		structTmpl = tmpl.TableQueryStruct
	}
	err = render(structTmpl, &buf, data.QueryStructMeta)
	if err != nil {
		return err
	}

	if g.judgeMode(WithQueryInterface) {
		err = render(tmpl.TableQueryIface, &buf, data)
		if err != nil {
			return err
		}
	}

	for _, method := range data.Interfaces {
		err = render(tmpl.DIYMethod, &buf, method)
		if err != nil {
			return err
		}
	}

	err = render(tmpl.CRUDMethod, &buf, data.QueryStructMeta)
	if err != nil {
		return err
	}

	defer g.info(fmt.Sprintf("generate query file: %s/%s.gen.go", g.OutPath, data.FileName))
	return g.output(fmt.Sprintf("%s/%s.gen.go", g.OutPath, data.FileName), buf.Bytes())
}

// generateQueryUnitTestFile generate unit test file for query
func (g *Generator) generateQueryUnitTestFile(data *genInfo) (err error) {
	var buf bytes.Buffer

	structPkgPath := data.StructInfo.PkgPath
	if structPkgPath == "" {
		structPkgPath = g.modelPkgPath
	}
	err = render(tmpl.Header, &buf, map[string]interface{}{
		"Package":        g.queryPkgName,
		"ImportPkgPaths": unitTestImportList.Add(structPkgPath).Add(data.ImportPkgPaths...).Paths(),
	})
	if err != nil {
		return err
	}

	err = render(tmpl.CRUDMethodTest, &buf, data.QueryStructMeta)
	if err != nil {
		return err
	}

	for _, method := range data.Interfaces {
		err = render(tmpl.DIYMethodTest, &buf, method)
		if err != nil {
			return err
		}
	}

	defer g.info(fmt.Sprintf("generate unit test file: %s/%s.gen_test.go", g.OutPath, data.FileName))
	return g.output(fmt.Sprintf("%s/%s.gen_test.go", g.OutPath, data.FileName), buf.Bytes())
}

// generateModelFile generate model structures and save to file
func (g *Generator) generateModelFile() error {
	if len(g.models) == 0 {
		return nil
	}

	modelOutPath, err := g.getModelOutputPath()
	if err != nil {
		return err
	}

	if err = os.MkdirAll(modelOutPath, os.ModePerm); err != nil {
		return fmt.Errorf("create model pkg path(%s) fail: %s", modelOutPath, err)
	}

	errChan := make(chan error)
	pool := pools.NewPool(concurrent)
	for _, data := range g.models {
		if data == nil || !data.Generated {
			continue
		}
		pool.Wait()
		go func(data *generate.QueryStructMeta) {
			defer pool.Done()

			var buf bytes.Buffer
			err := render(tmpl.Model, &buf, data)
			if err != nil {
				errChan <- err
				return
			}

			for _, method := range data.ModelMethods {
				err = render(tmpl.ModelMethod, &buf, method)
				if err != nil {
					errChan <- err
					return
				}
			}

			modelFile := modelOutPath + data.FileName + ".gen.go"
			err = g.output(modelFile, buf.Bytes())
			if err != nil {
				errChan <- err
				return
			}

			g.info(fmt.Sprintf("generate model file(table <%s> -> {%s.%s}): %s", data.TableName, data.StructInfo.Package, data.StructInfo.Type, modelFile))
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
		lines := strings.Split(string(content), "\n")
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		startLine, endLine := errLine-5, errLine+5
		fmt.Println("Format fail:", errLine, err)
		if startLine < 0 {
			startLine = 0
		}
		if endLine > len(lines)-1 {
			endLine = len(lines) - 1
		}
		for i := startLine; i <= endLine; i++ {
			fmt.Println(i, lines[i])
		}
		return fmt.Errorf("cannot format file: %w", err)
	}
	return ioutil.WriteFile(fileName, result, 0640)
}

func (g *Generator) pushQueryStructMeta(meta *generate.QueryStructMeta) (*genInfo, error) {
	structName := meta.ModelStructName
	if g.Data[structName] == nil {
		g.Data[structName] = &genInfo{QueryStructMeta: meta}
	}
	if g.Data[structName].Source != meta.Source {
		return nil, fmt.Errorf("cannot generate struct with the same name from different source:%s.%s and %s.%s",
			meta.StructInfo.Package, meta.ModelStructName, g.Data[structName].StructInfo.Package, g.Data[structName].ModelStructName)
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

func getImportPkgPaths(data *genInfo) []string {
	importPathMap := make(map[string]struct{})
	for _, path := range data.ImportPkgPaths {
		importPathMap[path] = struct{}{}
	}
	// imports.Process (called in Generator.output) will guess missing imports, and will be
	// much faster if import path is already specified. So add all imports from DIY interface package.
	for _, method := range data.Interfaces {
		for _, param := range method.Params {
			importPathMap[param.PkgPath] = struct{}{}
		}
	}
	importPkgPaths := make([]string, 0, len(importPathMap))
	for importPath := range importPathMap {
		importPkgPaths = append(importPkgPaths, importPath)
	}
	return importPkgPaths
}
