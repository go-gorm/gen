package gen

import (
	"bytes"
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

	"gorm.io/gen/internal/check"
	"gorm.io/gen/internal/parser"
	tmpl "gorm.io/gen/internal/template"
)

// TODO implement some unit tests

// T genric type
type T interface{}

// M map[string]interface{}
type M map[string]interface{}

// NewGenerator create a new generator
func NewGenerator(cfg Config) *Generator {
	if cfg.ModelPkgName == "" {
		cfg.ModelPkgName = check.ModelPkg
	}
	return &Generator{
		Config:           cfg,
		Data:             make(map[string]*genInfo),
		readInterfaceSet: new(parser.InterfaceSet),
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

func (i *genInfo) AppendMethods(methods []*check.InterfaceMethod) error {
	for _, newMethod := range methods {
		for _, infoMethod := range i.Interfaces {
			if infoMethod.MethodName == newMethod.MethodName && infoMethod.InterfaceName != newMethod.InterfaceName {
				return fmt.Errorf("can't generate method with the same name from different interface:%s.%s and %s.%s", infoMethod.InterfaceName, infoMethod.MethodName, newMethod.InterfaceName, newMethod.MethodName)
			}
		}
		i.Interfaces = append(i.Interfaces, newMethod)
	}
	return nil
}

// Generator code generator
type Generator struct {
	Config

	Data             map[string]*genInfo
	readInterfaceSet *parser.InterfaceSet
}

// UseDB set db connection
func (g *Generator) UseDB(db *gorm.DB) {
	g.db = db
}

// GenerateModel catch table info from db, return a BaseStruct
func (g *Generator) GenerateModel(tableName string) *check.BaseStruct {
	return g.GenerateModelAs(tableName, g.db.Config.NamingStrategy.SchemaName(tableName))
}

// GenerateModel catch table info from db, return a BaseStruct
func (g *Generator) GenerateModelAs(tableName string, modelName string) *check.BaseStruct {
	s, err := check.GenBaseStructs(g.db, g.Config.ModelPkgName, tableName, modelName, g.dbNameOpts...)
	if err != nil {
		log.Fatalf("check struct error: %s", err)
	}
	return s
}

// ApplyBasic specify models which will implement basic method
func (g *Generator) ApplyBasic(models ...interface{}) {
	g.ApplyInterface(func() {}, models...)
}

// ApplyInterfaces specifies one method interface on several model structures
// eg: g.ApplyInterfaces(model.User{}, func(model.Method1, model.Method2){})
func (g *Generator) ApplyInterfaces(model interface{}, fc interface{}) {
	g.ApplyInterface(fc, model)
}

// ApplyInterface specifies method interfaces on structures, implment codes will be generated after calling g.Execute()
// eg: g.ApplyInterface(func(model.Method){}, model.User{}, model.Company{})
func (g *Generator) ApplyInterface(fc interface{}, models ...interface{}) {
	structs, err := check.CheckStructs(g.db, models...)
	if err != nil {
		log.Fatalf("check struct error: %s", err)
	}
	g.apply(fc, structs)
}

func (g *Generator) apply(fc interface{}, structs []*check.BaseStruct) {
	interfacePaths, err := parser.GetInterfacePath(fc)
	if err != nil {
		log.Fatalf("can't get interface name or file: %s", err)
	}

	err = g.readInterfaceSet.ParseFile(interfacePaths)
	if err != nil {
		log.Fatalf("parser file error: %s", err)
	}

	for _, interfaceStruct := range structs {
		data, err := g.pushBaseStruct(interfaceStruct)
		if err != nil {
			log.Fatalf("gen struct error: %s", err)
		}

		functions, err := check.CheckInterface(g.readInterfaceSet, interfaceStruct)
		if err != nil {
			log.Fatalf("check interface error: %s", err)
		}
		err = data.AppendMethods(functions)
		if err != nil {
			log.Fatalf("gen Interface error: %s", err)
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
			log.Fatalf("mkdir failed: %s", err)
		}
	}
	g.queryPkgName = filepath.Base(g.OutPath)

	err = g.generatedBaseStruct()
	if err != nil {
		log.Fatalf("generate base struct fail: %s", err)
	}
	err = g.generatedToOutFile()
	if err != nil {
		log.Fatalf("generate to file fail: %s", err)
	}
	log.Println("Gorm generated query object file successful!")
	log.Println("Generated path：", g.OutPath)
	log.Println("Generated file：", g.OutFile)
}

// generatedToOutFile save generate code to file
func (g *Generator) generatedToOutFile() (err error) {
	var buf bytes.Buffer

	render := func(tmpl string, wr io.Writer, data interface{}) error {
		t, err := template.New(tmpl).Parse(tmpl)
		if err != nil {
			return err
		}
		return t.Execute(wr, data)
	}

	err = render(tmpl.HeaderTmpl, &buf, g.queryPkgName)
	if err != nil {
		return err
	}

	for _, data := range g.Data {
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
		return fmt.Errorf("can't format generated file: %w", err)
	}
	return outputFile(g.OutFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, result)
}

// generatedBaseStruct generate basic structures
func (g *Generator) generatedBaseStruct() (err error) {
	outPath, err := filepath.Abs(g.OutPath)
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
				log.Fatalf("mkdir failed: %s", err)
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
			return fmt.Errorf("can't format generated file: %w", err)
		}
		err = outputFile(modelFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, result)
		if err != nil {
			return nil
		}
		log.Printf("Generate struct [%s.%s] from table [%s]\n", data.StructInfo.Package, data.StructInfo.Type, data.TableName)
	}
	return nil
}

func (g *Generator) pushBaseStruct(base *check.BaseStruct) (*genInfo, error) {
	structName := base.StructName
	if g.Data[structName] == nil {
		g.Data[structName] = &genInfo{BaseStruct: base}
	}
	if g.Data[structName].Source != base.Source {
		return nil, fmt.Errorf("can't generate struct with the same name from different source:%s.%s and %s.%s", base.StructInfo.Package, base.StructName, g.Data[structName].StructInfo.Package, g.Data[structName].StructName)
	}
	return g.Data[structName], nil
}

func outputFile(filename string, flag int, data []byte) error {
	out, err := os.OpenFile(filename, flag, 0640)
	if err != nil {
		return fmt.Errorf("can't open out file: %w", err)
	}
	return output(out, data)
}

func output(wr io.WriteCloser, data []byte) (err error) {
	defer func() {
		if e := wr.Close(); e != nil {
			err = fmt.Errorf("can't close: %w", e)
		}
	}()

	if _, err = wr.Write(data); err != nil {
		return fmt.Errorf("can't write: %w", err)
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
