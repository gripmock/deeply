package deeply

import (
	"reflect"
	"regexp"
)

// Ranker is a function type used to rank matches between two values.
type ranker func(expect, actual interface{}) float64

// RankMatch is the main function used to rank matches between two values.
//
// This function recursively ranks the matches between maps and slices, and then
// ranks the matches between the remaining values. The function returns the sum
// of the matches between the maps and slices, and the matches between the
// remaining values.
//
// Parameters:
//   - expect: The expected value.
//   - actual: The actual value.
//
// Returns:
//   - The total match score between the expected and actual values.
func RankMatch(expect, actual any) float64 {
	// Initialize the total score to 0.
	result := 0.0

	// Rank the matches between the remaining values.
	result += rank(expect, *&actual)

	// Call slicesRankMatch to rank the matches between the maps and slices.
	result += slicesRankMatch(expect, *&actual, RankMatch)

	// Call mapRankMatch to rank the matches between the maps and slices.
	result += mapRankMatch(expect, *&actual, RankMatch)

	// Return the sum of the matches between the maps and slices, and the
	// matches between the remaining values.
	return result
}

// rank is a function that ranks the matches between two strings.
//
// It compares two strings and returns a float64 representing the match score.
//
// Parameters:
// - expect: The expected string.
// - actual: The actual string.
//
// Returns:
// - The match score between the expected and actual strings.
func rank(expect, actual interface{}) float64 {
	// Convert the expected and actual values to strings.
	var (
		expectedStr, expectedStringOk = expect.(string)
		actualStr, actualStringOk     = actual.(string)
	)

	// If the values are not strings, check if they are equal and return the
	// corresponding match score.
	if !expectedStringOk || !actualStringOk {
		if reflect.DeepEqual(expect, actual) {
			return 1 // Full match.
		}

		return 0 // No match.
	}

	// If the strings are equal, return the full match score.
	if expectedStr == actualStr {
		return 1
	}

	// Try to compile the expected string as a regular expression and find the
	// first match in the actual string. Return the match score based on the
	// length of the match.
	compile, err := regexp.Compile(expectedStr)
	if compile != nil && err == nil {
		results := compile.FindStringIndex(actualStr)

		// If a match is found, calculate the match score based on the length of
		// the match.
		if len(results) == 2 { //nolint:mnd
			return float64(results[1]-results[0]) / float64(len(actualStr))
		}
	}

	// If no match is found, calculate the match score based on the Levenshtein
	// distance between the two strings.
	return distance(expectedStr, actualStr)
}

// mapRankMatch calculates the match score between two maps.
//
// It iterates over the keys of the left map and finds the corresponding key in
// the right map. If a match is found, it calculates the match score between
// the values of the keys and adds it to the total score. It marks the keys
// that have been matched to avoid duplicate matches. The function returns the
// total score divided by the maximum number of keys in the two maps.
//
// Parameters:
//   - expect: The expected map.
//   - actual: The actual map.
//   - compare: The ranker function used to compare values.
//
// Returns:
//   - The match score between the expected and actual maps.
func mapRankMatch(expect, actual any, compare ranker) float64 {
	// Check if the types of the expected and actual values are the same.
	// If they are not, return 0.
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return 0
	}

	// Check if the types of the expected and actual values are nil.
	// If they are, return 1.
	if reflect.TypeOf(expect) == nil {
		return 1
	}

	// Check if the expected value is a map.
	// If it is not, return 0.
	if reflect.TypeOf(expect).Kind() != reflect.Map {
		return 0
	}

	// Convert the expected and actual values to reflect.Value.
	left := reflect.ValueOf(expect)
	right := reflect.ValueOf(actual)

	// Initialize the total score.
	var res float64

	// Calculate the maximum number of keys in the two maps.
	total := max(left.Len(), right.Len())

	// Create a map to keep track of the keys that have been matched.
	marked := make(map[reflect.Value]bool, total)

	// Iterate over the keys of the left map.
	for _, k := range left.MapKeys() {
		// If the corresponding key exists in the right map, calculate the match
		// score between the values and add it to the total score.
		// Mark the key as matched.
		if right.MapIndex(k).IsValid() {
			res += compare(left.MapIndex(k).Interface(), right.MapIndex(k).Interface())
			marked[right.MapIndex(k)] = true
		}
	}

	// Iterate over the keys of the right map.
	// If a key has not been marked as matched, calculate the match score between
	// the corresponding values in the left and right maps and add it to the total
	// score.
	for _, k := range right.MapKeys() {
		if _, ok := marked[k]; ok {
			continue
		}

		if left.MapIndex(k).IsValid() {
			res += compare(left.MapIndex(k).Interface(), right.MapIndex(k).Interface())
		}
	}

	// If the total score is 0 and the maximum number of keys is 0, return 1.
	if res == 0 && total == 0 {
		return 1
	}

	// Return the total score divided by the maximum number of keys.
	return res / float64(total)
}

