// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	"github.com/go-quicktest/qt"

	it "github.com/frankban/iterate"
)

func TestDropWhile(t *testing.T) {
	iter := it.DropWhile(it.Count(0, 10, 1), func(idx, v int) bool {
		return idx < 2 || idx > 4 || v == 2 || v == 7
	})

	vs, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(vs, []int{3, 4, 5, 6, 7, 8, 9}))
}

func TestDropWhileError(t *testing.T) {
	iter := it.DropWhile[string](&errorIterator[string]{
		v: "ok",
	}, func(idx int, v string) bool {
		return v == "ok"
	})

	vs, err := it.ToSlice(iter)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.IsNil(vs))
}

func TestTakeWhile(t *testing.T) {
	iter := it.TakeWhile(it.Count(0, 20, 1), func(idx, v int) bool {
		return idx < 3 || idx > 4 || v == 3 || v == 7
	})

	vs, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(vs, []int{0, 1, 2, 3}))

	// Consuming more items does not produce values.
	for i := 0; i < 5; i++ {
		v, err := it.Next(iter)
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.Equals(v, 0))
	}
}

func TestTakeWhileSkipValue(t *testing.T) {
	iter := it.TakeWhile(it.Count(0, 100, 10), func(idx, v int) bool {
		return idx < 5
	})
	iter.Next()
	iter.Next()
	qt.Assert(t, qt.Equals(iter.Value(), 10))
	qt.Assert(t, qt.Equals(iter.Value(), 10))
	qt.Assert(t, qt.IsNil(iter.Err()))
}

func TestTakeWhileError(t *testing.T) {
	iter := it.TakeWhile[string](&errorIterator[string]{
		v: "ok",
	}, func(idx int, v string) bool {
		return v == "ok"
	})

	vs, err := it.ToSlice(iter)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(vs, []string{"ok"}))
}

func TestLimit(t *testing.T) {
	iter := it.Limit(it.Count(0, 10, 1), 3)

	vs, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(vs, []int{0, 1, 2}))
}

func TestNext(t *testing.T) {
	iter := &errorIterator[string]{
		v: "ok",
	}

	v, err := it.Next[string](iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(v, "ok"))

	v, err = it.Next[string](iter)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.Equals(v, ""))
}
