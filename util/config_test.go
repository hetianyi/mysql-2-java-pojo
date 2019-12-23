package util

import (
	"github.com/hetianyi/mysql-2-java-pojo/common"
	"io/ioutil"
	"testing"
)

func TestConvert2Yaml(t *testing.T) {
	config := &common.Config{
		Host:                   "localhost",
		Port:                   3306,
		User:                   "root",
		Password:               "123456",
		DB:                     "test",
		Package:                "com.xxx",
		Author:                 "Jason He",
		Version:                "1.0.0",
		UseLombok:              false,
		AddSerializeAnnotation: true,
		BeanSuffix:             "DO",
		Options: map[string]string{
			"charset":   "utf8",
			"parseTime": "True",
			"loc":       "Local",
		},
		UseMybatisPlus: true,
		Tables:         []string{""},
	}

	bs, _ := Convert2Yaml(config)

	ioutil.WriteFile("config.yml", bs, 0666)
}
