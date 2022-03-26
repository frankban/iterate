// Licensed under the MIT license, see LICENSE file for details.

package iterate

import "golang.org/x/exp/constraints"

// Count returns an iterator counting consecutive values from start to stop with
// the given step. The returned error is always nil.
//
// For instance:
//
//     counter := it.Count(0, 10, 1)
//     for counter.Next() {
//         v := counter.Value()
//         // v is 0, then 1, then 2 and so on till 9.
//     }
//
func Count(start, stop, step int) Iterator[int] {
	return &counter{
		start: start - step,
		stop:  stop,
		step:  step,
	}
}

type counter struct {
	start, stop, step int
}

// Next implements Iterator[T].Next.
func (it *counter) Next() bool {
	it.start += it.step
	if it.start == it.stop {
		it.start, it.stop = 0, it.step
		return false
	}
	return true
}

// Value implements Iterator[T].Value by returning the next number in the count.
func (it *counter) Value() int {
	return it.start
}

// Err implements Iterator[T].Err. The returned error is always nil.
func (it *counter) Err() error {
	return nil
}

// Enumerate returns an iterator that produces key/value pairs in which the keys
// are iterator indexes (starting from 0) and the values are produced by the
// given iterator.
//
// For instance:
//
//     letters = it.Enumerate(it.FromSlice([]string{"a", "b", "c"}))
//     for letters.Next() {
//         kv := letters.Value()
//         // kv is (0, "a"), then (1, "b"), then (2, "c")
//     }
//
func Enumerate[T any](it Iterator[T]) Iterator[KeyValue[int, T]] {
	return &enumerator[T]{
		source: it,
		idx:    -1,
	}
}

type enumerator[T any] struct {
	source Iterator[T]
	idx    int
}

// Next implements Iterator[T].Next.
func (it *enumerator[T]) Next() bool {
	for it.source.Next() {
		it.idx++
		return true
	}
	it.idx = 0
	return false
}

// Value implements Iterator[T].Next enumerating values from the source iterator.
func (it *enumerator[T]) Value() KeyValue[int, T] {
	return KeyValue[int, T]{
		Key:   it.idx,
		Value: it.source.Value(),
	}
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *enumerator[T]) Err() error {
	return it.source.Err()
}

// Sum sums the values produced by the given iterator.
//
// For instance:
//
//     sum, err := it.Sum(it.Count(1, 11, 1)) // sum is 55
//     sum, err := it.Sum(it.FromSlice([]string{"hello", " ", "world"})) // sum is "hello world"
//
func Sum[T constraints.Ordered](it Iterator[T]) (sum T, err error) {
	return Reduce(it, func(a, v T) T {
		return a + v
	}, *new(T))
}
