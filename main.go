package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
)

func main() {

	err := initConfig(cf)
	if err != nil {
		log.Fatalf("解析%s配置失败，error: %v", *cfp, err)
		return
	}

	dbSQL, err = initSqlDB(cf)
	if err != nil {
		log.Fatalf("链接驱动源失败: %s, error: %v", cf.Source.DSN, err)
		return
	}

	call, ok := driversConvert[cf.Source.Driver]
	if !ok {
		log.Fatalf("没有配置对应的转换函数")
		return
	}

	//生成文件
	f, filepath, err := FilepathExistOrCrate(cf.Proto.Path, cf.Proto.Name)
	if err != nil {
		log.Fatalf("创建文件失败: %s, error: %v", filepath, err)
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalf("关闭文件失败: %s, error: %v", filepath, err)
			return
		}
	}()

	err = genProtoes(call, cf, filepath, f)
	if err != nil {
		log.Fatalf("生成文件失败: %s, error: %v", filepath, err)
		return
	}
}

func genProtoes(call Db2FieldElemListCall, cf *Config, filepath string, f *os.File) error {

	protoModels := make([]TableElem, 0, len(cf.Custom))

	//proto头生成
	var protoContent strings.Builder
	for _, str := range cf.Proto.Headers {
		protoContent.WriteString(fmt.Sprintf("%s\n", str))
	}
	protoContent.WriteString(fmt.Sprintf("\n\n"))

	for _, e := range cf.Custom {
		list, err := call(cf, e.Table)
		if err != nil {
			log.Fatal(err)
			return err
		}
		protoModels = append(protoModels, TableElem{
			Table:            e.Table,
			Data:             list,
			ProtoMessageName: Table2ProtoMessageName(e.Table, &cf.Proto.Model),
		})
	}

	if err := Gen(f, protoModels, &protoContent); err != nil {
		log.Fatalf("生成失败: %s, error: %v", filepath, err)
		return err
	}
	return nil
}
