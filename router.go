package fcgirouter

import "strings"

type Router struct {
	root Path
}

func NewRouter(yamlData map[interface{}]interface{}) *Router {
	normalized := Normalize(yamlData)
	path, err := ConvertValue(normalized)
	if err != nil {
		return nil
	}
	return &Router{
		root: path,
	}
}

func (r Router) Resolve(method string, path string) (string, map[string]string) {
	args := map[string]string{}
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) == 0 {
		return r.root.Match(method, "", "", nil, args), args
	} else if len(segments) == 1 {
		return r.root.Match(method, "", segments[0], nil, args), args
	}
	return r.root.Match(method, "", segments[0], segments[1:], args), args
}