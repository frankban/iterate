// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"fmt"
	"testing"

	"github.com/frankban/iterate"
	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestCount(t *testing.T) {
	c := qt.New(t)

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
		c.Run(fmt.Sprintf("%d, %d, %d\n", test.start, test.stop, test.step), func(c *qt.C) {
			counter := it.Count(test.start, test.stop, test.step)
			got, err := it.ToSlice(counter)
			c.Assert(err, qt.IsNil)
			c.Assert(got, qt.DeepEquals, test.want)

			// Further calls to next return false and produce the zero value.
			v := 42
			c.Assert(counter.Next(&v), qt.IsFalse)
			c.Assert(v, qt.Equals, 0)
		})
	}
}

func TestEnumerate(t *testing.T) {
	c := qt.New(t)

	iter := it.Enumerate(it.FromSlice([]string{"a", "b", "c"}))
	got, err := it.ToSlice(iter)
	c.Assert(err, qt.IsNil)
	c.Assert(got, qt.DeepEquals, []it.KeyValue[int, string]{{
		Key:   0,
		Value: "a",
	}, {
		Key:   1,
		Value: "b",
	}, {
		Key:   2,
		Value: "c",
	}})

	// Further calls to next return false and produce the zero value.
	v := it.KeyValue[int, string]{
		Key:   42,
		Value: "engage",
	}
	c.Assert(iter.Next(&v), qt.IsFalse)
	c.Assert(v, qt.Equals, it.KeyValue[int, string]{})
}

func TestEnumerateError(t *testing.T) {
	iter := it.Enumerate[string](&errorIterator[string]{v: "v"})
	got, err := it.ToSlice(iter)
	qt.Assert(t, err, qt.ErrorMatches, "bad wolf")
	qt.Assert(t, got, qt.DeepEquals, []iterate.KeyValue[int, string]{{
		Value: "v",
	}})
}

func TestSum(t *testing.T) {
	c := qt.New(t)

	sum1, err := it.Sum(it.FromSlice([]int{1, -1}))
	c.Assert(err, qt.IsNil)
	c.Assert(sum1, qt.Equals, 0)

	sum2, err := it.Sum(it.Count(1, 11, 1))
	c.Assert(err, qt.IsNil)
	c.Assert(sum2, qt.Equals, 55)

	sum3, err := it.Sum(it.FromSlice([]float64{0.0, 0.0}))
	c.Assert(err, qt.IsNil)
	c.Assert(sum3, qt.Equals, 0.0)

	sum4, err := it.Sum(it.FromSlice([]string{"hello", " ", "world"}))
	c.Assert(err, qt.IsNil)
	c.Assert(sum4, qt.Equals, "hello world")
}
