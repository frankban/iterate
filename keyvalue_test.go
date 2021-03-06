// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	"github.com/go-quicktest/qt"

	it "github.com/frankban/iterate"
)

func TestKeyValueSplit(t *testing.T) {
	kv := it.KeyValue[int, string]{
		Key:   42,
		Value: "engage",
	}
	k, v := kv.Split()
	qt.Assert(t, qt.Equals(k, 42))
	qt.Assert(t, qt.Equals(v, "engage"))
}

func TestZip(t *testing.T) {
	keys := it.FromSlice([]string{"a", "b", "c"})
	values := it.FromSlice([]int{4, 3, 2, 1})
	iter := it.Zip(keys, values)

	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, []it.KeyValue[string, int]{{
		Key:   "a",
		Value: 4,
	}, {
		Key:   "b",
		Value: 3,
	}, {
		Key:   "c",
		Value: 2,
	}}))

	// Further calls to next return false and produce the zero value.
	qt.Assert(t, qt.IsFalse(iter.Next()))
	qt.Assert(t, qt.Equals(iter.Value(), it.KeyValue[string, int]{}))

	// The remaining value can still be retrieved.
	vs, err := it.ToSlice(values)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(vs, []int{1}))
}

func TestZipError(t *testing.T) {
	keys := it.FromSlice([]string{"a", "b", "c"})
	iter := it.Zip[string, int](keys, &errorIterator[int]{
		v: 42,
	})

	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(got, []it.KeyValue[string, int]{{
		Key:   "a",
		Value: 42,
	}}))

}

func TestUnzip(t *testing.T) {
	iter := it.FromSlice(makeKeyValues())
	keys, values := it.Unzip(iter)

	qt.Assert(t, qt.IsTrue(keys.Next()))
	qt.Assert(t, qt.Equals(keys.Value(), 1))

	qt.Assert(t, qt.IsTrue(values.Next()))
	qt.Assert(t, qt.Equals(values.Value(), "these"))

	qt.Assert(t, qt.IsTrue(values.Next()))
	qt.Assert(t, qt.Equals(values.Value(), "are"))

	qt.Assert(t, qt.IsTrue(values.Next()))
	qt.Assert(t, qt.Equals(values.Value(), "the"))

	qt.Assert(t, qt.IsTrue(keys.Next()))
	qt.Assert(t, qt.Equals(keys.Value(), 2))

	qt.Assert(t, qt.IsTrue(values.Next()))
	qt.Assert(t, qt.Equals(values.Value(), "voyages"))

	qt.Assert(t, qt.IsFalse(values.Next()))
	qt.Assert(t, qt.Equals(values.Value(), ""))

	qt.Assert(t, qt.IsTrue(keys.Next()))
	qt.Assert(t, qt.Equals(keys.Value(), 42))

	qt.Assert(t, qt.IsTrue(keys.Next()))
	qt.Assert(t, qt.Equals(keys.Value(), 47))

	qt.Assert(t, qt.IsFalse(keys.Next()))
	qt.Assert(t, qt.Equals(keys.Value(), 0))

	qt.Assert(t, qt.IsNil(keys.Err()))
	qt.Assert(t, qt.IsNil(values.Err()))
}

func TestUnzipSkipValues(t *testing.T) {
	iter := it.FromSlice(makeKeyValues())
	keys, values := it.Unzip(iter)

	keys.Next()
	keys.Next()
	keys.Next()
	values.Next()
	values.Next()

	qt.Assert(t, qt.Equals(keys.Value(), 42))
	qt.Assert(t, qt.Equals(keys.Value(), 42))
	qt.Assert(t, qt.Equals(values.Value(), "are"))

	qt.Assert(t, qt.IsNil(keys.Err()))
	qt.Assert(t, qt.IsNil(values.Err()))
}

func TestUnzipError(t *testing.T) {
	keys, values := it.Unzip[int, string](&errorIterator[it.KeyValue[int, string]]{
		v: it.KeyValue[int, string]{
			Key:   1,
			Value: "engage",
		},
	})

	qt.Assert(t, qt.IsTrue(keys.Next()))
	qt.Assert(t, qt.Equals(keys.Value(), 1))

	qt.Assert(t, qt.IsFalse(keys.Next()))
	qt.Assert(t, qt.Equals(keys.Value(), 0))

	qt.Assert(t, qt.ErrorMatches(keys.Err(), "bad wolf"))
	// The values iterator is still able to produce values but it's already in
	// error.
	qt.Assert(t, qt.ErrorMatches(values.Err(), "bad wolf"))

	qt.Assert(t, qt.IsTrue(values.Next()))
	qt.Assert(t, qt.Equals(values.Value(), "engage"))

	qt.Assert(t, qt.IsFalse(values.Next()))
	qt.Assert(t, qt.Equals(values.Value(), ""))
}

func TestZipUnzipRoundtrip(t *testing.T) {
	s := makeKeyValues()
	got, err := it.ToSlice(it.Zip(it.Unzip(it.FromSlice(s))))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, s))
}

func TestToMap(t *testing.T) {
	iter := it.FromSlice(makeKeyValues())
	m, err := it.ToMap(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(m, map[int]string{
		1:  "these",
		2:  "are",
		42: "the",
		47: "voyages",
	}))

}

func TestToMapError(t *testing.T) {
	m, err := it.ToMap[int, string](&errorIterator[it.KeyValue[int, string]]{
		v: it.KeyValue[int, string]{
			Key:   1,
			Value: "engage",
		},
	})
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(m, map[int]string{
		1: "engage",
	}))

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
