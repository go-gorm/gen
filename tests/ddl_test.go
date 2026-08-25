package tests_test

import (
	"fmt"
	"os"
	"regexp"
)

const ddlPath = "tables.sql"

var reg, _ = regexp.Compile(`(DROP TABLE IF EXISTS \x60.*?\x60;)\s(CREATE TABLE [\s\S][^;]*;)`)

func GetDDL() (tableMetas [][2]string, err error) {
	data, err := os.ReadFile(ddlPath)
	if err != nil {
		return nil, fmt.Errorf("read DDL: %w", err)
	}

	results := reg.FindAllStringSubmatch(string(data), -1)
	for _, res := range results {
		tableMetas = append(tableMetas, [2]string{res[1], res[2]})
	}
	if len(tableMetas) == 0 {
		return nil, fmt.Errorf("no table definitions found in %s", ddlPath)
	}
	return tableMetas, nil
}
