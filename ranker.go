package deeply

import (
	"reflect"
	"regexp"
)

type ranker func(expect, actual interface{}) float64

func RankMatch(expect, actual any) float64 {
	return mapRankMatch(expect, actual, RankMatch) +
		slicesRankMatch(expect, actual, RankMatch) +
		rank(expect, actual)
}

func rank(expect, actual interface{}) float64 {
	var (
		expectedStr, expectedStringOk = expect.(string)
		actualStr, actualStringOk     = actual.(string)
	)

	if !expectedStringOk || !actualStringOk {
		if reflect.DeepEqual(expect, actual) {
			return 1
		}

		return 0
	}

	if expectedStr == actualStr {
		return 1
	}

	compile, err := regexp.Compile(expectedStr)
	if compile != nil && err == nil {
		results := compile.FindStringIndex(actualStr)

		if len(results) == 2 { //nolint:gomnd
			return float64(results[1]-results[0]) / float64(len(actualStr))
		}
	}

	return distance(expectedStr, actualStr)
}

func mapRankMatch(expect, actual any, compare ranker) float64 {
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return 0
	}

	if reflect.TypeOf(expect) == nil {
		return 1
	}

	if reflect.TypeOf(expect).Kind() != reflect.Map {
		return 0
	}

	left := reflect.ValueOf(expect)
	right := reflect.ValueOf(actual)

	var res float64

	for _, v := range right.MapKeys() {
		if left.MapIndex(v).IsValid() && left.MapIndex(v).CanInterface() {
			res += compare(left.MapIndex(v).Interface(), right.MapIndex(v).Interface())
		}
	}

	return res / float64(max(left.Len(), right.Len()))
}

func slicesRankMatch(expect, actual any, compare ranker) float64 {
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return 0
	}

	if reflect.TypeOf(expect) == nil {
		return 1
	}

	a := reflect.ValueOf(expect)
	b := reflect.ValueOf(actual)

	if a.Kind() != reflect.Slice || b.Kind() != reflect.Slice {
		return 0
	}

	var res float64

	for i := range max(a.Len(), b.Len()) {
		if i < b.Len() && b.Index(i).IsValid() && b.Index(i).CanInterface() {
			res += compare(a.Index(i).Interface(), b.Index(i).Interface())
		}
	}

	return res / float64(max(a.Len(), b.Len()))
}

func distance(s, t string) float64 {
	r1, r2 := []rune(s), []rune(t)
	column := make([]int, 1, 64) //nolint:gomnd

	for y := 1; y <= len(r1); y++ {
		column = append(column, y)
	}

	for x := 1; x <= len(r2); x++ {
		column[0] = x

		for y, lastDiag := 1, x-1; y <= len(r1); y++ {
			oldDiag := column[y]

			cost := 0
			if r1[y-1] != r2[x-1] {
				cost = 1
			}

			column[y] = min(column[y]+1, column[y-1]+1, lastDiag+cost)
			lastDiag = oldDiag
		}
	}

	length := float64(max(len(s), len(t)))

	return (length - float64(column[len(r1)])) / length
}
