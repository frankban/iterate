// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestKeyValueSplit(t *testing.T) {
	kv := it.KeyValue[int, string]{
		Key:   42,
		Value: "engage",
	}
	k, v := kv.Split()
	qt.Assert(t, k, qt.Equals, 42)
	qt.Assert(t, v, qt.Equals, "engage")
}

func TestZip(t *testing.T) {
	c := qt.New(t)

	keys := it.FromSlice([]string{"a", "b", "c"})
	values := it.FromSlice([]int{4, 3, 2, 1})
	iter := it.Zip(keys, values)

	got, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.DeepEquals, []it.KeyValue[string, int]{{
		Key:   "a",
		Value: 4,
	}, {
		Key:   "b",
		Value: 3,
	}, {
		Key:   "c",
		Value: 2,
	}})

	// Further calls to next return false and produce the zero value.
	v := it.KeyValue[string, int]{
		Key:   "engage",
		Value: 42,
	}
	c.Assert(iter.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, it.KeyValue[string, int]{})

	// The remaining value can still be retrieved.
	vs, err := it.ToSlice(values)
	c.Assert(err, qt.IsNil)
	c.Assert(vs, qt.DeepEquals, []int{1})
}

func TestZipError(t *testing.T) {
	c := qt.New(t)

	keys := it.FromSlice([]string{"a", "b", "c"})
	iter := it.Zip[string, int](keys, &errorIterator[int]{
		v: 42,
	})

	got, err := it.ToSlice(iter)
	c.Assert(err, qt.ErrorMatches, "bad wolf")
	c.Assert(got, qt.DeepEquals, []it.KeyValue[string, int]{{
		Key:   "a",
		Value: 42,
	}})
}

func TestUnzip(t *testing.T) {
	c := qt.New(t)

	iter := it.FromSlice(makeKeyValues())
	keys, values := it.Unzip(iter)

	var k int
	var v string

	c.Assert(keys.Next(&k), qt.IsTrue)
	c.Assert(k, qt.Equals, 1)

	c.Assert(values.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "these")

	c.Assert(values.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "are")

	c.Assert(values.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "the")

	c.Assert(keys.Next(&k), qt.IsTrue)
	c.Assert(k, qt.Equals, 2)

	c.Assert(values.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "voyages")

	c.Assert(values.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")

	c.Assert(keys.Next(&k), qt.IsTrue)
	c.Assert(k, qt.Equals, 42)

	c.Assert(keys.Next(&k), qt.IsTrue)
	c.Assert(k, qt.Equals, 47)

	c.Assert(keys.Next(&k), qt.IsFalse)
	c.Assert(k, qt.Equals, 0)

	c.Assert(keys.Err(), qt.IsNil)
	c.Assert(values.Err(), qt.IsNil)
}

func TestUnzipError(t *testing.T) {
	c := qt.New(t)

	keys, values := it.Unzip[int, string](&errorIterator[it.KeyValue[int, string]]{
		v: it.KeyValue[int, string]{
			Key:   1,
			Value: "engage",
		},
	})

	var k int
	var v string

	c.Assert(keys.Next(&k), qt.IsTrue)
	c.Assert(k, qt.Equals, 1)

	c.Assert(keys.Next(&k), qt.IsFalse)
	c.Assert(k, qt.Equals, 0)

	c.Assert(keys.Err(), qt.ErrorMatches, "bad wolf")
	// The values iterator is still able to produce values but it's already in
	// error.
	c.Assert(values.Err(), qt.ErrorMatches, "bad wolf")

	c.Assert(values.Next(&v), qt.IsTrue)
	c.Assert(v, qt.Equals, "engage")

	c.Assert(values.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, "")
}

func TestZipUnzipRoundtrip(t *testing.T) {
	s := makeKeyValues()
	got, err := it.ToSlice(it.Zip(it.Unzip(it.FromSlice(s))))
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, got, qt.DeepEquals, s)
}

func TestToMap(t *testing.T) {
	iter := it.FromSlice(makeKeyValues())
	m, err := it.ToMap(iter)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, m, qt.DeepEquals, map[int]string{
		1:  "these",
		2:  "are",
		42: "the",
		47: "voyages",
	})
}

func TestToMapError(t *testing.T) {
	m, err := it.ToMap[int, string](&errorIterator[it.KeyValue[int, string]]{
		v: it.KeyValue[int, string]{
			Key:   1,
			Value: "engage",
		},
	})
	qt.Assert(t, err, qt.ErrorMatches, "bad wolf")
	qt.Assert(t, m, qt.DeepEquals, map[int]string{
		1: "engage",
	})
}

func makeKeyValues() []it.KeyValue[int, string] {
	return []it.KeyValue[int, string]{{
		Key:   1,
		Value: "these",
	}, {
		Key:   2,
		Value: "are",
	}, {
		Key:   42,
		Value: "the",
	}, {
		Key:   47,
		Value: "voyages",
	}}
}
