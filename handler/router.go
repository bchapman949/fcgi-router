package handler

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/kr9ly/fcgirouter"
)

func loadRouter(path string) (*fcgirouter.Router, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var node map[interface{}]interface{}
	if err := yaml.Unmarshal([]byte(bytes), &node); err != nil {
		return nil, err
	}

	return fcgirouter.NewRouter(node), nil
}
