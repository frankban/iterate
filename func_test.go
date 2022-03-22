// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"errors"
	"strings"
	"testing"
	"unicode"

	"github.com/go-quicktest/qt"

	it "github.com/frankban/iterate"
)

func TestFilter(t *testing.T) {
	iter := it.FromSlice([]string{"these", "are", "the", "voyages"})
	iter = it.Filter(iter, func(v string) bool {
		return len(v) == 3
	})
	var v string
	var vs []string
	for iter.Next(&v) {
		vs = append(vs, v)
	}
	qt.Assert(t, qt.IsNil(iter.Err()))
	qt.Assert(t, qt.DeepEquals(vs, []string{"are", "the"}))

	// Further calls to next return false and produce the zero value.
	qt.Assert(t, qt.IsFalse(iter.Next(&v)))
	qt.Assert(t, qt.Equals(v, ""))
}

func TestFilterError(t *testing.T) {
	var iter it.Iterator[rune] = &errorIterator[rune]{
		v: 'r',
	}
	iter = it.Filter(iter, unicode.IsLower)
	var v rune
	var vs []rune
	for iter.Next(&v) {
		vs = append(vs, v)
	}
	qt.Assert(t, qt.ErrorMatches(iter.Err(), "bad wolf"))
	qt.Assert(t, qt.DeepEquals(vs, []rune{'r'}))
}

func TestMap(t *testing.T) {
	iter := it.FromSlice([]string{"these", "are", "the", "voyages"})
	iter = it.Map(iter, strings.ToUpper)
	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, []string{"THESE", "ARE", "THE", "VOYAGES"}))

	// Further calls to next return false and produce the zero value.
	v := "42"
	qt.Assert(t, qt.IsFalse(iter.Next(&v)))
	qt.Assert(t, qt.Equals(v, ""))
}

func TestMapDifferentTypes(t *testing.T) {
	type rectangle struct {
		x, y int
	}

	iter := it.FromSlice([]rectangle{{
		x: 1, y: 2,
	}, {
		x: 4, y: 5,
	}, {
		x: 10, y: 20,
	}})
	areas := it.Map(iter, func(v rectangle) int {
		return v.x * v.y
	})
	got, err := it.ToSlice(areas)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, []int{2, 20, 200}))
}

func TestMapError(t *testing.T) {
	iter := it.Chain[int](it.Count(1, 5, 1), &errorIterator[int]{v: 5})
	got, err := it.ToSlice(it.Map(iter, func(v int) int {
		return v * v
	}))
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(got, []int{1, 4, 9, 16, 25}))
}

func TestReduce(t *testing.T) {
	iter := it.FromSlice([]string{"hello", "world"})
	length, err := it.Reduce(iter, func(a int, v string) int {
		return a + len(v)
	}, 0)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(length, 10))
}

func TestReduceError(t *testing.T) {
	iter := it.Chain[int](it.Count(1, 5, 1), &errorIterator[int]{v: 5})
	got, err := it.Reduce(iter, func(a, v int) int {
		return a + v
	}, 0)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(got, 15))
}

func TestChain(t *testing.T) {
	iter := it.Chain(it.FromSlice([]string{"1", "2"}), it.FromSlice([]string{"3", "4"}))
	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, []string{"1", "2", "3", "4"}))

	// Further calls to next return false and produce the zero value.
	v := "42"
	qt.Assert(t, qt.IsFalse(iter.Next(&v)))
	qt.Assert(t, qt.Equals(v, ""))
}

func TestChainError(t *testing.T) {
	iter := it.Chain[int](it.Count(1, 5, 1), &errorIterator[int]{v: 5})
	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(got, []int{1, 2, 3, 4, 5}))
}

func TestRepeat(t *testing.T) {
	iter := it.Limit(it.Repeat(it.Count(3, 0, -1)), 9)
	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, []int{3, 2, 1, 3, 2, 1, 3, 2, 1}))
}

func TestRepeatError(t *testing.T) {
	iter := it.Repeat[string](&errorIterator[string]{v: "v"})
	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(got, []string{"v"}))
}

// errorIterator is an iterator returning the given value and then an error.
type errorIterator[T any] struct {
	v        T
	numcalls int
}

func (it *errorIterator[T]) Next(v *T) bool {
	it.numcalls++
	if it.numcalls > 1 {
		return false
	}
	*v = it.v
	return true
}

func (it *errorIterator[T]) Err() error {
	if it.numcalls > 1 {
		return errors.New("bad wolf")
	}
	return nil
}

func TestTee(t *testing.T) {
	var vs []int
	iter := it.Tee(it.Count(0, 10, 2), func(v int) {
		vs = append(vs, v)
	})

	v, err := it.Next(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(v, 0))
	qt.Assert(t, qt.DeepEquals(vs, []int{0}))

	s, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(s, []int{2, 4, 6, 8}))
	qt.Assert(t, qt.DeepEquals(vs, []int{0, 2, 4, 6, 8}))
}

func TestIteratorFunc(t *testing.T) {
	iter := it.Count(0, 10, 2)
	f := it.IteratorFunc[int](func(v *int) bool {
		return iter.Next(v)
	})
	got, err := it.ToSlice[int](f)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, []int{0, 2, 4, 6, 8}))
}
