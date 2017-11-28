package utils

import (
	"path"
	"reflect"
	"runtime"
	"testing"
	"fmt"
)

type Predicate func(interface{}, interface{}) bool
type Message string
type PredicateDefinition struct {
	Predicate
	Message
}

func AssertEquals(t *testing.T, expected, actual interface{}, predicates ...PredicateDefinition) {
	var isEqual bool
	var reason string

	if len(predicates) > 0 {
		for _, p := range predicates {
			isEqual = p.Predicate(expected, actual)
			if !isEqual {
				reason = fmt.Sprintf(string(p.Message), expected, actual)
				break
			}
		}
	} else {
		reason = fmt.Sprintf("Expected '%v' [%v], but get '%v' [%v]", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
		isEqual = expected == actual
	}

	if isEqual {
		return
	}

	_, caller, line, _ := runtime.Caller(1)
	t.Fatalf("%s (at %s:%d)", reason, path.Base(caller), line)
}

func AssertNotEquals(t *testing.T, unexpected, actual interface{}, predicates ...PredicateDefinition) {
	var isEqual bool
	var reason string

	if len(predicates) > 0 {
		for _, p := range predicates {
			isEqual = p.Predicate(unexpected, actual)
			if !isEqual {
				reason = fmt.Sprintf(string(p.Message), unexpected, actual)
				break
			}
		}
	} else {
		reason = fmt.Sprintf("Expected '%v' [%v], but get '%v' [%v]", unexpected, reflect.TypeOf(unexpected), actual, reflect.TypeOf(actual))
		isEqual = unexpected == actual
	}

	if !isEqual {
		return
	}

	_, caller, line, _ := runtime.Caller(1)
	t.Fatalf("%s (at %s:%d)", reason, path.Base(caller), line)
}
