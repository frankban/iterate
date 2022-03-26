// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"fmt"
	"testing"

	"github.com/frankban/iterate"
	"github.com/go-quicktest/qt"

	it "github.com/frankban/iterate"
)

func TestCount(t *testing.T) {
	tests := []struct {
		start, stop, step int
		want              []int
	}{{
		start: 0,
		stop:  10,
		step:  1,
		want:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}, {
		start: 0,
		stop:  0,
		step:  0,
		want:  nil,
	}, {
		start: 10,
		stop:  -10,
		step:  -5,
		want:  []int{10, 5, 0, -5},
	}}
	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("%d, %d, %d\n", test.start, test.stop, test.step), func(t *testing.T) {
			counter := it.Count(test.start, test.stop, test.step)
			got, err := it.ToSlice(counter)
			qt.Assert(t, qt.IsNil(err))
			qt.Assert(t, qt.DeepEquals(got, test.want))

			// Further calls to next return false and produce the zero value.
			qt.Assert(t, qt.IsFalse(counter.Next()))
			qt.Assert(t, qt.Equals(counter.Value(), 0))
		})
	}
}

func TestCountSkipValue(t *testing.T) {
	counter := it.Count(0, 100, 3)
	counter.Next()
	counter.Next()
	counter.Next()
	qt.Assert(t, qt.Equals(counter.Value(), 6))
	qt.Assert(t, qt.Equals(counter.Value(), 6))
	qt.Assert(t, qt.IsNil(counter.Err()))
}

func TestEnumerate(t *testing.T) {
	iter := it.Enumerate(it.FromSlice([]string{"a", "b", "c"}))
	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.DeepEquals(got, []it.KeyValue[int, string]{{
		Key:   0,
		Value: "a",
	}, {
		Key:   1,
		Value: "b",
	}, {
		Key:   2,
		Value: "c",
	}}))

	// Further calls to next return false and produce the zero value.
	qt.Assert(t, qt.IsFalse(iter.Next()))
	qt.Assert(t, qt.DeepEquals(iter.Value(), it.KeyValue[int, string]{}))
}

func TestEnumerateSkipValue(t *testing.T) {
	iter := it.Enumerate(it.FromSlice([]string{"a", "b", "c"}))
	iter.Next()
	iter.Next()
	want := it.KeyValue[int, string]{1, "b"}
	qt.Assert(t, qt.DeepEquals(iter.Value(), want))
	qt.Assert(t, qt.DeepEquals(iter.Value(), want))
	qt.Assert(t, qt.IsNil(iter.Err()))
}

func TestEnumerateError(t *testing.T) {
	iter := it.Enumerate[string](&errorIterator[string]{v: "v"})
	got, err := it.ToSlice(iter)
	qt.Assert(t, qt.ErrorMatches(err, "bad wolf"))
	qt.Assert(t, qt.DeepEquals(got, []iterate.KeyValue[int, string]{{
		Value: "v",
	}}))

}

func TestSum(t *testing.T) {
	sum1, err := it.Sum(it.FromSlice([]int{1, -1}))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(sum1, 0))

	sum2, err := it.Sum(it.Count(1, 11, 1))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(sum2, 55))

	sum3, err := it.Sum(it.FromSlice([]float64{0.0, 0.0}))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(sum3, 0.0))

	sum4, err := it.Sum(it.FromSlice([]string{"hello", " ", "world"}))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(sum4, "hello world"))
}
