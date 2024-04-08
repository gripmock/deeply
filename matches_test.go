package deeply_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gripmock/deeply"
)

func TestMatches_Simple(t *testing.T) {
	require.True(t, deeply.Matches("a", "a"))
	require.False(t, deeply.Matches("a", "b"))

	require.True(t, deeply.Matches([]int{1, 2, 3}, []int{1, 2, 3}))
	require.False(t, deeply.Matches([]int{1, 2, 3}, []int{1, 3, 2}))
}

func TestMatches_Map_Left(t *testing.T) {
	a := map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": map[string]interface{}{
			"f": []string{"a", "b", "c"},
			"d": "d",
			"e": []int{1, 2, 3},
		},
		"name":   "^grip.*$",
		"cities": []string{"Jakarta", "Istanbul", ".*grad$"},
	}

	b := map[string]interface{}{
		"c": map[string]interface{}{
			"d": "d",
			"e": []int{1, 2, 3},
			"f": []string{"a", "b", "c"},
		},
		"b":      "b",
		"a":      "a",
		"name":   "gripmock",
		"cities": []string{"Jakarta", "Istanbul", "Stalingrad"},
	}

	require.True(t, deeply.Matches(a, b))

	delete(a, "a")

	require.True(t, deeply.Matches(a, b))

	a["a"] = true

	require.False(t, deeply.Matches(a, b))
}

func TestMatches_Map_Right(t *testing.T) {
	a := map[string]interface{}{
		"a": "[a-z]",
		"b": "b",
		"c": map[string]interface{}{
			"f": []string{"[a-z]", "[0-9]", "c"},
			"d": "d",
			"e": []int{1, 2, 3},
		},
	}

	b := map[string]interface{}{
		"c": map[string]interface{}{
			"d": "d",
			"e": []int{1, 2, 3},
			"f": []string{"d", "1", "c"},
		},
		"b": "b",
		"a": "c",
	}

	require.True(t, deeply.Matches(a, b))

	delete(b, "a")

	require.False(t, deeply.Matches(a, b))

	b["a"] = true

	require.False(t, deeply.Matches(a, b))
}

func TestMatches_Slices_Left(t *testing.T) {
	require.True(t, deeply.Matches([]int{1, 2, 3}, []int{1, 2, 3}))

	require.False(t, deeply.Matches([]int{1, 3, 2}, []int{1, 2, 3}))
	require.True(t, deeply.Matches([]int{1, 2}, []int{1, 2, 3}))

	require.True(t, deeply.Matches([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))

	require.False(t, deeply.Matches([]interface{}{1, 3, 2}, []interface{}{1, 2, 3}))
	require.True(t, deeply.Matches([]interface{}{1, 2}, []interface{}{1, 2, 3}))
}

func TestMatches_Slices_Right(t *testing.T) {
	require.False(t, deeply.Matches([]string{"^hello$"}, []string{"hell!"}))

	require.True(t, deeply.Matches([]int{1, 2, 3}, []int{1, 2, 3}))

	require.False(t, deeply.Matches([]int{1, 2, 3}, []int{1, 3, 2}))
	require.False(t, deeply.Matches([]int{1, 2, 3}, []int{1, 2}))

	require.True(t, deeply.Matches([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))

	require.False(t, deeply.Matches([]interface{}{1, 2, 3}, []interface{}{1, 3, 2}))
	require.False(t, deeply.Matches([]interface{}{1, 2, 3}, []interface{}{1, 2}))
}

func TestMatches_Slices_OrderIgnore(t *testing.T) {
	require.False(t, deeply.MatchesIgnoreArrayOrder([]string{"a", "b", "c"}, []int{1, 2, 3}))

	require.True(t, deeply.MatchesIgnoreArrayOrder([]string{"a", "b", "c"}, []string{"b", "a", "c"}))

	require.True(t, deeply.MatchesIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 2, 3}))
	require.True(t, deeply.MatchesIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 3, 2}))
	require.True(t, deeply.MatchesIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))
	require.True(t, deeply.MatchesIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 3, 2}))

	require.False(t, deeply.MatchesIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 2}))
	require.False(t, deeply.MatchesIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 2}))
}

func TestMatches_Boundary(t *testing.T) {
	require.False(t, deeply.Matches([]string{"a", "a", "a"}, []string{"a", "b", "c"}))
	require.False(t, deeply.Matches([]string{"a", "b", "c"}, []string{"a", "a", "a"}))
}
