package main

import (
	"database/sql"
	"flag"
	"fmt"
)

const (
	mysqlDriver = "mysql"
)

var (
	cf = &Config{}

	dbSQL = &sql.DB{}

	cfReader = NewYaml()

	driversConvert = map[string]Db2FieldElemListCall{
		mysqlDriver: Mysql2FieldElemList,
	}
)

var (
	cfp = flag.String("cfp", "./config", "配置文件路径")
	cff = flag.String("cff", "config.yaml", "配置文件名称")
)

func initConfig(cf *Config) error {

	flag.Parse()

	if *cfp == "" {
		return fmt.Errorf("未指定配置文件路径")
	}
	if *cff == "" {
		return fmt.Errorf("未指定配置文件名称")
	}
	cfPath := fmt.Sprintf("%s/%s", *cfp, *cff)
	body, err := cfReader.Read(cfPath)
	if err != nil {
		return err
	}
	if err = cfReader.Decode(body, cf); err != nil {
		return err
	}

	if cf.typeMapping == nil {
		cf.typeMapping = map[string]map[string]string{
			cf.Source.Driver: make(map[string]string),
		}
	}

	if cf.customMapping == nil {
		cf.customMapping = make(map[string]*CustomElem)
	}

	driverCfPath := fmt.Sprintf("%s/%s.yaml", *cfp, cf.Source.Driver)
	body, err = cfReader.Read(driverCfPath)
	if err != nil {
		return err
	}

	if err = cfReader.Decode(body, cf.typeMapping); err != nil {
		return err
	}

	for _, elem := range cf.Custom {
		cf.customMapping[elem.Table] = elem
	}

	return nil
}

func initSqlDB(cf *Config) (*sql.DB, error) {
	return sql.Open(cf.Source.Driver, cf.Source.DSN)
}
