// Licensed under the MIT license, see LICENSE file for details.

package iterate

// GroupBy returns an iterator returning key/value pairs in which the key is the
// key used for grouping elements using the given function, and the value is an
// iterator of values with that key.
//
// For instance:
//
//     words := it.FromSlice([]string{"a", "be", "it", "no", "hello", "the", "are"})
//     // Group words by length.
// 	   groups := it.GroupBy(words, func(v string) int {
// 	       return len(v)
//     })
//     for groups.Next() {
//         kv := groups.Value()
//	       // kv is (1, Iterator("a")), then (2, Iterator("be", "it", "no")),
//         // then (5, Iterator("hello")), then (3, Iterator("the", "are"))
//     }
//
func GroupBy[T any, K comparable](it Iterator[T], f func(v T) K) Iterator[KeyValue[K, Iterator[T]]] {
	return &grouper[T, K]{
		source:        it,
		f:             f,
		pendingValues: make(map[int][]T),
	}
}

type grouper[T any, K comparable] struct {
	source        Iterator[T]
	f             func(v T) K
	id            int
	key           K
	pendingValues map[int][]T
	iter          Iterator[T]
}

// Next implements Iterator[T].Next by iterating over groups.
func (it *grouper[T, K]) Next() bool {
	if !it.next() {
		it.iter = nil
		it.key = *new(K)
		return false
	}
	if it.iter != nil {
		return true
	}
	return it.Next()
}

func (it *grouper[T, K]) next() bool {
	if !it.source.Next() {
		// The iteration is done.
		return false
	}

	val := it.source.Value()
	key := it.f(val)
	if it.id == 0 || key != it.key {
		// We either are at the beginning of the iteration, or the key just
		// changed. Store a new iterator to be returned on the next call to
		// Value.
		it.id++
		it.iter = &groupKeyIterator[T, K]{
			source: it,
			id:     it.id,
		}
		it.key = key
	}

	// Store the produced value waiting for group iterators to retrieve it.
	it.pendingValues[it.id] = append(it.pendingValues[it.id], val)
	return true
}

// Value implements Iterator[T].Value by iterating over groups.
func (it *grouper[T, K]) Value() KeyValue[K, Iterator[T]] {
	return KeyValue[K, Iterator[T]]{
		Key:   it.key,
		Value: it.iter,
	}
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *grouper[T, K]) Err() error {
	return it.source.Err()
}

// groupKeyIterator is the iterator returned for generating values for a
// specific group key.
type groupKeyIterator[T any, K comparable] struct {
	source *grouper[T, K]
	id     int
	value  T
}

// Next implements Iterator[T].Next.
func (it *groupKeyIterator[T, K]) Next() bool {
	for {
		// Check whether there are pending values already.
		if len(it.source.pendingValues[it.id]) > 0 {
			it.value, it.source.pendingValues[it.id] = it.source.pendingValues[it.id][0], it.source.pendingValues[it.id][1:]
			return true
		}

		// Check whether the grouper is still iterating over this id, in which case we
		// can progress the iteration further and retry.
		if it.source.id == it.id && it.source.next() {
			continue
		}

		// The iteration for this group is over.
		it.value = *new(T)
		return false
	}
}

// Value implements Iterator[T].Value by returning values for a specific key.
func (it *groupKeyIterator[T, K]) Value() T {
	return it.value
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *groupKeyIterator[T, K]) Err() error {
	return it.source.Err()
}
