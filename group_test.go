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
	var kv it.KeyValue[int, it.Iterator[string]]
	var v string

	// Forward the groups iterator a couple of times.
	qt.Assert(t, qt.IsTrue(groups.Next(&kv)))
	k, group1 := kv.Split()
	qt.Assert(t, qt.Equals(k, 1))

	qt.Assert(t, qt.IsTrue(groups.Next(&kv)))
	k, group2 := kv.Split()
	qt.Assert(t, qt.Equals(k, 2))

	// Consume the second group.
	qt.Assert(t, qt.IsTrue(group2.Next(&v)))
	qt.Assert(t, qt.Equals(v, "be"))

	qt.Assert(t, qt.IsTrue(group2.Next(&v)))
	qt.Assert(t, qt.Equals(v, "it"))

	qt.Assert(t, qt.IsTrue(group2.Next(&v)))
	qt.Assert(t, qt.Equals(v, "no"))

	// Produce another group.
	qt.Assert(t, qt.IsTrue(groups.Next(&kv)))
	k, group3 := kv.Split()
	qt.Assert(t, qt.Equals(k, 5))

	// Consume the first group.
	qt.Assert(t, qt.IsTrue(group1.Next(&v)))
	qt.Assert(t, qt.Equals(v, "a"))

	// Consume the third group.
	qt.Assert(t, qt.IsTrue(group3.Next(&v)))
	qt.Assert(t, qt.Equals(v, "hello"))

	qt.Assert(t, qt.IsFalse(group3.Next(&v)))
	qt.Assert(t, qt.Equals(v, ""))

	// Produce the fourth group.
	qt.Assert(t, qt.IsTrue(groups.Next(&kv)))
	k, group4 := kv.Split()
	qt.Assert(t, qt.Equals(k, 3))

	// Consume the fourth group.
	qt.Assert(t, qt.IsTrue(group4.Next(&v)))
	qt.Assert(t, qt.Equals(v, "the"))

	qt.Assert(t, qt.IsTrue(group4.Next(&v)))
	qt.Assert(t, qt.Equals(v, "are"))

	qt.Assert(t, qt.IsFalse(group4.Next(&v)))
	qt.Assert(t, qt.Equals(v, ""))

	// Produce the fifth group.
	qt.Assert(t, qt.IsTrue(groups.Next(&kv)))
	k, group5 := kv.Split()
	qt.Assert(t, qt.Equals(k, 2))

	// Produce the sixth group.
	qt.Assert(t, qt.IsTrue(groups.Next(&kv)))
	k, group6 := kv.Split()
	qt.Assert(t, qt.Equals(k, 5))

	// There are no other groups.
	qt.Assert(t, qt.IsFalse(groups.Next(&kv)))
	k, group7 := kv.Split()
	qt.Assert(t, qt.Equals(k, 0))
	qt.Assert(t, qt.IsNil(group7))

	// Consume the sixth group.
	qt.Assert(t, qt.IsTrue(group6.Next(&v)))
	qt.Assert(t, qt.Equals(v, "again"))

	// Consume the fifth group.
	qt.Assert(t, qt.IsTrue(group5.Next(&v)))
	qt.Assert(t, qt.Equals(v, "be"))

	qt.Assert(t, qt.IsTrue(group5.Next(&v)))
	qt.Assert(t, qt.Equals(v, "it"))

	qt.Assert(t, qt.IsFalse(group5.Next(&v)))
	qt.Assert(t, qt.Equals(v, ""))

	// Check errors.
	qt.Assert(t, qt.IsNil(groups.Err()))
	qt.Assert(t, qt.IsNil(group1.Err()))
	qt.Assert(t, qt.IsNil(group2.Err()))
	qt.Assert(t, qt.IsNil(group3.Err()))
	qt.Assert(t, qt.IsNil(group4.Err()))
	qt.Assert(t, qt.IsNil(group5.Err()))
	qt.Assert(t, qt.IsNil(group4.Err()))
}
