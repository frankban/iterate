// Licensed under the MIT license, see LICENSE file for details.

package iterate_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"

	it "github.com/frankban/iterate"
)

func TestLines(t *testing.T) {
	c := qt.New(t)

	r := strings.NewReader("hello\nworld")
	lines := it.Lines(r)
	var line string

	c.Assert(lines.Next(&line), qt.IsTrue)
	c.Assert(line, qt.Equals, "hello")
	c.Assert(lines.Next(&line), qt.IsTrue)
	c.Assert(line, qt.Equals, "world")

	// Further calls to next return false and produce the zero value.
	c.Assert(lines.Next(&line), qt.IsFalse)
	c.Assert(line, qt.Equals, "")

	c.Assert(lines.Err(), qt.IsNil)
}

func TestLinesError(t *testing.T) {
	c := qt.New(t)

	lines := it.Lines(errReader{})
	var line string
	c.Assert(lines.Next(&line), qt.IsFalse)
	c.Assert(line, qt.Equals, "")
	c.Assert(lines.Err(), qt.ErrorMatches, "bad wolf")
}

func TestBytes(t *testing.T) {
	c := qt.New(t)

	r := bytes.NewReader([]byte("hello\nworld"))
	b, err := it.ToSlice(it.Bytes(r))
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, "hello\nworld")
}

func TestBytesError(t *testing.T) {
	c := qt.New(t)

	bytes := it.Bytes(errReader{})
	var b byte
	c.Assert(bytes.Next(&b), qt.IsFalse)
	c.Assert(b, qt.Equals, uint8(0))
	c.Assert(bytes.Err(), qt.ErrorMatches, "bad wolf")
}

// errReader is a io.Reader implementation that always return an error.
type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("bad wolf")
}
