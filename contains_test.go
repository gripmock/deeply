package deeply_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gripmock/deeply"
)

func TestContains_Simple(t *testing.T) {
	require.True(t, deeply.Contains("a", "a"))
	require.False(t, deeply.Contains("a", "b"))

	require.True(t, deeply.Contains([]int{1, 2, 3}, []int{1, 2, 3}))
	require.False(t, deeply.Contains([]int{1, 2, 3}, []int{1, 3, 2}))
}

func TestContains_Map_Left(t *testing.T) {
	a := map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": map[string]interface{}{
			"f": []string{"a", "b", "c"},
			"d": "d",
			"e": []int{1, 2, 3},
		},
	}

	b := map[string]interface{}{
		"c": map[string]interface{}{
			"d": "d",
			"e": []int{1, 2, 3},
			"f": []string{"a", "b", "c"},
		},
		"b": "b",
		"a": "a",
	}

	require.True(t, deeply.Contains(a, b))

	delete(a, "a")

	require.True(t, deeply.Contains(a, b))

	a["a"] = true

	require.False(t, deeply.Contains(a, b))
}

func TestContains_Map_Right(t *testing.T) {
	a := map[string]interface{}{
		"a": "a",
		"b": "b",
		"c": map[string]interface{}{
			"f": []string{"a", "b", "c"},
			"d": "d",
			"e": []int{1, 2, 3},
		},
	}

	b := map[string]interface{}{
		"c": map[string]interface{}{
			"d": "d",
			"e": []int{1, 2, 3},
			"f": []string{"a", "b", "c"},
		},
		"b": "b",
		"a": "a",
	}

	require.True(t, deeply.Contains(a, b))

	delete(b, "a")

	require.False(t, deeply.Contains(a, b))

	b["a"] = true

	require.False(t, deeply.Contains(a, b))
}

func TestContains_Slices_Left(t *testing.T) {
	require.True(t, deeply.Contains([]int{1, 2, 3}, []int{1, 2, 3}))

	require.False(t, deeply.Contains([]int{1, 3, 2}, []int{1, 2, 3}))
	require.False(t, deeply.Contains([]int{1, 2}, []int{1, 2, 3}))

	require.True(t, deeply.Contains([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))

	require.False(t, deeply.Contains([]interface{}{1, 3, 2}, []interface{}{1, 2, 3}))
	require.False(t, deeply.Contains([]interface{}{1, 2}, []interface{}{1, 2, 3}))
}

func TestContains_Slices_Right(t *testing.T) {
	require.True(t, deeply.Contains([]int{1, 2, 3}, []int{1, 2, 3}))

	require.False(t, deeply.Contains([]int{1, 2, 3}, []int{1, 3, 2}))
	require.False(t, deeply.Contains([]int{1, 2, 3}, []int{1, 2}))

	require.True(t, deeply.Contains([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))

	require.False(t, deeply.Contains([]interface{}{1, 2, 3}, []interface{}{1, 3, 2}))
	require.False(t, deeply.Contains([]interface{}{1, 2, 3}, []interface{}{1, 2}))
}

func TestContains_MapStable(t *testing.T) {
	a := map[string][]interface{}{
		"items": {
			map[string]interface{}{
				"high": json.Number("72057594037927936"),
				"low":  json.Number("18446744073709551615"),
			},
			map[string]interface{}{
				"low":  json.Number("2"),
				"high": json.Number("1"),
			},
		},
	}

	b := map[string][]interface{}{
		"items": {
			map[string]interface{}{
				"low":  json.Number("18446744073709551615"),
				"high": json.Number("72057594037927936"),
			},
			map[string]interface{}{
				"high": json.Number("1"),
				"low":  json.Number("2"),
			},
		},
	}

	require.True(t, deeply.Contains(a, b))
	require.True(t, deeply.Contains(b, a))

	a["items"][0], a["items"][1] = a["items"][1], a["items"][0]

	require.False(t, deeply.Contains(a, b))
	require.False(t, deeply.Contains(b, a))

	require.True(t, deeply.ContainsIgnoreArrayOrder(a, b))
	require.True(t, deeply.ContainsIgnoreArrayOrder(b, a))

	require.False(t, deeply.ContainsIgnoreArrayOrder(
		[]string{"a", "a", "a"}, []string{"a", "b", "c", "a"}))
	require.False(t, deeply.ContainsIgnoreArrayOrder(
		[]string{"a", "a", "a"}, []string{"a", "b", "c"}))

	require.False(t, deeply.Contains([]string{"a", "c", "b"}, []string{"a", "b", "c"}))

	require.True(t, deeply.ContainsIgnoreArrayOrder([]string{"a", "c", "b"}, []string{"a", "b", "c"}))
}

func TestContains_Slices_OrderIgnore(t *testing.T) {
	require.False(t, deeply.ContainsIgnoreArrayOrder([]string{"a", "b", "c"}, []int{1, 2, 3}))

	require.True(t, deeply.ContainsIgnoreArrayOrder([]string{"a", "b", "c"}, []string{"b", "a", "c"}))

	require.True(t, deeply.ContainsIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 2, 3}))
	require.True(t, deeply.ContainsIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 3, 2}))
	require.True(t, deeply.ContainsIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))
	require.True(t, deeply.ContainsIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 3, 2}))

	require.False(t, deeply.ContainsIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 2}))
	require.False(t, deeply.ContainsIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 2}))
}

func TestContains_Boundary(t *testing.T) {
	require.False(t, deeply.Contains([]string{"a", "a", "a"}, []string{"a", "b", "c"}))
	require.False(t, deeply.Contains([]string{"a", "b", "c"}, []string{"a", "a", "a"}))
	require.False(t, deeply.Contains(nil, false))

	require.True(t, deeply.Contains(nil, nil))

	require.False(t, deeply.Contains(map[string]interface{}{
		"field1": "hello",
	}, map[string]interface{}{
		"field2": "hello field1",
	}))

	require.True(t, deeply.Contains(map[string]interface{}{
		"name": "Afra Gokce",
		"age":  1,
		"girl": true,
		"null": nil,
		"greetings": map[string]interface{}{
			"hola":    "mundo",
			"merhaba": "dunya",
		},
		"cities": []interface{}{
			"Istanbul",
			"Jakarta",
		},
	}, map[string]interface{}{
		"name": "Afra Gokce",
		"age":  1,
		"girl": true,
		"null": nil,
		"greetings": map[string]interface{}{
			"hola":    "mundo",
			"merhaba": "dunya",
		},
		"cities": []interface{}{
			"Istanbul",
			"Jakarta",
		},
	}))
}
