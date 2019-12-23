package util

import (
	"github.com/hetianyi/gox/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ParseYaml parses yaml file.
func ParseYaml(configFile string, target interface{}) error {
	fi, err := file.GetFile(configFile)
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(fi)
	if err != nil {
		return err
	}
	return ParseYamlFromString(bs, target)
}

func ParseYamlFromString(input []byte, target interface{}) error {
	return yaml.Unmarshal(input, target)
}

func Convert2Yaml(c interface{}) ([]byte, error) {
	return yaml.Marshal(c)
}
