package utils

import (
	"reflect"
	"testing"
)

type Predicate func(interface{}, interface{}) bool

func AssertEquals(t *testing.T, expected, actual interface{}, predicates ...Predicate) {
	var isEqual bool

	if len(predicates) > 0 {
		for _, p := range predicates {
			isEqual = p(expected, actual)
			if !isEqual {
				break
			}
		}
	} else {
		isEqual = expected == actual
	}

	if isEqual {
		return
	}

	t.Fatalf("Expected '%v' [%v], but get '%v' [%v]", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
}

func AssertNotEquals(t *testing.T, unexpected, actual interface{}, predicates ...Predicate) {
	var isEqual bool

	if len(predicates) > 0 {
		for _, p := range predicates {
			isEqual = p(unexpected, actual)
			if !isEqual {
				break
			}
		}
	} else {
		isEqual = unexpected == actual
	}

	if !isEqual {
		return
	}

	t.Fatalf("Expected something different than '%v' [%v]", unexpected, reflect.TypeOf(unexpected))
}
