package ssatest

import (
	"regexp"
	"strconv"
)

type Value struct {
	Reference string
	Value     int
}

var (
	identifierRegex = regexp.MustCompile(`^[a-z]+$`)
	ssaRegex        = regexp.MustCompile(`^[a-z]+(-[0-9]+)?$`)
)

func getValue(value string) Value {
	if identifierRegex.MatchString(value) {
		return Value{Reference: value}
	}

	v, _ := strconv.Atoi(value)

	return Value{Value: v}
}

func getValueSSA(value string) Value {
	if ssaRegex.MatchString(value) {
		return Value{Reference: value}
	}

	v, _ := strconv.Atoi(value)

	return Value{Value: v}
}
