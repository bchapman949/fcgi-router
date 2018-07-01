package fcgirouter

import (
	"testing"
	"gopkg.in/yaml.v2"
)

func BenchmarkRouter_Resolve(b *testing.B) {
	var node map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(testYaml), &node)
	if err != nil {
		b.Fatal(err)
	}
	router := NewRouter(node)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.Resolve("POST", "/admin/users/123/update")
	}
}