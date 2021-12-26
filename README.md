[![GoDoc](https://godoc.org/github.com/frankban/iterate?status.svg)](https://godoc.org/github.com/frankban/iterate)
[![Build Status](https://github.com/frankban/iterate/actions/workflows/ci.yaml/badge.svg)](https://github.com/frankban/iterate/actions/workflows/ci.yaml)

# iterate

`go get github.com/frankban/iterate`

Some experiments with generics and lazy evaluation in Go.
Mostly for fun and figuring our how such an API would feel.

## Usage

All functions work with implementations of the Iterator interface.

#### type Iterator

```go
type Iterator[T any] interface {
	// Next produces the next iterator value and assign it to the variable
	// pointed by v. When the iterator is done, no values are assigned and false
	// is returned. Further calls to Next should just return false with no other
	// side effects. When iterating produces an error, false is returned, and
	// Err() returns the error.
	Next(v *T) bool
	// Err returns the first error occurred while iterating.
	Err() error
}
```
Iterator is implemented by types producing values of type T. Implementations are
typically used in for loops, for instance:
```go
    var v SomeType
    for iterator.Next(&v) {
        // Do something with v.
    }
    if err := iterator.Err(); err != nil {
        // Handle error.
    }
```
Depending on the implementation, producing values might lead to errors. For this
reason it is important to always check Err() after iterating.nnnngg

#### func Next

```go
func Next[T any](it Iterator[T]) (T, error)
```
Next returns the next value produced by the iterator.

#### func Filter

```go
func Filter[T any](it Iterator[T], predicate func(v T) bool) Iterator[T]
```
Filter returns an iterator producing values from the given interator, for which
predicate(v) is true.

#### func Map

```go
func Map[S, D any](source Iterator[S], f func(v S) D) Iterator[D]
```
Map returns an iterator that computes the given function using values from the
given iterator.

#### func Reduce

```go
func Reduce[T, A any](it Iterator[T], f func(a A, v T) A, initial A) (A, error)
```
Reduce applies f cumulatively to the values of the given iterator, from left to
right, so as to reduce the iterable to a single value. The first argument is the
accumulated value and the second argument is the value from the iterator.

For instance, for calculating the overall length of a slice of strings:
```go
        iter := it.FromSlice([]string{"hello", "world"})
        length, err := it.Reduce(iter, func(a int, v string) int {
    	       return a + len(v)
        }, 0)
```

#### func Chain

```go
func Chain[T any](base Iterator[T], others ...Iterator[T]) Iterator[T]
```
Chain returns an iterator producing values from the concatenation of all the
given iterators. The iteration is stopped when all iterators are consumed or
when any of them has an error.

#### func Repeat

```go
func Repeat[T any](it Iterator[T]) Iterator[T]
```
Repeat returns an iterator repeating values from the given iterator endlessly.

#### func Tee

```go
func Tee[T any](it Iterator[T], f func(v T)) Iterator[T]
```
Tee returns an iterator that causes the given function to be called each time a
value is produced.

#### func DropWhile

```go
func DropWhile[T any](it Iterator[T], predicate func(idx int, v T) bool) Iterator[T]
```
DropWhile returns an iterator discarding values from the given iterator while
predicate(v) is true.

#### func TakeWhile

```go
func TakeWhile[T any](it Iterator[T], predicate func(idx int, v T) bool) Iterator[T]
```
TakeWhile returns an iterator producing values from the given iterator while
predicate(v) is true.

#### func Limit

```go
func Limit[T any](it Iterator[T], limit int) Iterator[T]
```
Limit returns an iterator limiting the number of values returned by the given
iterator.

#### func GroupBy

```go
func GroupBy[T any, K comparable](it Iterator[T], f func(v T) K) Iterator[KeyValue[K, Iterator[T]]]
```
GroupBy returns an iterator returning key/value pairs in which the key is the
key used for grouping elements using the given function, and the value is an
iterator of values with that key.

For instance:
```go
        words := it.FromSlice([]string{"a", "be", "it", "no", "hello", "the", "are"})
        // Group words by length.
    	   groups := it.GroupBy(words, func(v string) int {
    	       return len(v)
        })
        var kv it.KeyValue[int, it.Iterator[string]]
        for groups.Next(&kv) {
    	       // kv is (1, Iterator("a")), then (2, Iterator("be", "it", "no")),
            // then (5, Iterator("hello")), then (3, Iterator("the", "are"))
        }
```

#### func FromSlice

```go
func FromSlice[T any](s []T) Iterator[T]
```
FromSlice returns an iterator producing values from the given slice.

#### func ToSlice

```go
func ToSlice[T any](it Iterator[T]) (s []T, err error)
```
ToSlice consumes the given iterator and returns a slice of produced values or an
error occurred while iterating. This function should not be used with infinite
iterators (see TakeWhile or Limit).

#### func FromChannel

```go
func FromChannel[T any](ch <-chan T) Iterator[T]
```
FromChannel returns an iterator producing values from the given channel. The
iteration stops when the channel is closed. The returned error is always nil.

#### func Count

```go
func Count(start, stop, step int) Iterator[int]
```
Count returns an iterator counting consecutive values from start to stop with
the given step. The returned error is always nil.

For instance:
```go
    counter := it.Count(0, 10, 1)
    var v int
    for counter.Next(&v) {
        // v is 0, then 1, then 2 and so on till 9.
    }
```

#### func Enumerate

```go
func Enumerate[T any](it Iterator[T]) Iterator[KeyValue[int, T]]
```
Enumerate returns an iterator that produces key/value pairs in which the keys
are iterator indexes (starting from 0) and the values are produced by the given
iterator.

For instance:
```
    letters := it.FromSlice([]string{"a", "b", "c"})
    var kv KeyValue[int, string]
    for it.Enumerate(letters).Next(&kv) {
        // kv is (0, "a"), then (1, "b"), then (2, "c")
    }
```

#### func Sum

```go
func Sum[T constraints.Ordered](it Iterator[T]) (sum T, err error)
```
Sum sums the values produced by the given iterator.

For instance:
```go
    sum, err := it.Sum(it.Count(1, 11, 1)) // sum is 55
    sum, err := it.Sum(it.FromSlice([]string{"hello", " ", "world"})) // sum is "hello world"
```

#### func Lines

```go
func Lines(r io.Reader) Iterator[string]
```
Lines returns an iterator producing lines from the given reader.

#### func Bytes

```go
func Bytes(r io.Reader) Iterator[byte]
```
Bytes returns an iterator producing bytes from the given reader.

#### func ToMap

```go
func ToMap[K comparable, V any](it Iterator[KeyValue[K, V]]) (map[K]V, error)
```
ToMap returns a map with the values produced by the given key/value iterator. An
error is returned if the iterator returns an error, in which case the returned
map includes the key/value pairs already consumed.

#### type KeyValue

```go
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}
```

KeyValue represents a key value pair.

#### func (kv KeyValue[K, V]) Split

```go
func (kv KeyValue[K, V]) Split() (K, V)
```
Split return the key and the value.

#### func Zip

```go
func Zip[K comparable, V any](keys Iterator[K], values Iterator[V]) Iterator[KeyValue[K, V]]
```
Zip returns an iterator of key/value pairs produced by the given key/value
iterators. The shorter of the two iterators is used. An error is returned if any
of the two iterators returns an error.

#### func Unzip

```go
func Unzip[K comparable, V any](kvs Iterator[KeyValue[K, V]]) (Iterator[K], Iterator[V])
```
Unzip returns a key iterator and a value iterator with pairs produced by the
given key/value iterator.

#### type IteratorFunc

```go
type IteratorFunc[T any] func(v *T) bool
```

IteratorFunc is a fun type that implements Iterator. The resulting iterator
Err() is always nil.

#### func (f IteratorFunc[T]) Err

```go
func (f IteratorFunc[T]) Err() error
```
Err implements Iterator[T].Err. The returned error is always nil.

#### func (f IteratorFunc[T]) Next

```go
func (f IteratorFunc[T]) Next(v *T) bool
```
Next implements Iterator[T].Next by calling the func.
