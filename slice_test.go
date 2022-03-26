// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"strings"
	"testing"

	"github.com/go-quicktest/qt"

	it "github.com/frankban/iterate"
)

func TestFromSlice(t *testing.T) {
	iter := it.FromSlice([]byte("these are the voyages"))
	var b strings.Builder
	for iter.Next() {
		b.WriteByte(iter.Value())
	}
	qt.Assert(t, qt.IsNil(iter.Err()))
	qt.Assert(t, qt.Equals(b.String(), "these are the voyages"))

	// Further calls to next return false and produce the zero value.
	qt.Assert(t, qt.IsFalse(iter.Next()))
	qt.Assert(t, qt.Equals(iter.Value(), 0))
}

func TestFromSliceSkipValue(t *testing.T) {
	iter := it.FromSlice([]int{1, 2, 3})
	iter.Next()
	iter.Next()
	qt.Assert(t, qt.Equals(iter.Value(), 2))
	qt.Assert(t, qt.Equals(iter.Value(), 2))
	qt.Assert(t, qt.IsNil(iter.Err()))
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
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(s, []int{0, 3, 6, 9, 12, 15, 18}))
}

func TestToSliceError(t *testing.T) {
	s, err := it.ToSlice[string](&errorIterator[string]{
		v: "engage",
	})
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(s, []string{"engage"}))
}

func TestSliceRoundtrip(t *testing.T) {
	want := []string{"these", "are", "the", "voyages"}
	got, err := it.ToSlice(it.FromSlice(want))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, want))
}
