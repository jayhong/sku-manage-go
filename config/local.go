package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func (this *Config) loadFromFile(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if strings.Contains(file, "json") {
		return json.Unmarshal(b, &this)
	}

	return yaml.Unmarshal(b, &this)
}
