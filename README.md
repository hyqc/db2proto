# db2proto

生成指定数据库的表的对应protobuf文件

## 目标

1. 可配置数据源
2. 可配置生成protobuf的：
    - 路径
    - 名称
    - 头信息
    - message的名称前缀、后缀
3. 可配置指定表的：**（优先级：字段的protobuf类型映射 > 类型的protobuf类型映射）**
    - 指定类型的protobuf类型映射
    - 指定字段的protobuf类型映射

## 命令行

- cfp: 配置文件的路径（文件夹地址，默认路径：./config）
- cff: 配置文件的文件名（默认：config.yaml）

## 配置说明

#### 完整配置示例
文件夹./proto为生成的示例文件，其中web_proto.proto为自动生成，common.proto为手动定义

- config.yaml: 数据源配置

```yaml
#数据源
source:
  #数据库驱动名称
  driver: mysql
  #数据库名称
  dbname: website
  #数据源名称
  dsn: root:123456@tcp(192.168.9.111:3306)/website
#生成的proto配置
proto:
  #生成的proto文件路径
  path: ./proto
  #生成的proto文件名称
  name: web_proto.proto
  #proto的header配置
  headers:
    - syntax = "proto3";
    - package model;
    - option go_package = "./model_proto;model";
    - import "common.proto";
    #- import \"google/protobuf/timestamp.proto\";
  #生成的proto对应model的前缀和后缀,示例表web_account，prefix=Test suffix=Model，则生成TestWebAccountModel的proto的message
  model:
    prefix: DB
    suffix: Model
#自定义数据源 对特定表的特定字段自定义类型
custom:
  #表名
  - table: web_account
    #表中的字段类型，优先级fields-> types
    types:
      datetime: int64
      date: int64
    #表中的特定字段
    fields:
      #表中的字段名称
      - field: bind_at
        #表中的类型转成proto的类型
        type: int64
      - field: created_at
        type: int64
      - field: updated_at
        type: int64
  - table: web_index_config
    types:
      datetime: int64
      date: int64
    fields:
      - field: channels
        type: repeated ChannelElem
```

- mysql.yaml: 对应的数据库的类型映射

```
#mysql 字段类型到 proto字段类型的映射  dsn: root:123456@tcp(192.168.9.111:3306)/website
mysql:
  int: int32
  tinyint: int32
  smallint: int32
  mediumint: int32
  bigint: int64
  float: float
  double: double
  decimal: string
  varchar: string
  char: string
  text: string
  longtext: string
  date: int64
  datetime: int64
  timestamp: int64
  time: string
  blob: bytes
  json: string

```


#### 数据源

- source: 数据源的键名
    - driver: 数据库驱动名称
    - dbname: 数据库名称
    - dsn: 数据源链接配置

- 示例：

```yaml
#数据源
source:
  #数据库驱动名称
  driver: mysql
  #数据库名称
  dbname: website
  #数据源链接配置，以%s结束，用于
  dsn: root:123456@tcp(192.168.9.111:3306)/website
```

#### 生成的proto配置

- proto: 生成的proto配置键名
    - path: 生成的proto文件路径
    - name: 生成的proto文件名称
    - headers: protobuf文件的头（数据类型：字符串数组）
    - model: 生成的表对应的protobuf的message名称的前缀、后缀
        - prefix: 前缀
        - suffix: 后缀

- 示例：

```yaml
#生成的proto配置
proto:
  #生成的proto文件路径
  path: ./proto
  #生成的proto文件名称
  name: web_proto.proto
  #proto的header配置
  headers:
    - syntax = "proto3";
    - package model;
    - option go_package = "./model_proto;model";
    - import "common.proto";
    #- import \"google/protobuf/timestamp.proto\";
  #生成的proto对应model的前缀和后缀,示例表web_account，prefix=Test suffix=Model，则生成TestWebAccountModel的proto的message
  model:
    prefix: DB
    suffix: Model
```

#### 自定义类型映射

- custom: 自定义配置键名（数据类型：结构体数组）
    - table: 自定义配置的表名称
    - types:  表中的字段类型，优先级fields-> types
    - fields: 自定义配置的字段数组（数据类型：结构体数组）
        - field: 自定义类型的字段
        - type: field字段转protobuf的message的字段类型

- 示例：

```yaml
#自定义数据源 对特定表的特定字段自定义类型
custom:
  #表名
  - table: web_account
    #表中的字段类型，优先级fields-> types
    types:
    datetime: int64
    date: int64
    #表中的特定字段
    fields:
      #表中的字段名称
      - field: bind_at
        #表中的类型转成proto的类型
        type: int64
      - field: created_at
        type: int64
      - field: updated_at
        type: int64
  - table: web_index_config
    types:
    datetime: int64
    date: int64
    fields:
      - field: channels
        type: repeated ChannelElem
```

#### 数据源对应的默认类型映射

配置数据源的默认类型映射配置。如果source的driver配置是mysql，则需要单独配置 mysql.yaml 文件作为 mysql的默认映射

- xxx: 数据源名称，如 mysql 等 （数据类型：数据源的类型:protobuf的类型 map 映射）
  int: int32 （示例）

示例：（数据源为mysql）

```yaml
#mysql 字段类型到 proto字段类型的映射（map类型）
mysql:
  int: int32
  tinyint: int32
  smallint: int32
  mediumint: int32
  bigint: int64
  float: float
  double: double
  decimal: string
  varchar: string
  char: string
  text: string
  longtext: string
  date: int64
  datetime: int64
  timestamp: int64
  time: string
  blob: bytes
  json: string
```

