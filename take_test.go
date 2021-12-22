// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestDropWhile(t *testing.T) {
	c := qt.New(t)

	iter := it.DropWhile(it.Count(0, 10, 1), func(idx, v int) bool {
		return idx < 2 || idx > 4 || v == 2 || v == 7
	})

	vs, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(vs, qt.DeepEquals, []int{3, 4, 5, 6, 7, 8, 9})
}

func TestDropWhileError(t *testing.T) {
	c := qt.New(t)

	iter := it.DropWhile[string](&errorIterator[string]{
		v: "ok",
	}, func(idx int, v string) bool {
		return v == "ok"
	})

	vs, err := it.ToSlice(iter)
	c.Assert(err, qt.ErrorMatches, "bad wolf")
	c.Assert(vs, qt.IsNil)
}

func TestTakeWhile(t *testing.T) {
	c := qt.New(t)

	iter := it.TakeWhile(it.Count(0, 20, 1), func(idx, v int) bool {
		return idx < 3 || idx > 4 || v == 3 || v == 7
	})

	vs, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(vs, qt.DeepEquals, []int{0, 1, 2, 3})

	// Consuming more items does not produce values.
	for i := 0; i < 5; i++ {
		v, err := it.Next(iter)
		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.Equals, 0)
	}
}

func TestTakeWhileError(t *testing.T) {
	c := qt.New(t)

	iter := it.TakeWhile[string](&errorIterator[string]{
		v: "ok",
	}, func(idx int, v string) bool {
		return v == "ok"
	})

	vs, err := it.ToSlice(iter)
	c.Assert(err, qt.ErrorMatches, "bad wolf")
	c.Assert(vs, qt.DeepEquals, []string{"ok"})
}

func TestLimit(t *testing.T) {
	c := qt.New(t)

	iter := it.Limit(it.Count(0, 10, 1), 3)

	vs, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(vs, qt.DeepEquals, []int{0, 1, 2})
}

func TestNext(t *testing.T) {
	c := qt.New(t)

	iter := &errorIterator[string]{
		v: "ok",
	}

	v, err := it.Next[string](iter)
	c.Assert(err, qt.IsNil)
	c.Assert(v, qt.Equals, "ok")

	v, err = it.Next[string](iter)
	c.Assert(err, qt.ErrorMatches, "bad wolf")
	c.Assert(v, qt.Equals, "")
}
