// Licensed under the MIT license, see LICENSE file for details.

package iterate

// DropWhile returns an iterator discarding values from the given iterator while
// predicate(v) is true.
func DropWhile[T any](it Iterator[T], predicate func(idx int, v T) bool) Iterator[T] {
	return &dropper[T]{
		Iterator:  it,
		predicate: predicate,
	}
}

type dropper[T any] struct {
	Iterator[T]
	predicate func(idx int, v T) bool
	idx       int
	started   bool
}

// Next implements Iterator[T].Next by discarding values while they satisfy the
// predicate.
func (it *dropper[T]) Next() bool {
	for it.Iterator.Next() {
		if !it.started && it.predicate(it.idx, it.Iterator.Value()) {
			it.idx++
			continue
		}
		it.started = true
		return true
	}
	return false
}

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
	stopped   bool
}

// Next implements Iterator[T].Next.
func (it *taker[T]) Next() bool {
	if !it.stopped && it.source.Next() && it.predicate(it.idx, it.source.Value()) {
		it.idx++
		return true
	}
	it.stopped = true
	return false
}

// Value implements Iterator[T].Value by returning values until they satisfy the
// predicate.
func (it *taker[T]) Value() T {
	if it.stopped {
		return *new(T)
	}
	return it.source.Value()
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
	it.Next()
	return it.Value(), it.Err()
}
