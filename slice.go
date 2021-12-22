// Licensed under the MIT license, see LICENSE file for details.

package iterate

// FromSlice returns an iterator producing values from the given slice.
func FromSlice[T any](s []T) Iterator[T] {
	return &sliceIterator[T]{
		s: s,
	}
}

type sliceIterator[T any] struct {
	s   []T
	idx int
}

// Next implements Iterator[T].Next by producing values from a slice.
func (it *sliceIterator[T]) Next(v *T) bool {
	if it.idx == len(it.s) {
		*v = *new(T)
		return false
	}
	*v = it.s[it.idx]
	it.idx++
	return true
}

// Err implements Iterator[T].Err. The returned error is always nil.
func (it *sliceIterator[T]) Err() error {
	return nil
}

// ToSlice consumes the given iterator and returns a slice of produced values or
// an error occurred while iterating. This function should not be used with
// infinite iterators (see TakeWhile or Limit).
func ToSlice[T any](it Iterator[T]) (s []T, err error) {
	var v T
	for it.Next(&v) {
		s = append(s, v)
	}
	return s, it.Err()
}
