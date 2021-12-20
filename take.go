// Licensed under the MIT license, see LICENSE file for details.

package iterate

// TakeWhile returns an iterator producing values from the given iterator while
// predicate(v) is true.
func TakeWhile[T any](it Iterator[T], predicate func(idx int, v T) bool) Iterator[T] {
	return &taker[T]{
		source:    it,
		predicate: predicate,
	}
}

type taker[T any] struct {
	source    Iterator[T]
	predicate func(idx int, v T) bool
	idx       int
}

// Next implements Iterator[T].Next by producing values until they satisfy the
// predicate.
func (it *taker[T]) Next(v *T) bool {
	var val T
	if it.source.Next(&val) && it.predicate(it.idx, val) {
		*v = val
		it.idx++
		return true
	}
	return false
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *taker[T]) Err() error {
	return it.source.Err()
}

// Limit returns an iterator limiting the number of values returned by the given
// iterator.
func Limit[T any](it Iterator[T], limit int) Iterator[T] {
	return TakeWhile(it, func(idx int, v T) bool {
		return idx < limit
	})
}

// Next returns the next value produced by the iterator.
func Next[T any](it Iterator[T]) (T, error) {
	var v T
	it.Next(&v)
	return v, it.Err()
}
