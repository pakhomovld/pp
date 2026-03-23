package sniff

import (
	"bytes"
	"io"
)

const DefaultSize = 8192

// Reader peeks the first N bytes from a reader without consuming them.
// After sniffing, Reader() returns an io.Reader that replays the peeked
// bytes followed by the remaining original reader.
type Reader struct {
	sample []byte
	rest   io.Reader
}

// NewReader reads up to size bytes from r into a sample buffer.
func NewReader(r io.Reader, size int) (*Reader, error) {
	buf := make([]byte, size)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	sample := buf[:n]
	return &Reader{
		sample: sample,
		rest:   r,
	}, nil
}

// Sample returns the sniffed bytes.
func (s *Reader) Sample() []byte {
	return s.sample
}

// Reader returns an io.Reader that replays the sample then continues
// with the remaining original reader.
func (s *Reader) Reader() io.Reader {
	return io.MultiReader(bytes.NewReader(s.sample), s.rest)
}
