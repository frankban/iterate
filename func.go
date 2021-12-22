// Licensed under the MIT license, see LICENSE file for details.

package iterate

// Filter returns an iterator producing values from the given interator, for
// which predicate(v) is true.
func Filter[T any](it Iterator[T], predicate func(v T) bool) Iterator[T] {
	return &filter[T]{
		source:    it,
		predicate: predicate,
	}
}

type filter[T any] struct {
	source    Iterator[T]
	predicate func(v T) bool
}

// Next implements Iterator[T].Next by producing the next value satisfying the
// predicate.
func (it *filter[T]) Next(v *T) bool {
	var val T
	for it.source.Next(&val) {
		if it.predicate(val) {
			*v = val
			return true
		}
	}
	*v = val
	return false
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *filter[T]) Err() error {
	return it.source.Err()
}

// Map returns an iterator that computes the given function using values from
// the given iterator.
func Map[S, D any](source Iterator[S], f func(v S) D) Iterator[D] {
	return &mapper[S, D]{
		source: source,
		f:      f,
	}
}

type mapper[S, D any] struct {
	source Iterator[S]
	f      func(v S) D
}

// Next implements Iterator[T].Next by producing the next corresponding value.
func (it *mapper[S, D]) Next(v *D) bool {
	var s S
	if it.source.Next(&s) {
		*v = it.f(s)
		return true
	}
	*v = *new(D)
	return false
}

// Err implements Iterator[T].Err by propagating the error from the source
// iterator.
func (it *mapper[S, D]) Err() error {
	return it.source.Err()
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
	var v T
	for it.Next(&v) {
		initial = f(initial, v)
	}
	return initial, it.Err()
}

// Chain returns an iterator producing values from the concatenation of all the
// given iterators. The iteration is stopped when all iterators are consumed or
// when any of them has an error.
func Chain[T any](base Iterator[T], others ...Iterator[T]) Iterator[T] {
	return &chain[T]{
		base:   base,
		others: others,
	}
}

type chain[T any] struct {
	base   Iterator[T]
	others []Iterator[T]
}

// Next implements Iterator[T].Next by producing values until all iterators are
// consumed.
func (it *chain[T]) Next(v *T) bool {
	if it.base.Next(v) {
		return true
	}
	if len(it.others) == 0 || it.base.Err() != nil {
		*v = *(new(T))
		return false
	}
	it.base, it.others = it.others[0], it.others[1:]
	return it.Next(v)
}

// Err implements Iterator[T].Err by propagating the errors from the iterators.
func (it *chain[T]) Err() error {
	return it.base.Err()
}

// Repeat returns an iterator repeating values from the given iterator
// endlessly.
func Repeat[T any](it Iterator[T]) Iterator[T] {
	return &repeater[T]{
		source: it,
	}
}

type repeater[T any] struct {
	source Iterator[T]
	values []T
	idx    int
}

// Next implements Iterator[T].Next by repeating values endlessly.
func (it *repeater[T]) Next(v *T) bool {
	if it.source.Next(v) {
		it.values = append(it.values, *v)
		return true
	}
	if len(it.values) == 0 || it.source.Err() != nil {
		return false
	}
	*v = it.values[it.idx]
	if it.idx+1 == len(it.values) {
		it.idx = 0
	} else {
		it.idx++
	}
	return true
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

// IteratorFunc is a fun type that implements Iterator.
// The resulting iterator Err() is always nil.
type IteratorFunc[T any] func(v *T) bool

// Next implements Iterator[T].Next by calling the func.
func (f IteratorFunc[T]) Next(v *T) bool {
	return f(v)
}

// Err implements Iterator[T].Err. The returned error is always nil.
func (f IteratorFunc[T]) Err() error {
	return nil
}
