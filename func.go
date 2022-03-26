// Licensed under the MIT license, see LICENSE file for details.

package iterate

// Filter returns an iterator producing values from the given interator, for
// which predicate(v) is true.
func Filter[T any](it Iterator[T], predicate func(v T) bool) Iterator[T] {
	return &filter[T]{
		Iterator:  it,
		predicate: predicate,
	}
}

type filter[T any] struct {
	Iterator[T]
	predicate func(v T) bool
}

// Next implements Iterator[T].Next by producing the next value satisfying the
// predicate.
func (it *filter[T]) Next() bool {
	for it.Iterator.Next() {
		if it.predicate(it.Iterator.Value()) {
			return true
		}
	}
	return false
}

// Map returns an iterator that computes the given function using values from
// the given iterator.
func Map[S, D any](source Iterator[S], f func(v S) D) Iterator[D] {
	return &mapper[S, D]{
		Iterator: source,
		f:        f,
	}
}

type mapper[S, D any] struct {
	Iterator[S]
	f func(v S) D
}

// Value implements Iterator[T].Value by returning the next corresponding value.
func (it *mapper[S, D]) Value() D {
	return it.f(it.Iterator.Value())
}

// Reduce applies f cumulatively to the values of the given iterator, from left
// to right, so as to reduce the iterable to a single value. The first argument
// is the accumulated value and the second argument is the value from the
// iterator.
//
// For instance, for calculating the overall length of a slice of strings:
//
//     iter := it.FromSlice([]string{"hello", "world"})
//     length, err := it.Reduce(iter, func(a int, v string) int {
// 	       return a + len(v)
//     }, 0)
//
func Reduce[T, A any](it Iterator[T], f func(a A, v T) A, initial A) (A, error) {
	for it.Next() {
		initial = f(initial, it.Value())
	}
	return initial, it.Err()
}

// Chain returns an iterator producing values from the concatenation of all the
// given iterators. The iteration is stopped when all iterators are consumed or
// when any of them has an error.
func Chain[T any](base Iterator[T], others ...Iterator[T]) Iterator[T] {
	return &chain[T]{
		Iterator: base,
		others:   others,
	}
}

type chain[T any] struct {
	Iterator[T]
	others []Iterator[T]
}

// Next implements Iterator[T].Next by producing values until all iterators
// are consumed.
func (it *chain[T]) Next() bool {
	if it.Iterator.Next() {
		return true
	}
	if len(it.others) == 0 || it.Iterator.Err() != nil {
		return false
	}
	it.Iterator, it.others = it.others[0], it.others[1:]
	return it.Next()
}

// Repeat returns an iterator repeating values from the given iterator
// endlessly.
func Repeat[T any](it Iterator[T]) Iterator[T] {
	return &repeater[T]{
		source: it,
		idx:    -1,
	}
}

type repeater[T any] struct {
	source Iterator[T]
	idx    int
	values []T
}

// Next implements Iterator[T].Next.
func (it *repeater[T]) Next() bool {
	if it.source.Next() {
		it.values = append(it.values, it.source.Value())
		it.idx++
		return true
	}
	if len(it.values) == 0 || it.source.Err() != nil {
		return false
	}
	if it.idx+1 == len(it.values) {
		it.idx = 0
	} else {
		it.idx++
	}
	return true
}

// Value implements Iterator[T].Value by returning values endlessly.
func (it *repeater[T]) Value() T {
	return it.values[it.idx]
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *repeater[T]) Err() error {
	return it.source.Err()
}

// Tee returns an iterator that causes the given function to be called each time
// a value is produced.
func Tee[T any](it Iterator[T], f func(v T)) Iterator[T] {
	return Map(it, func(v T) T {
		f(v)
		return v
	})
}
