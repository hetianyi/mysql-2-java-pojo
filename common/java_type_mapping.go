package common

import (
	"github.com/hetianyi/gox/logger"
	"strings"
)

var (
	typeMapping = make(map[string]string)
)

func init() {
	typeMapping["int"] = "Integer"
	typeMapping["integer"] = "Integer"
	typeMapping["tinyint"] = "Integer"
	typeMapping["smallint"] = "Integer"
	typeMapping["mediumint"] = "Integer"
	typeMapping["bigint"] = "Long"

	typeMapping["double"] = "Double"
	typeMapping["float"] = "Float"

	typeMapping["char"] = "String"
	typeMapping["varchar"] = "String"
	typeMapping["tinytext"] = "String"
	typeMapping["text"] = "String"
	typeMapping["mediumtext"] = "String"
	typeMapping["longtext"] = "String"
	typeMapping["json"] = "String"

	typeMapping["mediumblob"] = "byte[]"
	typeMapping["datetime"] = "java.util.Date"
	typeMapping["date"] = "java.util.Date"
	typeMapping["timestamp"] = "java.sql.Timestamp"
	typeMapping["time"] = "java.sql.Timestamp"
	typeMapping["decimal"] = "java.math.BigDecimal"
	typeMapping["varbinary"] = "byte[]"
	typeMapping["binary"] = "byte[]"
}

func GetType(dbType string) (importPackage, typeName string) {
	if typeMapping[dbType] == "" {
		logger.Fatal("缺少类型映射，请补充：", dbType, " -> ?")
	}
	typeN := typeMapping[dbType]
	if strings.Index(typeN, ".") == -1 {
		return "", typeN
	} else {
		return typeN, typeN[strings.LastIndex(typeN, ".")+1:]
	}
}
