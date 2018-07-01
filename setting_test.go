package fcgirouter

import (
	"gopkg.in/yaml.v2"
	"testing"
	"encoding/json"
)

var testYaml = `/:
  GET: index
/admin:
  /users:
    GET: admin/get_users
    POST: admin/create_user
  /users/:id:
    POST: admin/update_user
  /users/:id/delete:
    POST: admin/delete_user
/topics:
  GET: topic/get_topic
  POST: topic/create_topics
/topics/:id:
  GET: topic/get_topic
/sample/:foo/:bar:
  GET: sample/get
/sample/:foo/:bar/update:
  POST: sample/update
/sample/:foo/:bar/delete:
  POST: sample/delete
/sample2/:foo/:bar:
  GET: sample2/get
/sample2/:foo/:bar/update:
  POST: sample2/update
/sample2/:foo/:bar/delete:
  POST: sample2/delete`

func TestNormalize(t *testing.T) {
	var node map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(testYaml), &node)
	if err != nil {
		t.Fatal(err)
	}
	result := Normalize(node)
	bytes, err := json.Marshal(convertJsonMap(result))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))
}

func convertJsonMap(source map[interface{}]interface{}) map[string]interface{} {
	dest := map[string]interface{}{}
	for k,v := range source {
		key, ok := k.(string)
		if ok {
			if value, ok := v.(string); ok {
				dest[key] = value
			} else if value, ok := v.(map[interface{}]interface{}); ok {
				dest[key] = convertJsonMap(value)
			}
		}
	}
	return dest
}
