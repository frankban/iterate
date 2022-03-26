// Licensed under the MIT license, see LICENSE file for details.

package iterate

// FromSlice returns an iterator producing values from the given slice.
func FromSlice[T any](s []T) Iterator[T] {
	return &sliceIterator[T]{
		s:     s,
		first: true,
	}
}

type sliceIterator[T any] struct {
	s     []T
	first bool
}

// Next implements Iterator[T].Next.
func (it *sliceIterator[T]) Next() bool {
	if len(it.s) == 0 {
		return false
	}
	if it.first {
		it.first = false
	} else {
		it.s = it.s[1:]
	}
	return len(it.s) > 0
}

// Value implements Iterator[T].Value by returning values from a slice.
func (it *sliceIterator[T]) Value() T {
	if len(it.s) == 0 {
		return *new(T)
	}
	return it.s[0]
}

// Err implements Iterator[T].Err. The returned error is always nil.
func (it *sliceIterator[T]) Err() error {
	return nil
}

// ToSlice consumes the given iterator and returns a slice of produced values or
// an error occurred while iterating. This function should not be used with
// infinite iterators (see TakeWhile or Limit).
func ToSlice[T any](it Iterator[T]) (s []T, err error) {
	for it.Next() {
		s = append(s, it.Value())
	}
	return s, it.Err()
}
