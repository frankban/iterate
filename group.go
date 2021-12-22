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
//     var kv it.KeyValue[int, it.Iterator[string]]
//     for groups.Next(&kv) {
//	       // kv is (1, Iterator("a")), then (2, Iterator("be", "it", "no")),
//         // then (5, Iterator("hello")), then (3, Iterator("the", "are"))
//     }
//
func GroupBy[T any, K comparable](it Iterator[T], f func(v T) K) Iterator[KeyValue[K, Iterator[T]]] {
	return &grouper[T, K]{
		source: it,
		f:      f,
	}
}

type grouper[T any, K comparable] struct {
	source          Iterator[T]
	f               func(v T) K
	id              int
	key             K
	pendingValues   map[int][]T
	pendingIterator Iterator[T]
}

// Next implements Iterator[T].Next by iterating over groups.
func (it *grouper[T, K]) Next(v *KeyValue[K, Iterator[T]]) bool {
	if !it.next() {
		*v = KeyValue[K, Iterator[T]]{}
		return false
	}
	if it.pendingIterator != nil {
		*v = KeyValue[K, Iterator[T]]{
			Key:   it.key,
			Value: it.pendingIterator,
		}
		it.pendingIterator = nil
		return true
	}
	return it.Next(v)
}

func (it *grouper[T, K]) next() bool {
	var val T
	if !it.source.Next(&val) {
		// The iteration is done.
		return false
	}

	key := it.f(val)
	if it.id == 0 || key != it.key {
		// We either are at the beginning of the iteration, or the key just
		// changed. Store a new iterator to be produced on the next call to
		// Next.
		it.id++
		id := it.id
		it.pendingIterator = IteratorFunc[T](func(v *T) bool {
			for {
				// Check whether there are pending values already.
				if len(it.pendingValues[id]) > 0 {
					*v, it.pendingValues[id] = it.pendingValues[id][0], it.pendingValues[id][1:]
					return true
				}

				// Check whether the grouper is still iterating over this id, in which case we
				// can progress the iteration further and retry.
				if it.id == id && it.next() {
					continue
				}

				// The iteration for this group is over.
				*v = *new(T)
				return false
			}
		})
	}
	it.key = key

	// Store the produced value waiting for group iterators to retrieve it.
	if it.pendingValues == nil {
		it.pendingValues = make(map[int][]T)
	}
	it.pendingValues[it.id] = append(it.pendingValues[it.id], val)
	return true
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *grouper[T, K]) Err() error {
	return it.source.Err()
}
