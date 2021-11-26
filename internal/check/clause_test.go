package check

import (
	"fmt"
	"gorm.io/gen/internal/parser"
	"strconv"
	"strings"
	"testing"
)

func checkBuildExpr(t *testing.T, SQL string, result []string, i *InterfaceMethod) {
	i.SqlString = SQL
	err := i.sqlStateCheck()
	if err != nil {
		t.Errorf("err:%s", err)
	}

	if strings.Join(i.Sections.Tmpl, "\n") != strings.Join(result, "\n") {
		for _, res := range i.Sections.Tmpl {
			fmt.Println(strconv.Quote(res) + ",")
		}
		t.Errorf("Sql expects \nexp:%v \ngot:%v", result, i.Sections.Tmpl)
	}

}
func TestClause(t *testing.T) {

	testcases := []struct {
		SQL    string
		Result []string
	}{
		{
			SQL: "select * from @@table",
			Result: []string{
				"generateSQL.WriteString(\"select * from users\")",
			},
		},
		{
			SQL: "select * from @@table {{where}} id>@id{{end}}",
			Result: []string{
				"generateSQL.WriteString(\"select * from users\")",
				"var whereClause0 strings.Builder",
				"whereClause0.WriteString(\" id>@id\")",
				"generateSQL.WriteString(helper.WhereTrim(whereClause0.String()))",
			},
		},
		{
			SQL: "select * from @@table {{where}}{{if id>0}} id>@id{{end}}{{end}}",
			Result: []string{
				"generateSQL.WriteString(\"select * from users\")",
				"var whereClause0 strings.Builder",
				"if id > 0 {",
				"whereClause0.WriteString(\" id>@id\")",
				"}",
				"generateSQL.WriteString(helper.WhereTrim(whereClause0.String()))",
			},
		},
		{
			SQL: "select * from @@table {{where}}{{if id>0}} id>@id{{end}}{{end}}",
			Result: []string{
				"generateSQL.WriteString(\"select * from users\")",
				"var whereClause0 strings.Builder",
				"if id > 0 {",
				"whereClause0.WriteString(\" id>@id\")",
				"}",
				"generateSQL.WriteString(helper.WhereTrim(whereClause0.String()))",
			},
		},
	}
	inface := m()
	for _, testcase := range testcases {
		checkBuildExpr(t, testcase.SQL, testcase.Result, inface)
	}
}

var m = func() *InterfaceMethod {
	var m = new(InterfaceMethod)
	m.Table = "users"
	m.Params = []parser.Param{
		{
			Type: "int",
			Name: "id",
		},
	}

	return m

}
