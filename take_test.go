// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestTakeWhile(t *testing.T) {
	c := qt.New(t)

	iter := it.TakeWhile(it.Count(0, 20, 1), func(idx, v int) bool {
		return idx < 3 || v == 3
	})

	vs, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(vs, qt.DeepEquals, []int{0, 1, 2, 3})
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
