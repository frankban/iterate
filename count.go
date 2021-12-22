// Licensed under the MIT license, see LICENSE file for details.

package iterate

import "constraints"

// Count returns an iterator counting consecutive values from start to stop with
// the given step. The returned error is always nil.
//
// For instance:
//
//     counter := it.Count(0, 10, 1)
//     var v int
//     for counter.Next(&v) {
//         // v is 0, then 1, then 2 and so on till 9.
//     }
//
func Count(start, stop, step int) Iterator[int] {
	return &counter{
		start: start,
		stop:  stop,
		step:  step,
	}
}

type counter struct {
	start, stop, step int
}

// Next implements Iterator[T].Next by producing the next number in the count.
func (it *counter) Next(v *int) bool {
	if it.start == it.stop {
		*v = 0
		return false
	}
	*v = it.start
	it.start += it.step
	return true
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
//     letters := it.FromSlice([]string{"a", "b", "c"})
//     var kv KeyValue[int, string]
//     for it.Enumerate(letters).Next(&kv) {
//         // kv is (0, "a"), then (1, "b"), then (2, "c")
//     }
//
func Enumerate[T any](it Iterator[T]) Iterator[KeyValue[int, T]] {
	return &enumerator[T]{
		source: it,
	}
}

type enumerator[T any] struct {
	source Iterator[T]
	idx    int
}

// Next implements Iterator[T].Next enumerating values from the source iterator.
func (it *enumerator[T]) Next(v *KeyValue[int, T]) bool {
	var val T
	for it.source.Next(&val) {
		*v = KeyValue[int, T]{
			Key:   it.idx,
			Value: val,
		}
		it.idx++
		return true
	}
	*v = KeyValue[int, T]{}
	return false
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
