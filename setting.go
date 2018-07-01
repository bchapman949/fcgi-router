package fcgirouter

import (
	"strings"
)

func Normalize(node map[interface{}]interface{}) map[interface{}]interface{} {
	result := map[interface{}]interface{}{}
	for k, v := range node {
		key, ok := k.(string)
		if !ok {
			continue
		}
		trimmed := strings.Trim(key, "/")
		var head string
		var tail []string
		if strings.Contains(trimmed, "/") {
			segments := strings.Split(trimmed, "/")
			head = segments[0]
			tail = segments[1:]
		} else {
			head = trimmed
			tail = []string{}
		}

		left, ok := result[head]
		right := ConvertSegments(tail, v)
		if ok {
			result[head] = Merge(left, right)
		} else {
			result[head] = right
		}
	}
	return result
}

func Merge(left interface{}, right interface{}) interface{} {
	leftMap, lok := left.(map[interface{}]interface{})
	rightMap, rok := right.(map[interface{}]interface{})
	if lok && rok {
		for k, v := range rightMap {
			e, ok := leftMap[k]
			if ok {
				leftMap[k] = Merge(e, v)
			} else {
				leftMap[k] = v
			}
		}
	}
	return leftMap
}

func ConvertSegments(tail []string, value interface{}) interface{} {
	if len(tail) == 0 {
		_, ok := value.(string)
		if ok {
			return value
		}
		node, ok := value.(map[interface{}]interface{})
		if ok {
			return Normalize(node)
		}
		return ""
	}

	var nextHead string
	var nextTail []string
	if len(tail) == 1 {
		nextHead = tail[0]
		nextTail = []string{}
	} else {
		nextHead = tail[0]
		nextTail = tail[1:]
	}

	result := map[interface{}]interface{}{}
	result[nextHead] = ConvertSegments(nextTail, value)
	return result
}
