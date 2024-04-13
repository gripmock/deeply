package deeply_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gripmock/deeply"
)

func ranker(expect any, actual []any) []any {
	sort.SliceStable(actual, func(i, j int) bool {
		return deeply.RankMatch(expect, actual[i]) > deeply.RankMatch(expect, actual[j])
	})

	return actual
}

func TestRankMatch_Simple(t *testing.T) {
	require.Equal(t, []any{"a", "ab", "abc"}, ranker("a", []any{"a", "ab", "abc"}))
	require.Equal(t, []any{"a", "ab", "abc"}, ranker("a", []any{"a", "abc", "ab"}))

	require.Equal(t, []any{"hello", "world", "zzzzz"}, ranker("hella", []any{"world", "hello", "zzzzz"}))
	require.Equal(t, []any{"hello", "world", "zzzzz"}, ranker("hella", []any{"world", "zzzzz", "hello"}))
	require.Equal(t, []any{"hello", "world", "zzzzz"}, ranker("hella", []any{"hello", "zzzzz", "world"}))

	require.Equal(t,
		[]any{[]int{1, 2, 3}, []int{3, 2, 1}, []int{1}},
		ranker([]int{1, 2, 3}, []any{[]int{1, 2, 3}, []int{3, 2, 1}, []int{1}}))
}

func TestRankMatch_Map_Left(t *testing.T) {
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

	c := map[string]interface{}{
		"c": map[string]interface{}{
			"d": "d",
			"e": []int{1, 2, 3},
			"f": []string{"a", "b", "c"},
		},
		"b":      "b",
		"a":      "a",
		"name":   "gripmock",
		"cities": []string{"Istanbul", "Stalingrad"},
	}

	require.Equal(t, []any{b, c}, ranker(a, []any{c, b}))
}

func TestRankMatch_Map_Right(t *testing.T) {
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

	c := map[string]interface{}{
		"c": map[string]interface{}{
			"d": "d",
			"e": []int{1, 2, 3},
			"f": []string{"a", "b", "c"},
		},
		"b":      "b",
		"a":      "a",
		"name":   "gripmock",
		"cities": []string{"Istanbul", "Stalingrad"},
	}

	require.Equal(t, []any{b, c}, ranker(a, []any{c, b}))
}

func TestRankMatch_Boundary(t *testing.T) {
	require.Equal(t, []any{nil, false, true, 0, 1}, ranker(nil, []any{false, true, 0, 1, nil}))
	require.Equal(t,
		[]any{[]string{"a", "b", "c"}, []string{"a", "b", "d"}, []string{"a", "c", "d"}},
		ranker(
			[]string{"[a]", "[b]", "[cd]"},
			[]any{[]string{"a", "b", "c"}, []string{"a", "b", "d"}, []string{"a", "c", "d"}}))

	require.Greater(t, deeply.RankMatch(map[string]interface{}{
		"field1": "hello field1",
		"field3": "hello field3",
	}, map[string]interface{}{
		"field1": "hello field1",
	}), 0.)

	require.Greater(t, deeply.RankMatch(map[string]interface{}{}, map[string]interface{}{}), 0.)
}
