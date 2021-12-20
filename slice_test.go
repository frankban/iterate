// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestFromSlice(t *testing.T) {
	c := qt.New(t)

	iter := it.FromSlice([]byte("these are the voyages"))
	var b strings.Builder
	var v byte
	for iter.Next(&v) {
		b.WriteByte(v)
	}
	c.Assert(iter.Err(), qt.IsNil)
	c.Assert(b.String(), qt.Equals, "these are the voyages")

	// Further calls to next return false and produce the zero value.
	c.Assert(iter.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, uint8(0))
}

func TestToSlice(t *testing.T) {
	// Let's take advantage of this for testing some composition as well.
	multipleOf3 := func(v int) bool {
		return v%3 == 0
	}
	lessThan500 := func(_, v int) bool {
		return v < 500
	}
	s, err := it.ToSlice(it.TakeWhile(it.Filter(it.Count(0, 20, 1), multipleOf3), lessThan500))
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, s, qt.DeepEquals, []int{0, 3, 6, 9, 12, 15, 18})
}

func TestToSliceError(t *testing.T) {
	s, err := it.ToSlice[string](&errorIterator[string]{
		v: "engage",
	})
	qt.Assert(t, err, qt.ErrorMatches, "bad wolf")
	qt.Assert(t, s, qt.DeepEquals, []string{"engage"})
}

func TestSliceRoundtrip(t *testing.T) {
	want := []string{"these", "are", "the", "voyages"}
	got, err := it.ToSlice(it.FromSlice(want))
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, got, qt.DeepEquals, want)
}
