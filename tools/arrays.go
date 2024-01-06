package tools

import "strings"

func CopyStringArraySkipEmpty(source []string) []string {
	if len(source) == 0 {
		return source
	}
	result := make([]string, 0, len(source))
	for _, str := range source {
		str = strings.TrimSpace(str)
		if len(str) > 0 {
			result = append(result, str)
		}
	}
	if len(result) == len(source) {
		return result
	}
	res := make([]string, len(result))
	copy(res, result)
	return res
}
func CopyArray[T any](source []T) []T {
	if len(source) == 0 {
		return source
	}
	result := make([]T, len(source))
	copy(result, source)
	return result
}
