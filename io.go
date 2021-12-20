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

// Next implements Iterator[T].Next by producing the next line in the reader.
func (it *lineReader) Next(v *string) bool {
	if it.scanner.Scan() {
		*v = it.scanner.Text()
		return true
	}
	*v = ""
	return false
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
	err error
}

// Next implements Iterator[T].Next by producing the next byte in the reader.
func (it *byteReader) Next(v *byte) bool {
	b, err := it.r.ReadByte()
	if err != nil {
		if err != io.EOF {
			it.err = err
		}
		*v = 0
		return false
	}
	*v = b
	return true
}

// Err implements Iterator[T].Err by propagating any errors occurred while
// reading, except for io.EOF.
func (it *byteReader) Err() error {
	return it.err
}
