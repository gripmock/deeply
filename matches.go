package deeply

import (
	"log"
	"reflect"
	"regexp"
)

func Matches(expect, actual any) bool {
	return mapDeepMatches(expect, actual, Matches) ||
		slicesDeepMatches(expect, actual, Matches) ||
		regexMatch(expect, actual) ||
		reflect.DeepEqual(expect, actual)
}

func MatchesIgnoreArrayOrder(expect, actual any) bool {
	return mapDeepMatches(expect, actual, MatchesIgnoreArrayOrder) ||
		slicesDeepContains(expect, actual, MatchesIgnoreArrayOrder) ||
		regexMatch(expect, actual) ||
		reflect.DeepEqual(expect, actual)
}

func slicesDeepMatches(expect, actual any, compare cmp) bool {
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return false
	}

	a := reflect.ValueOf(expect)
	b := reflect.ValueOf(actual)

	if a.Kind() != reflect.Slice || b.Kind() != reflect.Slice || a.Len() > b.Len() {
		return false
	}

	return slicesDeepEquals(a, b, compare)
}

func mapDeepMatches(expect, actual any, compare cmp) bool {
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

	return mapDeepEquals(left, right, compare)
}

func regexMatch(expect, actual interface{}) bool {
	var (
		expectedStr, expectedStringOk = expect.(string)
		actualStr, actualStringOk     = actual.(string)
	)

	if !expectedStringOk || !actualStringOk {
		return false
	}

	match, err := regexp.MatchString(expectedStr, actualStr)
	if err != nil {
		log.Printf("Error on matching regex %s with %s error:%v\n", expect, actual, err)
	}

	return match
}