// slicesRankMatch is a function that calculates the match score between two
// slices or maps. It takes the expected and actual values and a ranker
// function that compares two values and returns a match score between 0 and 1.
//
// The ranker function is called for each pair of values in the slices or maps,
// and the match scores are accumulated. The function returns the accumulated
// match score divided by the maximum number of values in the slices or maps.
//
// If the types of the expected and actual values are not equal, the function
// returns 0. If either the expected or actual value is nil, the function
// returns 1. If the types of the expected and actual values are not slice or
// map, the function returns 0.
//
// The ranker function is called for each pair of values in the slices or maps,
// and the match scores are accumulated. The function returns the accumulated
// match score divided by the maximum number of values in the slices or maps.
//
// The function uses a marked algorithm to avoid redundant comparisons.
//
//nolint:cyclop
func slicesRankMatch(expect, actual any, compare ranker) float64 {
	// Check if the types of the expected and actual values are equal.
	if reflect.TypeOf(expect) != reflect.TypeOf(actual) {
		return 0
	}

	// If both values are nil, return 1.
	if reflect.TypeOf(expect) == nil {
		return 1
	}

	// Convert the expected and actual values to reflect.Value.
	a := reflect.ValueOf(expect)
	b := reflect.ValueOf(actual)

	// If the types of the expected and actual values are not slice, return 0.
	if a.Kind() != reflect.Slice || b.Kind() != reflect.Slice {
		return 0
	}

	var res float64 // Initialize the total score.

	marked := make(map[int]struct{}, b.Len()) // Create a map to keep track of the keys that have been marked as matched.

	// Iterate over the values of the left slice.
	for i := range a.Len() {
		// Iterate over the values of the right slice.
		for j := range b.Len() {
			// Skip the value if it has already been marked as matched.
			if _, ok := marked[j]; ok {
				continue
			}

			// Calculate the match score between the values of the indices and
			// add it to the total score if the result is not 0.
			if result := compare(a.Index(i).Interface(), b.Index(j).Interface()); result != 0 {
				res += result
				marked[j] = struct{}{}
			}
		}
	}

	total := max(a.Len(), b.Len()) // Calculate the maximum number of values in the two slices.

	// If the total score is 0 and the maximum number of values is 0, return 1.
	if res == 0 && total == 0 {
		return 1
	}

	// Return the total score divided by the maximum number of values.
	return res / float64(total)
}

// distance calculates the Levenshtein distance between two strings.
// It returns a float64 representing the distance normalized by the length of the
// longer string.
//
// The Levenshtein distance is a measure of the number of single-character edits
// needed to transform one string into another, such as insertion, deletion, or
// substitution.
//
// Parameters:
// - s: The first string.
// - t: The second string.
func distance(s, t string) float64 {
	// Convert the strings to runes.
	r1, r2 := []rune(s), []rune(t)

	// Create a column of integers to store the edit distances.
	column := make([]int, 1, 64) //nolint:mnd // Initial size is 1, capacity is 64.

	// Populate the column with the row indices (0, 1, 2, ...).
	for y := 1; y <= len(r1); y++ {
		column = append(column, y)
	}

	// Iterate over the columns of the matrix.
	for x := 1; x <= len(r2); x++ {
		// Set the first element of the column to the current column index.
		column[0] = x

		// Iterate over the rows of the matrix.
		for y, lastDiag := 1, x-1; y <= len(r1); y++ {
			// Store the previous diagonal value.
			oldDiag := column[y]

			// Calculate the cost of the current edit operation.
			cost := 0
			if r1[y-1] != r2[x-1] {
				cost = 1
			}

			// Update the current cell of the column.
			column[y] = min(column[y]+1, column[y-1]+1, lastDiag+cost)

			// Store the previous diagonal value for the next iteration.
			lastDiag = oldDiag
		}
	}

	// Calculate the length of the longer string.
	length := float64(max(len(s), len(t)))

	// Return the normalized Levenshtein distance.
	return (length - float64(column[len(r1)])) / length
}
