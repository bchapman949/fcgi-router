package fcgirouter

import (
	"testing"
	"gopkg.in/yaml.v2"
	"strings"
)

func TestConvertValue(t *testing.T) {
	var node map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(testYaml), &node)
	if err != nil {
		t.Fatal(err)
	}
	normalized := Normalize(node)
	path, err := ConvertValue(normalized)

	if resolve(path, "GET", "/", map[string]string{}) != "index" {
		t.Error()
	}
	if resolve(path, "GET", "/topics", map[string]string{}) != "topic/get_topic" {
		t.Error()
	}
	args := map[string]string{}
	if resolve(path, "GET", "/sample/a/b", args) != "sample/get" {
		t.Error()
	}
	if args["foo"] != "a" {
		t.Error()
	}
	if args["bar"] != "b" {
		t.Error()
	}
	args = map[string]string{}
	if resolve(path, "POST", "/sample/d/c/update", args) != "sample/update" {
		t.Error()
	}
	if args["foo"] != "d" {
		t.Error()
	}
	if args["bar"] != "c" {
		t.Error()
	}
	args = map[string]string{}
	if resolve(path, "POST", "/sample/e/f/delete", args) != "sample/delete" {
		t.Error()
	}
	if args["foo"] != "e" {
		t.Error()
	}
	if args["bar"] != "f" {
		t.Error()
	}
	if resolve(path, "POST", "/admin/users", map[string]string{}) != "admin/create_user" {
		t.Error()
	}
	args = map[string]string{}
	if resolve(path, "POST", "/admin/users/456", args) != "admin/update_user" {
		t.Error()
	}
	if args["id"] != "456" {
		t.Error()
	}
	args = map[string]string{}
	if resolve(path, "POST", "/admin/users/123/delete", args) != "admin/delete_user" {
		t.Error()
	}
	if args["id"] != "123" {
		t.Error()
	}
}

func resolve(target Path, method string, path string, args map[string]string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) == 0 {
		return target.Match(method, "", "", []string{}, args)
	} else if len(segments) == 1 {
		return target.Match(method, "", segments[0], []string{}, args)
	}
	return target.Match(method, "", segments[0], segments[1:], args)
}
