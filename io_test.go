// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/go-quicktest/qt"

	it "github.com/frankban/iterate"
)

func TestLines(t *testing.T) {
	r := strings.NewReader("hello\nworld")
	lines := it.Lines(r)

	qt.Assert(t, qt.IsTrue(lines.Next()))
	qt.Assert(t, qt.Equals(lines.Value(), "hello"))
	qt.Assert(t, qt.IsTrue(lines.Next()))
	qt.Assert(t, qt.Equals(lines.Value(), "world"))

	// Further calls to next return false and produce the zero value.
	qt.Assert(t, qt.IsFalse(lines.Next()))
	qt.Assert(t, qt.Equals(lines.Value(), ""))

	qt.Assert(t, qt.IsNil(lines.Err()))
}

func TestLinesError(t *testing.T) {
	lines := it.Lines(errReader{})
	qt.Assert(t, qt.IsFalse(lines.Next()))
	qt.Assert(t, qt.Equals(lines.Value(), ""))
	qt.Assert(t, qt.ErrorMatches(lines.Err(), "bad wolf"))
}

func TestBytes(t *testing.T) {
	r := bytes.NewReader([]byte("hello\nworld"))
	b, err := it.ToSlice(it.Bytes(r))
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.Equals(string(b), "hello\nworld"))
}

func TestBytesError(t *testing.T) {
	bytes := it.Bytes(errReader{})
	qt.Assert(t, qt.IsFalse(bytes.Next()))
	qt.Assert(t, qt.Equals(bytes.Value(), 0))
	qt.Assert(t, qt.ErrorMatches(bytes.Err(), "bad wolf"))
}

// errReader is a io.Reader implementation that always return an error.
type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("bad wolf")
}
