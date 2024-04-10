package deeply

import (
	"reflect"
)

type cmp func(expect, actual any) bool

func slicesDeepEquals(expect, actual reflect.Value, compare cmp) bool {
	for i := range expect.Len() {
		if !compare(expect.Index(i).Interface(), actual.Index(i).Interface()) {
			return false
		}
	}

	return true
}

func slicesDeepEqualContains(expect, actual reflect.Value, compare cmp) bool {
	marks := make([]bool, actual.Len())
	res := 0

	for i := range expect.Len() {
		for j := range actual.Len() {
			if !marks[j] && compare(expect.Index(i).Interface(), actual.Index(j).Interface()) {
				marks[j] = true
				res++
			}
		}
	}

	return res == expect.Len()
}

func mapDeepEquals(expect, actual reflect.Value, compare cmp) bool {
	for _, v := range expect.MapKeys() {
		if actual.MapIndex(v).Kind() == reflect.Invalid ||
			!compare(expect.MapIndex(v).Interface(), actual.MapIndex(v).Interface()) {
			return false
		}
	}

	return true
}
