package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Db2FieldElemListCall func(cf *Config, table string) (data []*TableFieldElem, err error)

// TableElem 每个表的转后的结构
type TableElem struct {
	Table            string
	Data             []*TableFieldElem
	ProtoMessageName string //生成的message名称
}

// TableFieldElem 最终转成的结构
type TableFieldElem struct {
	Table   string `json:"table"`
	Field   string `json:"field"`
	ToField string `json:"to_field"`
	Type    string `json:"type"`
	ToType  string `json:"to_type"`
}

func Table2ProtoMessageName(table string, m *ProtoModel) string {
	return fmt.Sprintf("%s%s%s", m.Prefix, SnakeToCamel(table, true), m.Suffix)
}

// SnakeToCamel 将蛇形命名转换为小驼峰命名
func SnakeToCamel(snakeStr string, firstUpper bool) string {
	var camelStrBuilder strings.Builder
	capitalizeNext := false

	for i, r := range snakeStr {
		if r == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext && i > 0 {
				camelStrBuilder.WriteRune(unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				if i == 0 && firstUpper {
					camelStrBuilder.WriteRune(unicode.ToUpper(r))
				} else {
					camelStrBuilder.WriteRune(unicode.ToLower(r))
				}
			}
		}
	}

	return camelStrBuilder.String()
}

// FilepathExistOrCrate 文件路径是否存在，不存在则创建
func FilepathExistOrCrate(path string, name string) (*os.File, string, error) {
	filepath := fmt.Sprintf("%s/%s", path, name)
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, "", err
		}
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, "", err
		}
	}
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	return f, filepath, err
}

func Gen(f *os.File, data []TableElem, protoContent *strings.Builder) error {
	for _, table := range data {
		generateProtoFile(protoContent, table)
	}
	_, err := f.WriteString(protoContent.String())
	return err
}

func generateProtoFile(protoContent *strings.Builder, table TableElem) {
	protoContent.WriteString(fmt.Sprintf("message %s {\n", table.ProtoMessageName))
	index := 1
	for _, elem := range table.Data {
		protoContent.WriteString(fmt.Sprintf("  %s %s = %d;\n", elem.ToType, elem.ToField, index))
		index++
	}
	protoContent.WriteString("}\n\n\n")
}

// 生成proto结构体数据
func makeTableFieldElemData(cf *Config, table, fieldName, fieldType string) (data *TableFieldElem) {
	elem := &TableFieldElem{
		Table:   table,
		Field:   fieldName,
		ToField: SnakeToCamel(fieldName, false),
		Type:    fieldType,
		ToType:  fieldType,
	}

	if t, ok := cf.typeMapping[cf.Source.Driver][fieldType]; ok {
		elem.ToType = t
	}

	if conf, ok := cf.customMapping[table]; ok {
		if t, ok := conf.Types[fieldType]; ok {
			elem.ToType = t
		}
		for _, e := range conf.Fields {
			if e.Field == fieldName {
				elem.ToType = e.Type
				break
			}
		}
	}
	return elem
}
