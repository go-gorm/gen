package tests_test

import (
	"io/ioutil"
	"log"
	"regexp"
)

const ddlPath = "tables.sql"

var reg, _ = regexp.Compile(`(DROP TABLE IF EXISTS \x60.*?\x60;)\s(CREATE TABLE [\s\S][^;]*;)`)

func GetDDL() (tableMetas [][2]string) {
	data, err := ioutil.ReadFile(ddlPath)
	if err != nil {
		log.Fatalf("read ddl fail: %s", err)
		return nil
	}

	results := reg.FindAllStringSubmatch(string(data), -1)
	for _, res := range results {
		tableMetas = append(tableMetas, [2]string{res[1], res[2]})
	}
	return tableMetas
}
