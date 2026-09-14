package generate

import (
	"errors"
	"strings"
	"testing"

	"gorm.io/gen/internal/diagnostic"
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
			SQL: "update @@table set modify_time=@wrong.Now() where id=@wrong.Id",
			SplitResult: []string{
				"\"update \"",
				"\"users\"",
				"\" set modify_time=\"",
				"wrong.Now()",
				"\" where id=\"",
				"wrong.Id",
			},
			GenerateResult: []string{
				"params = append(params,wrong.Now())",
				"params = append(params,wrong.Id)",
				"generateSQL.WriteString(\"update users set modify_time=? where id=? \")",
			},
		},
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
				"params = append(params,id)",
				"whereSQL0.WriteString(\"id>? \")",
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
				"params = append(params,id)",
				"whereSQL0.WriteString(\"id>? \")",
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
				"params = append(params,name)",
				"setSQL0.WriteString(\"name=? \")",
				"}",
				"setSQL0.WriteString(\", \")",
				"if id>0 {",
				"params = append(params,id)",
				"setSQL0.WriteString(\"id=? \")",
				"}",
				"helper.JoinSetBuilder(&generateSQL,setSQL0)",
				"params = append(params,id)",
				"generateSQL.WriteString(\"where id=? \")",
			},
		},
		{
			SQL: "select * from @@table {{where}} {{for _, name := range names}}name=@name{{end}}{{end}}",
			SplitResult: []string{
				"\"select * from \"",
				"\"users\"",
				"where",
				"for _, name := range names",
				"\"name=\"",
				"name",
				"end",
				"end",
			},
			GenerateResult: []string{
				"generateSQL.WriteString(\"select * from users \")",
				"var whereSQL0 strings.Builder",
				"for _, name := range names{",
				"params = append(params,name)",
				"whereSQL0.WriteString(\"name=? \")",
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
		{
			Type:    "time.Time",
			Package: "time",
			Name:    "wrong",
		},
	}

	return m

}

func TestCheckSQLVarParamValidation(t *testing.T) {
	newMethod := func() *InterfaceMethod {
		mm := m()
		mm.S = "u"
		mm.InterfaceName = "IF"
		mm.MethodName = "M"
		return mm
	}

	mustFailWith := func(t *testing.T, sql, causePart string) {
		t.Helper()
		mm := newMethod()
		mm.SQLString = sql
		err := mm.sqlStateCheckAndSplit()
		if err == nil {
			t.Fatalf("expected generation-time validation error for %q, got nil", sql)
		}
		var de *diagnostic.Error
		if !errors.As(err, &de) {
			t.Fatalf("expected *diagnostic.Error, got %T: %v", err, err)
		}
		if de.Diag.Code != diagnostic.CodeSQLVar {
			t.Fatalf("expected code %q, got %q", diagnostic.CodeSQLVar, de.Diag.Code)
		}
		if de.Diag.Method != mm.MethodName || de.Diag.Interface != mm.InterfaceName {
			t.Fatalf("diagnostic should carry the failing method: %+v", de.Diag)
		}
		cause := errors.Unwrap(de)
		if cause == nil {
			t.Fatalf("diagnostic should wrap the param-level cause: %v", err)
		}
		if !strings.Contains(cause.Error(), causePart) {
			t.Fatalf("cause should mention %q, got %q", causePart, cause.Error())
		}
	}

	mustGenerate := func(t *testing.T, sql string, wantParts ...string) {
		t.Helper()
		mm := newMethod()
		mm.SQLString = sql
		if err := mm.sqlStateCheckAndSplit(); err != nil {
			t.Fatalf("unexpected validation error for %q: %v", sql, err)
		}
		if _, err := mm.Section.BuildSQL(); err != nil {
			t.Fatalf("unexpected build error for %q: %v", sql, err)
		}
		joined := strings.Join(mm.Section.Tmpls, "\n")
		for _, part := range wantParts {
			if !strings.Contains(joined, part) {
				t.Fatalf("generated code for %q should contain %q, got:\n%s", sql, part, joined)
			}
		}
	}

	t.Run("typo data param rejected", func(t *testing.T) {
		mustFailWith(t, "select * from @@table where name=@nmae", "unknow variable param:nmae")
	})

	t.Run("typo variable param rejected", func(t *testing.T) {
		mustFailWith(t, "select * from @@table where @@nmae = @name", "unknow variable param:nmae")
	})

	t.Run("non string variable param rejected", func(t *testing.T) {
		mustFailWith(t, "select * from @@table where @@id = @id", "variable name must be string")
	})

	t.Run("array variable param rejected", func(t *testing.T) {
		mustFailWith(t, "select * from @@table where @@names = @name", "variable name must be string")
	})

	t.Run("for value and index params accepted", func(t *testing.T) {
		mm := newMethod()
		mm.Params = append(mm.Params, parser.Param{Type: "map[string]string", Name: "attrs"})
		mm.SQLString = "select * from @@table {{where}}{{for k, v := range attrs}}and @@k = @v{{end}}{{end}}"
		if err := mm.sqlStateCheckAndSplit(); err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
		if _, err := mm.Section.BuildSQL(); err != nil {
			t.Fatalf("unexpected build error: %v", err)
		}
		joined := strings.Join(mm.Section.Tmpls, "\n")
		if !strings.Contains(joined, "u.Quote(k)") {
			t.Fatalf("variable for-index should be quoted at generation time:\n%s", joined)
		}
		if !strings.Contains(joined, "params = append(params,v)") {
			t.Fatalf("for-value should be appended to params:\n%s", joined)
		}
	})

	t.Run("string variable param accepted", func(t *testing.T) {
		mustGenerate(t, "select * from @@table where @@name = @name",
			"u.Quote(name)", "params = append(params,name)")
	})

	t.Run("struct field data param accepted", func(t *testing.T) {
		mustGenerate(t, "update @@table set modify_time=@wrong.Now() where id=@wrong.Id",
			"params = append(params,wrong.Now())", "params = append(params,wrong.Id)")
	})
}

func TestCheckResultRejectsInterfaceReturn(t *testing.T) {
	for _, typ := range []string{"interface{}", "any"} {
		m := &InterfaceMethod{
			InterfaceName: "IF",
			MethodName:    "M",
			OriginStruct:  parser.Param{Package: "model", Type: "User"},
		}
		err := m.checkResult([]parser.Param{{Type: typ}, {Type: "error"}})
		if err == nil {
			t.Errorf("checkResult(%s): expected rejection, got nil", typ)
		}
	}
}
