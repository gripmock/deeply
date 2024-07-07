package deeply

import (
	"reflect"
)

// Contains checks if the expected value is contained in the actual value.
// It returns true if any of the following conditions are met:
//   - The expected and actual values are deeply equal using reflect.DeepEqual.
//   - The expected and actual values are maps and all keys and values in the expected map
//     are contained in the actual map.
//   - The expected and actual values are slices and the expected slice is completely
//     contained in the actual slice.
func Contains(expect, actual any) bool {
	return mapDeepContains(expect, actual, Contains) || reflect.DeepEqual(expect, actual)
}

// ContainsIgnoreArrayOrder checks if the expected value is contained in the actual value.
// It returns true if any of the following conditions are met:
//   - The expected and actual values are deeply equal using reflect.DeepEqual.
//   - The expected and actual values are maps and all keys and values in the expected map
//     are contained in the actual map.
//   - The expected and actual values are slices and the expected slice is partially
//     contained in the actual slice. The order of elements in the slice is not important.
func ContainsIgnoreArrayOrder(expect, actual any) bool {
	return mapDeepContains(expect, actual, ContainsIgnoreArrayOrder) ||
		slicesDeepContains(expect, actual, ContainsIgnoreArrayOrder) ||
		reflect.DeepEqual(expect, actual)
}

// mapDeepContains checks if the expected map is contained in the actual map.
// It returns true if all keys and values in the expected map are contained in the actual map.
func mapDeepContains(expect, actual any, compare cmp) bool {
	// Check if the types are the same.
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	// Check if both values are nil.
	if reflect.TypeOf(expect) == nil {
		return true
	}

	// Check if both values are maps.
	if reflect.TypeOf(expect).Kind() != reflect.Map {
		return false
	}

	left := reflect.ValueOf(expect)
	right := reflect.ValueOf(actual)

	// Check if the length of the expected map is less than or equal to the length of the actual map.
	if left.Len() > right.Len() {
		return false
	}

	return mapDeepEquals(left, right, compare)
}

// slicesDeepContains checks if the expected slice is contained in the actual slice.
// It returns true if the expected slice is completely contained in the actual slice.
func slicesDeepContains(expect, actual any, compare cmp) bool {
	// Check if the types are the same.
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	// Check if both values are nil.
	if reflect.TypeOf(expect) == nil {
		return true
	}

	// Check if both values are slices.
	if reflect.TypeOf(expect).Kind() != reflect.Slice {
		return false
	}

	a := reflect.ValueOf(expect)
	b := reflect.ValueOf(actual)

	// Check if the length of the expected slice is less than or equal to the length of the actual slice.
	if a.Len() > b.Len() {
		return false
	}

	return slicesDeepEqualContains(a, b, compare)
}
