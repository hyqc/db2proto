package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	//数据源配置
	Source Source `yaml:"source"`
	//生成的proto配置
	Proto Proto `yaml:"proto"`
	//用户的自定义配置
	Custom Custom `yaml:"custom"`
	//数据源的数据类型映射proto类型, string1=driver,string2=driver的字段类型
	typeMapping map[string]map[string]string
	//自定义映射
	customMapping map[string]*CustomElem
}

type Source struct {
	//数据库驱动名称
	Driver string `yaml:"driver"`
	//数据源名称
	DSN string `yaml:"dsn"`
	//数据库名称
	Dbname string `yaml:"dbname"`
}

type Proto struct {
	//生成的proto文件路径
	Path string `yaml:"path"`
	//生成的proto文件名称
	Name string `yaml:"name"`
	//proto的header配置
	Headers []string `yaml:"headers"`
	//生成的proto对应model的前缀和后缀,示例表web_account，prefix=Test suffix=Model，
	//则生成TestWebAccountModel的proto的message
	Model ProtoModel `yaml:"model"`
}

type Custom []*CustomElem

type ProtoModel struct {
	Prefix string `yaml:"prefix"`
	Suffix string `yaml:"suffix"`
}

type CustomElem struct {
	//表名
	Table string `yaml:"table"`
	//中的字段名称
	Fields []CustomFieldsElem `yaml:"fields"`
	//表中的字段类型，优先级fields-> types
	Types map[string]string `yaml:"types"`
}

type CustomFieldsElem struct {
	//表中的字段名称
	Field string `yaml:"field"`
	//表中的类型转成proto的类型
	Type string `yaml:"type"`
}

type IConfig interface {
	Read(source string) ([]byte, error)
	Decode(data []byte, conf any) error
}

type Yaml struct {
}

func NewYaml() IConfig {
	return &Yaml{}
}

func (Yaml) Read(path string) ([]byte, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (Yaml) Decode(data []byte, conf any) error {
	err := yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	return nil
}
