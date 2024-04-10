package deeply

import (
	"reflect"
)

func Equals(expect, actual any) bool {
	return mapDeepEqual(expect, actual, Equals) || reflect.DeepEqual(expect, actual)
}

func EqualsIgnoreArrayOrder(expect, actual any) bool {
	return mapDeepEqual(expect, actual, EqualsIgnoreArrayOrder) ||
		slicesDeepEqual(expect, actual, EqualsIgnoreArrayOrder) ||
		reflect.DeepEqual(expect, actual)
}

func mapDeepEqual(expect, actual any, compare cmp) bool {
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	if reflect.TypeOf(expect) == nil {
		return true
	}

	if reflect.TypeOf(expect).Kind() != reflect.Map {
		return false
	}

	left := reflect.ValueOf(expect)
	right := reflect.ValueOf(actual)

	if left.Len() != right.Len() {
		return false
	}

	return mapDeepEquals(left, right, compare)
}

func slicesDeepEqual(expect, actual any, compare cmp) bool {
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	if reflect.TypeOf(expect) == nil {
		return true
	}

	if reflect.TypeOf(expect).Kind() != reflect.Slice {
		return false
	}

	a := reflect.ValueOf(expect)
	b := reflect.ValueOf(actual)

	if a.Len() != b.Len() {
		return false
	}

	return slicesDeepEqualContains(a, b, compare)
}
