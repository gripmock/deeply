package deeply

import "reflect"

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
	res := 0

	for i := range expect.Len() {
		for j := range actual.Len() {
			if compare(expect.Index(i).Interface(), actual.Index(j).Interface()) {
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

func mapDeepEqualContains(expect, actual reflect.Value, compare cmp) bool {
	res := 0

	for _, v1 := range expect.MapKeys() {
		for _, v2 := range actual.MapKeys() {
			if compare(expect.MapIndex(v1).Interface(), actual.MapIndex(v2).Interface()) {
				res++
			}
		}
	}

	return res == expect.Len()
}
