package utils

import (
	"testing"
	"reflect"
)

type Predicate func(interface{}, interface{}) bool

func AssertEquals(t *testing.T, expected, actual interface{}, predicate ...Predicate) {
	var pred Predicate

	if len(predicate) > 0 {
		pred = predicate[0]
	} else {
		pred = func(a interface{}, b interface{}) bool {
			return a == b
		}
	}

	if pred(expected, actual) {
		return
	}

	t.Fatalf("Expected '%v' [%v], but get '%v' [%v]", expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
}

func AssertNotEquals(t *testing.T, unexpected, actual interface{}, predicate ...Predicate) {
	var pred Predicate

	if len(predicate) > 0 {
		pred = predicate[0]
	} else {
		pred = func(a interface{}, b interface{}) bool {
			return a == b
		}
	}

	if !pred(unexpected, actual) {
		return
	}

	t.Fatalf("Expected something different than '%v' (%v)", unexpected, reflect.TypeOf(unexpected))
}
