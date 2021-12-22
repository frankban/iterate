// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestGroupBy(t *testing.T) {
	c := qt.New(t)

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
	c.Assert(groups.Next(&kv), qt.IsTrue)
	k, group1 := kv.Split()
	c.Assert(k, qt.Equals, 1)

	c.Assert(groups.Next(&kv), qt.IsTrue)
	k, group2 := kv.Split()
	c.Assert(k, qt.Equals, 2)

	// Consume the second group.
	c.Assert(group2.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "be")

	c.Assert(group2.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "it")

	c.Assert(group2.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "no")

	// Produce another group.
	c.Assert(groups.Next(&kv), qt.IsTrue)
	k, group3 := kv.Split()
	c.Assert(k, qt.Equals, 5)

	// Consume the first group.
	c.Assert(group1.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "a")

	// Consume the third group.
	c.Assert(group3.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "hello")

	c.Assert(group3.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")

	// Produce the fourth group.
	c.Assert(groups.Next(&kv), qt.IsTrue)
	k, group4 := kv.Split()
	c.Assert(k, qt.Equals, 3)

	// Consume the fourth group.
	c.Assert(group4.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "the")

	c.Assert(group4.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "are")

	c.Assert(group4.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")

	// Produce the fifth group.
	c.Assert(groups.Next(&kv), qt.IsTrue)
	k, group5 := kv.Split()
	c.Assert(k, qt.Equals, 2)

	// Produce the sixth group.
	c.Assert(groups.Next(&kv), qt.IsTrue)
	k, group6 := kv.Split()
	c.Assert(k, qt.Equals, 5)

	// There are no other groups.
	c.Assert(groups.Next(&kv), qt.IsFalse)
	k, group7 := kv.Split()
	c.Assert(k, qt.Equals, 0)
	c.Assert(group7, qt.IsNil)

	// Consume the sixth group.
	c.Assert(group6.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "again")

	// Consume the fifth group.
	c.Assert(group5.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "be")

	c.Assert(group5.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "it")

	c.Assert(group5.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")

	// Check errors.
	c.Assert(groups.Err(), qt.IsNil)
	c.Assert(group1.Err(), qt.IsNil)
	c.Assert(group2.Err(), qt.IsNil)
	c.Assert(group3.Err(), qt.IsNil)
	c.Assert(group4.Err(), qt.IsNil)
	c.Assert(group5.Err(), qt.IsNil)
	c.Assert(group4.Err(), qt.IsNil)
}
