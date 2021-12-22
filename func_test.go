// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"errors"
	"strings"
	"testing"
	"unicode"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestFilter(t *testing.T) {
	c := qt.New(t)

	iter := it.FromSlice([]string{"these", "are", "the", "voyages"})
	iter = it.Filter(iter, func(v string) bool {
		return len(v) == 3
	})
	var v string
	var vs []string
	for iter.Next(&v) {
		vs = append(vs, v)
	}
	c.Assert(iter.Err(), qt.IsNil)
	c.Assert(vs, qt.DeepEquals, []string{"are", "the"})

	// Further calls to next return false and produce the zero value.
	c.Assert(iter.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")
}

func TestFilterError(t *testing.T) {
	c := qt.New(t)

	var iter it.Iterator[rune] = &errorIterator[rune]{
		v: 'r',
	}
	iter = it.Filter(iter, unicode.IsLower)
	var v rune
	var vs []rune
	for iter.Next(&v) {
		vs = append(vs, v)
	}
	c.Assert(iter.Err(), qt.ErrorMatches, "bad wolf")
	c.Assert(vs, qt.DeepEquals, []rune{'r'})
}

func TestMap(t *testing.T) {
	c := qt.New(t)

	iter := it.FromSlice([]string{"these", "are", "the", "voyages"})
	iter = it.Map(iter, strings.ToUpper)
	got, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.DeepEquals, []string{"THESE", "ARE", "THE", "VOYAGES"})

	// Further calls to next return false and produce the zero value.
	v := "42"
	c.Assert(iter.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")
}

func TestMapDifferentTypes(t *testing.T) {
	c := qt.New(t)

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
	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.DeepEquals, []int{2, 20, 200})
}

func TestMapError(t *testing.T) {
	iter := it.Chain[int](it.Count(1, 5, 1), &errorIterator[int]{v: 5})
	got, err := it.ToSlice(it.Map(iter, func(v int) int {
		return v * v
	}))
	qt.Assert(t, err, qt.ErrorMatches, "bad wolf")
	qt.Assert(t, got, qt.DeepEquals, []int{1, 4, 9, 16, 25})
}

func TestReduce(t *testing.T) {
	iter := it.FromSlice([]string{"hello", "world"})
	length, err := it.Reduce(iter, func(a int, v string) int {
		return a + len(v)
	}, 0)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, length, qt.Equals, 10)
}

func TestReduceError(t *testing.T) {
	iter := it.Chain[int](it.Count(1, 5, 1), &errorIterator[int]{v: 5})
	got, err := it.Reduce(iter, func(a, v int) int {
		return a + v
	}, 0)
	qt.Assert(t, err, qt.ErrorMatches, "bad wolf")
	qt.Assert(t, got, qt.DeepEquals, 15)
}

func TestChain(t *testing.T) {
	c := qt.New(t)

	iter := it.Chain(it.FromSlice([]string{"1", "2"}), it.FromSlice([]string{"3", "4"}))
	got, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.DeepEquals, []string{"1", "2", "3", "4"})

	// Further calls to next return false and produce the zero value.
	v := "42"
	c.Assert(iter.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")
}

func TestChainError(t *testing.T) {
	iter := it.Chain[int](it.Count(1, 5, 1), &errorIterator[int]{v: 5})
	got, err := it.ToSlice(iter)
	qt.Assert(t, err, qt.ErrorMatches, "bad wolf")
	qt.Assert(t, got, qt.DeepEquals, []int{1, 2, 3, 4, 5})
}

func TestRepeat(t *testing.T) {
	iter := it.Limit(it.Repeat(it.Count(3, 0, -1)), 9)
	got, err := it.ToSlice(iter)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, got, qt.DeepEquals, []int{3, 2, 1, 3, 2, 1, 3, 2, 1})
}

func TestRepeatError(t *testing.T) {
	iter := it.Repeat[string](&errorIterator[string]{v: "v"})
	got, err := it.ToSlice(iter)
	qt.Assert(t, err, qt.ErrorMatches, "bad wolf")
	qt.Assert(t, got, qt.DeepEquals, []string{"v"})
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
	c := qt.New(t)

	var vs []int
	iter := it.Tee(it.Count(0, 10, 2), func(v int) {
		vs = append(vs, v)
	})

	v, err := it.Next(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, 0)
	c.Assert(vs, qt.DeepEquals, []int{0})

	s, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(s, qt.DeepEquals, []int{2, 4, 6, 8})
	c.Assert(vs, qt.DeepEquals, []int{0, 2, 4, 6, 8})
}

func TestIteratorFunc(t *testing.T) {
	iter := it.Count(0, 10, 2)
	f := it.IteratorFunc[int](func(v *int) bool {
		return iter.Next(v)
	})
	got, err := it.ToSlice[int](f)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, got, qt.DeepEquals, []int{0, 2, 4, 6, 8})
}
