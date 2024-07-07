package deeply

import (
	"reflect"
)

// Equals checks if the expected and actual values are deeply equal.
// It returns true if any of the following conditions are met:
//   - The expected and actual values are both maps and have the same number of keys.
//   - The expected and actual values are both slices and have the same length.
//   - The expected and actual values are deeply equal using reflect.DeepEqual.
func Equals(expect, actual any) bool {
	return mapDeepEqual(expect, actual, Equals) || reflect.DeepEqual(expect, actual)
}

// EqualsIgnoreArrayOrder checks if the expected and actual values are deeply equal
// ignoring the order of arrays. It behaves similarly to Equals except that it
// uses slicesDeepEqualContains instead of slicesDeepEqual to compare slices.
func EqualsIgnoreArrayOrder(expect, actual any) bool {
	return mapDeepEqual(expect, actual, EqualsIgnoreArrayOrder) ||
		slicesDeepEqual(expect, actual, EqualsIgnoreArrayOrder) ||
		reflect.DeepEqual(expect, actual)
}

// mapDeepEqual checks if the expected and actual values are deeply equal as maps.
// It returns true if any of the following conditions are met:
//   - The expected and actual values are both nil.
//   - The expected and actual values are both maps and have the same number of keys.
//   - The expected and actual values are deeply equal using the provided compare function.
func mapDeepEqual(expect, actual any, compare cmp) bool {
	// Check if the types of the expected and actual values are equal.
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	// If both values are nil, return true.
	if reflect.TypeOf(expect) == nil {
		return true
	}

	// If the expected value is not a map, return false.
	if reflect.TypeOf(expect).Kind() != reflect.Map {
		return false
	}

	left := reflect.ValueOf(expect)
	right := reflect.ValueOf(actual)

	// If the number of keys in the maps are not equal, return false.
	if left.Len() != right.Len() {
		return false
	}

	// Compare the values of the maps.
	return mapDeepEquals(left, right, compare)
}

// slicesDeepEqual checks if the expected and actual values are deeply equal as slices.
// It returns true if any of the following conditions are met:
//   - The expected and actual values are both nil.
//   - The expected and actual values are both slices and have the same length.
//   - The expected and actual values are deeply equal using the provided compare function.
func slicesDeepEqual(expect, actual any, compare cmp) bool {
	// Check if the types of the expected and actual values are equal.
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	// If both values are nil, return true.
	if reflect.TypeOf(expect) == nil {
		return true
	}

	// If the expected value is not a slice, return false.
	if reflect.TypeOf(expect).Kind() != reflect.Slice {
		return false
	}

	a := reflect.ValueOf(expect)
	b := reflect.ValueOf(actual)

	// If the lengths of the slices are not equal, return false.
	if a.Len() != b.Len() {
		return false
	}

	// Compare the values of the slices.
	return slicesDeepEqualContains(a, b, compare)
}
