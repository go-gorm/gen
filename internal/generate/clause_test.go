package generate

import (
	"testing"

	"gorm.io/gen/internal/parser"
)

func checkBuildExpr(t *testing.T, SQL string, splitResult, generateResult []string, i *InterfaceMethod) {
	i.SQLString = SQL
	err := i.sqlStateCheckAndSplit()
	if err != nil {
		t.Errorf("err:%s\n", err)
	}

	if len(i.Section.members) != len(splitResult) {
		t.Errorf("SQL length exp:%v got:%v", len(generateResult), len(i.Section.members))
	}
	for index := range splitResult {
		if splitResult[index] != i.Section.members[index].Value {
			t.Errorf("SQL expects \nexp:%v \ngot:%v", splitResult[index], i.Section.members[index].Value)
		}
	}
	_, err = i.Section.BuildSQL()
	if err != nil {
		t.Errorf("err:%s", err)
	}

	if len(i.Section.Tmpls) != len(generateResult) {
		t.Errorf("SQL length exp:%v got:%v", len(i.Section.Tmpls), len(generateResult))
	}
	for index := range generateResult {
		if generateResult[index] != i.Section.Tmpls[index] {
			t.Errorf("SQL expects \nexp:%v \ngot:%v", generateResult[index], i.Section.Tmpls[index])
		}
	}

}
func TestClause(t *testing.T) {

	testcases := []struct {
		SQL            string
		SplitResult    []string
		GenerateResult []string
	}{
		{
			SQL: "select * from @@table",
			SplitResult: []string{
				"\"select * from \"",
				"\"users\"",
			},
			GenerateResult: []string{
				"generateSQL.WriteString(\"select * from users \")",
			},
		},
		{
			SQL: "select * from @@table {{where}} id>@id{{end}}",
			SplitResult: []string{
				"\"select * from \"",
				"\"users\"",
				"where",
				"\" id>\"",
				"id",
				"end",
			},
			GenerateResult: []string{
				"generateSQL.WriteString(\"select * from users \")",
				"var whereSQL0 strings.Builder",
				"params[\"id\"] = id",
				"whereSQL0.WriteString(\"id>@id \")",
				"helper.JoinWhereBuilder(&generateSQL,whereSQL0)",
			},
		},
		{
			SQL: "select * from @@table {{where}}{{if id > 0}} id>@id{{end}}{{end}}",
			SplitResult: []string{
				"\"select * from \"",
				"\"users\"",
				"where",
				"if id > 0",
				"\" id>\"",
				"id",
				"end",
				"end",
			},
			GenerateResult: []string{
				"generateSQL.WriteString(\"select * from users \")",
				"var whereSQL0 strings.Builder",
				"if id > 0 {",
				"params[\"id\"] = id",
				"whereSQL0.WriteString(\"id>@id \")",
				"}",
				"helper.JoinWhereBuilder(&generateSQL,whereSQL0)",
			},
		},
		{
			SQL: "update @@table {{set}}{{if name != \"\"}}name=@name{{end}},{{if id>0}}id=@id{{end}}{{end}} where id=@id",
			SplitResult: []string{
				"\"update \"",
				"\"users\"",
				"set",
				"if name != \"\"",
				"\"name=\"",
				"name",
				"end",
				"\",\"",
				"if id>0",
				"\"id=\"",
				"id",
				"end",
				"end",
				"\" where id=\"",
				"id",
			},
			GenerateResult: []string{
				"generateSQL.WriteString(\"update users \")",
				"var setSQL0 strings.Builder",
				"if name != \"\" {",
				"params[\"name\"] = name",
				"setSQL0.WriteString(\"name=@name \")",
				"}",
				"setSQL0.WriteString(\", \")",
				"if id>0 {",
				"params[\"id\"] = id",
				"setSQL0.WriteString(\"id=@id \")",
				"}",
				"helper.JoinSetBuilder(&generateSQL,setSQL0)",
				"params[\"id\"] = id",
				"generateSQL.WriteString(\"where id=@id \")",
			},
		},
		{
			SQL: "select * from @@table {{where}} {{for _, name := range names}}name=@name{{end}}{{end}}",
			SplitResult: []string{
				"\"select * from \"",
				"\"users\"",
				"where",
				"for _index, name := range names",
				"\"name=\"",
				"name",
				"end",
				"end",
			},
			GenerateResult: []string{
				"generateSQL.WriteString(\"select * from users \")",
				"var whereSQL0 strings.Builder",
				"for _index, name := range names{",
				"params[\"nameForWhereSQL0_\"+strconv.Itoa(_index)]=name",
				"whereSQL0.WriteString(\"name=@nameForWhereSQL0_\"+strconv.Itoa(_index)+\" \")",
				"}",
				"helper.JoinWhereBuilder(&generateSQL,whereSQL0)",
			},
		},
	}
	inface := m()
	for _, testcase := range testcases {
		checkBuildExpr(t, testcase.SQL, testcase.SplitResult, testcase.GenerateResult, inface)
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
		{
			Type: "string",
			Name: "name",
		},
		{
			Type:    "string",
			Name:    "names",
			IsArray: true,
		},
	}

	return m

}
