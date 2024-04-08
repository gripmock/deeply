package deeply

import (
	"reflect"
)

func Contains(expect, actual any) bool {
	return mapDeepContains(expect, actual, Contains) || reflect.DeepEqual(expect, actual)
}

func ContainsIgnoreArrayOrder(expect, actual any) bool {
	return mapDeepContains(expect, actual, ContainsIgnoreArrayOrder) ||
		slicesDeepContains(expect, actual, ContainsIgnoreArrayOrder) ||
		reflect.DeepEqual(expect, actual)
}

func mapDeepContains(expect, actual any, compare cmp) bool {
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	if reflect.TypeOf(expect).Kind() != reflect.Map {
		return false
	}

	left := reflect.ValueOf(expect)
	right := reflect.ValueOf(actual)

	if left.Len() > right.Len() {
		return false
	}

	return mapDeepEqualContains(left, right, compare)
}

func slicesDeepContains(expect, actual any, compare cmp) bool {
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	a := reflect.ValueOf(expect)
	b := reflect.ValueOf(actual)

	if a.Kind() != reflect.Slice || b.Kind() != reflect.Slice || a.Len() > b.Len() {
		return false
	}

	return slicesDeepEqualContains(a, b, compare)
}
