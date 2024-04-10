package deeply

import (
	"reflect"
)

type cmp func(expect, actual any) bool

func slicesDeepEquals(expect, actual reflect.Value, compare cmp) bool {
	res := 0

	for i := range expect.Len() {
		if compare(expect.Index(i).Interface(), actual.Index(i).Interface()) {
			res++
		}
	}

	return res == expect.Len()
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
	res := 0

	for _, v := range expect.MapKeys() {
		if compare(expect.MapIndex(v).Interface(), actual.MapIndex(v).Interface()) {
			res++
		}
	}

	return res == expect.Len()
}
