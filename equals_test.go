package deeply_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gripmock/deeply"
)

func TestEquals_Simple(t *testing.T) {
	require.True(t, deeply.Equals("a", "a"))
	require.False(t, deeply.Equals("a", "b"))

	require.True(t, deeply.Equals([]int{1, 2, 3}, []int{1, 2, 3}))
}

func TestEquals_Map_Left(t *testing.T) {
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

	require.True(t, deeply.Equals(a, b))

	delete(a, "a")

	require.False(t, deeply.Equals(a, b))

	a["a"] = true

	require.False(t, deeply.Equals(a, b))
}

func TestEquals_Map_Right(t *testing.T) {
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

	require.True(t, deeply.Equals(a, b))

	delete(b, "a")

	require.False(t, deeply.Equals(a, b))

	b["a"] = true

	require.False(t, deeply.Equals(a, b))
}

func TestEquals_Slices_Left(t *testing.T) {
	require.True(t, deeply.Equals([]int{1, 2, 3}, []int{1, 2, 3}))

	require.False(t, deeply.Equals([]int{1, 3, 2}, []int{1, 2, 3}))
	require.False(t, deeply.Equals([]int{1, 2}, []int{1, 2, 3}))

	require.True(t, deeply.Equals([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))

	require.False(t, deeply.Equals([]interface{}{1, 3, 2}, []interface{}{1, 2, 3}))
	require.False(t, deeply.Equals([]interface{}{1, 2}, []interface{}{1, 2, 3}))
}

func TestEquals_Slices_Right(t *testing.T) {
	require.True(t, deeply.Equals([]int{1, 2, 3}, []int{1, 2, 3}))

	require.False(t, deeply.Equals([]int{1, 2, 3}, []int{1, 3, 2}))
	require.False(t, deeply.Equals([]int{1, 2, 3}, []int{1, 2}))

	require.True(t, deeply.Equals([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))

	require.False(t, deeply.Equals([]interface{}{1, 2, 3}, []interface{}{1, 3, 2}))
	require.False(t, deeply.Equals([]interface{}{1, 2, 3}, []interface{}{1, 2}))
}

func TestEquals_Slices_OrderIgnore(t *testing.T) {
	require.False(t, deeply.EqualsIgnoreArrayOrder([]string{"a", "b", "c"}, []int{1, 2, 3}))

	require.True(t, deeply.EqualsIgnoreArrayOrder([]string{"a", "b", "c"}, []string{"b", "a", "c"}))

	require.True(t, deeply.EqualsIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 2, 3}))
	require.True(t, deeply.EqualsIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 3, 2}))
	require.True(t, deeply.EqualsIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))
	require.True(t, deeply.EqualsIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 3, 2}))

	require.False(t, deeply.EqualsIgnoreArrayOrder([]int{1, 2, 3}, []int{1, 2}))
	require.False(t, deeply.EqualsIgnoreArrayOrder([]interface{}{1, 2, 3}, []interface{}{1, 2}))
}

func TestEquals_Boundary(t *testing.T) {
	require.False(t, deeply.Equals([]string{"a", "a", "a"}, []string{"a", "b", "c"}))
	require.False(t, deeply.Equals([]string{"a", "b", "c"}, []string{"a", "a", "a"}))
	require.False(t, deeply.Equals(nil, false))

	require.True(t, deeply.Equals(nil, nil))

	require.True(t, deeply.Equals(map[string]interface{}{
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
