// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

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
