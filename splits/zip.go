// https://github.com/jkratz55/slices
// LICENSE: MIT

package splits

// Pair is a type representing a pair of values
type Pair[T, U any] struct {
	First  T
	Second U
}

// Zip accepts two arrays/slices and zip the values together returning a slice of
// Pairs. If the two arrays/slices are not of equal lengths this function will
// panic.
func Zip[T, U any](left []T, right []U) []Pair[T, U] {
	if len(left) != len(right) {
		panic("cannot zip slices of different lengths")
	}
	pairs := make([]Pair[T, U], 0, len(left))
	for idx, item := range left {
		pairs = append(pairs, Pair[T, U]{
			First:  item,
			Second: right[idx],
		})
	}
	return pairs
}
