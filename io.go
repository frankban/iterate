// Licensed under the MIT license, see LICENSE file for details.

package iterate

import (
	"bufio"
	"io"
)

// Lines returns an iterator producing lines from the given reader.
func Lines(r io.Reader) Iterator[string] {
	return &lineReader{
		scanner: *bufio.NewScanner(r),
	}
}

type lineReader struct {
	scanner bufio.Scanner
}

// Next implements Iterator[T].Next.
func (it *lineReader) Next() bool {
	return it.scanner.Scan()
}

// Value implements Iterator[T].Value by returning the next line from the
// reader.
func (it *lineReader) Value() string {
	return it.scanner.Text()
}

// Err implements Iterator[T].Err by propagating any errors occurred while
// reading, except for io.EOF.
func (it *lineReader) Err() error {
	return it.scanner.Err()
}

// Bytes returns an iterator producing bytes from the given reader.
func Bytes(r io.Reader) Iterator[byte] {
	return &byteReader{
		r: *bufio.NewReader(r),
	}
}

type byteReader struct {
	r   bufio.Reader
	b   byte
	err error
}

// Next implements Iterator[T].Next by producing the next byte in the reader.
func (it *byteReader) Next() bool {
	if it.err != nil {
		return false
	}
	it.b, it.err = it.r.ReadByte()
	if it.err == nil {
		return true
	}
	if it.err == io.EOF {
		it.err = nil
	}
	return false
}

// Value implements Iterator[T].Value by returning the next byte from the
// reader.
func (it *byteReader) Value() byte {
	return it.b
}

// Err implements Iterator[T].Err by propagating any errors occurred while
// reading, except for io.EOF.
func (it *byteReader) Err() error {
	return it.err
}
