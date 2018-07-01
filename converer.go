package fcgirouter

import (
	"errors"
	"strings"
)

func ConvertValue(value interface{}) (Path, error) {
	n, ok := value.(map[interface{}]interface{})
	if ok {
		return ConvertNode(n)
	}
	s, ok := value.(string)
	if ok {
		return leafPath{s}, nil
	}
	return nil, errors.New("invalid value")
}

func ConvertNode(node map[interface{}]interface{}) (Path, error) {
	p := nodePath{
		staticPath: map[string]Path{},
		methods:    map[string]Path{},
		argPath:    []Path{},
	}
	for k, v := range node {
		key, ok := k.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(key, ":") {
			s, err := ConvertValue(v)
			if err != nil {
				return nil, err
			}
			p.argPath = append(p.argPath, argPath{strings.Trim(key, ":"), s})
			continue
		}

		switch key {
		case "GET":
			fallthrough
		case "POST":
			fallthrough
		case "PUT":
			fallthrough
		case "DELETE":
			s, err := ConvertValue(v)
			if err != nil {
				return nil, err
			}
			p.methods[key] = s
		default:
			s, err := ConvertValue(v)
			if err != nil {
				return nil, err
			}
			p.staticPath[key] = s
		}
	}
	return p, nil
}

type Path interface {
	Match(method string, key string, head string, tail []string, args map[string]string) string
}

type nodePath struct {
	staticPath map[string]Path
	methods    map[string]Path
	argPath    []Path
}

func (p nodePath) Match(method string, key string, head string, tail []string, args map[string]string) string {
	var single Path
	var multiple []Path
	sPath, ok := p.staticPath[head]
	for {
		if ok {
			single = sPath
			break
		}
		if head == "" {
			m, ok := p.methods[method]
			if ok {
				single = m
				break
			}
		}
		if p.argPath != nil {
			multiple = p.argPath
			break
		}
		break
	}
	if single != nil {
		var value string
		if len(tail) == 0 {
			value = single.Match(method, head, "", nil, args)
		} else if len(tail) == 1 {
			value = single.Match(method, head, tail[0], nil, args)
		} else {
			value = single.Match(method, head, tail[0], tail[1:], args)
		}
		return value
	}
	for _, p := range multiple {
		var value string
		if len(tail) == 0 {
			value = p.Match(method, head, "", nil, args)
		} else if len(tail) == 1 {
			value = p.Match(method, head, tail[0], nil, args)
		} else {
			value = p.Match(method, head, tail[0], tail[1:], args)
		}
		if value != "" {
			return value
		}
	}
	return ""
}

type argPath struct {
	key   string
	child Path
}

func (p argPath) Match(method string, key string, head string, tail []string, args map[string]string) string {
	args[p.key] = key
	return p.child.Match(method, "", head, tail, args)
}

type leafPath struct {
	path string
}

func (p leafPath) Match(method string, key string, head string, tail []string, args map[string]string) string {
	return p.path
}
