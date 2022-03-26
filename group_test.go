// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	"github.com/go-quicktest/qt"

	it "github.com/frankban/iterate"
)

func TestGroupBy(t *testing.T) {
	// Group words by length.
	words := it.FromSlice([]string{
		"a",
		"be", "it", "no",
		"hello",
		"the", "are",
		"be", "it",
		"again",
	})
	groups := it.GroupBy(words, func(v string) int {
		return len(v)
	})

	// Consume groups.

	// Forward the groups iterator a couple of times.
	qt.Assert(t, qt.IsTrue(groups.Next()))
	k, group1 := groups.Value().Split()
	qt.Assert(t, qt.Equals(k, 1))

	qt.Assert(t, qt.IsTrue(groups.Next()))
	k, group2 := groups.Value().Split()
	qt.Assert(t, qt.Equals(k, 2))

	// Consume the second group.
	qt.Assert(t, qt.IsTrue(group2.Next()))
	qt.Assert(t, qt.Equals(group2.Value(), "be"))

	qt.Assert(t, qt.IsTrue(group2.Next()))
	qt.Assert(t, qt.Equals(group2.Value(), "it"))

	qt.Assert(t, qt.IsTrue(group2.Next()))
	qt.Assert(t, qt.Equals(group2.Value(), "no"))

	// Produce another group.
	qt.Assert(t, qt.IsTrue(groups.Next()))
	k, group3 := groups.Value().Split()
	qt.Assert(t, qt.Equals(k, 5))

	// Consume the first group.
	qt.Assert(t, qt.IsTrue(group1.Next()))
	qt.Assert(t, qt.Equals(group1.Value(), "a"))

	// Consume the third group.
	qt.Assert(t, qt.IsTrue(group3.Next()))
	qt.Assert(t, qt.Equals(group3.Value(), "hello"))

	qt.Assert(t, qt.IsFalse(group3.Next()))
	qt.Assert(t, qt.Equals(group3.Value(), ""))

	// Produce the fourth group.
	qt.Assert(t, qt.IsTrue(groups.Next()))
	k, group4 := groups.Value().Split()
	qt.Assert(t, qt.Equals(k, 3))

	// Consume the fourth group.
	qt.Assert(t, qt.IsTrue(group4.Next()))
	qt.Assert(t, qt.Equals(group4.Value(), "the"))

	qt.Assert(t, qt.IsTrue(group4.Next()))
	qt.Assert(t, qt.Equals(group4.Value(), "are"))

	qt.Assert(t, qt.IsFalse(group4.Next()))
	qt.Assert(t, qt.Equals(group4.Value(), ""))

	// Produce the fifth group.
	qt.Assert(t, qt.IsTrue(groups.Next()))
	k, group5 := groups.Value().Split()
	qt.Assert(t, qt.Equals(k, 2))

	// Produce the sixth group.
	qt.Assert(t, qt.IsTrue(groups.Next()))
	k, group6 := groups.Value().Split()
	qt.Assert(t, qt.Equals(k, 5))

	// There are no other groups.
	qt.Assert(t, qt.IsFalse(groups.Next()))
	k, group7 := groups.Value().Split()
	qt.Assert(t, qt.Equals(k, 0))
	qt.Assert(t, qt.IsNil(group7))

	// Consume the sixth group.
	qt.Assert(t, qt.IsTrue(group6.Next()))
	qt.Assert(t, qt.Equals(group6.Value(), "again"))

	// Consume the fifth group.
	qt.Assert(t, qt.IsTrue(group5.Next()))
	qt.Assert(t, qt.Equals(group5.Value(), "be"))

	qt.Assert(t, qt.IsTrue(group5.Next()))
	qt.Assert(t, qt.Equals(group5.Value(), "it"))

	qt.Assert(t, qt.IsFalse(group5.Next()))
	qt.Assert(t, qt.Equals(group5.Value(), ""))

	// Check errors.
	qt.Assert(t, qt.IsNil(groups.Err()))
	qt.Assert(t, qt.IsNil(group1.Err()))
	qt.Assert(t, qt.IsNil(group2.Err()))
	qt.Assert(t, qt.IsNil(group3.Err()))
	qt.Assert(t, qt.IsNil(group4.Err()))
	qt.Assert(t, qt.IsNil(group5.Err()))
	qt.Assert(t, qt.IsNil(group4.Err()))
}
