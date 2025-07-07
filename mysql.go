package main

import (
	"fmt"
	"strings"
)

// Mysql2FieldElemList 根据数据源读取配置
func Mysql2FieldElemList(cf *Config, table string) (data []*TableFieldElem, err error) {
	rows, err := dbSQL.Query(fmt.Sprintf("DESCRIBE %s", table))
	if err != nil {
		return
	}
	for rows.Next() {
		var (
			fieldName, fieldType  string
			null, key, def, extra any
		)
		if err = rows.Scan(&fieldName, &fieldType, &null, &key, &def, &extra); err != nil {
			return
		}
		index := strings.Index(fieldType, "(")
		if index == -1 {
			index = len(fieldType)
		}
		fieldType = fieldType[:index]

		data = append(data, makeTableFieldElemData(cf, table, fieldName, fieldType))
	}
	return
}
