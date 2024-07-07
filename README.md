# deeply

`deeply` is a package that provides stub filtering algorithms, which is used in `bavix/gripmock`.

Stub filtering algorithms is a technique used to filter stubs based on expectations. It is a powerful tool for testing and development.

These algorithms are written in a form of a function that takes an expectation and an actual value and returns a boolean result.

In a nutshell, `Matches` checks if the actual value matches the expectation, `Contains` checks if the actual value contains the expectation, `ContainsIgnoreArrayOrder` checks if the actual value contains the expectation ignoring the order of elements in arrays, and `MatchesIgnoreArrayOrder` checks if the actual value matches the expectation ignoring the order of elements in arrays.

Furthermore, `deeply` provides methods `Equals` and `EqualsIgnoreArrayOrder` that checks if the expected and actual values are deeply equal. The `Equals` method checks if the expected and actual values are deeply equal and ignores the order of arrays, while the `EqualsIgnoreArrayOrder` method checks if the expected and actual values are deeply equal ignoring the order of arrays.

Both `Equals` and `EqualsIgnoreArrayOrder` methods return true if any of the following conditions are met:

- The expected and actual values are both maps and have the same number of keys.
- The expected and actual values are both slices and have the same length.
- The expected and actual values are deeply equal using reflect.DeepEqual.
