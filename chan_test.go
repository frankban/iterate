// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestFromChan(t *testing.T) {
	ch := make(chan int)
	iter := it.FromChannel(ch)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	vs, err := it.ToSlice(iter)
	qt.Assert(t, err, qt.IsNil)
	qt.Assert(t, vs, qt.DeepEquals, []int{0, 1, 2, 3, 4})
}
