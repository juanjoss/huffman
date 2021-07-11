package bits

import (
	"bufio"
	"errors"
	"io"
)

type Order int

const (
	// Least Significant Bits first (right to left)
	LSB Order = iota

	// Most Significant Bits first (left to right)
	MSB
)

/* Writing */

type writer interface {
	io.ByteWriter
	Flush() error
}

type Writer struct {
	w     writer
	order Order
	write func(uint32, uint) error
	bits  uint32
	nBits uint
	err   error
}

func (w *Writer) writeLSB(c uint32, width uint) error {
	w.bits |= c << w.nBits
	w.nBits += width

	for w.nBits >= 8 {
		if err := w.w.WriteByte(uint8(w.bits)); err != nil {
			return err
		}

		w.bits >>= 8
		w.nBits -= 8
	}

	return nil
}

func (w *Writer) writeMSB(c uint32, width uint) error {
	w.bits |= c << (32 - width - w.nBits)
	w.nBits += width

	for w.nBits >= 8 {
		if err := w.w.WriteByte(uint8(w.bits >> 24)); err != nil {
			return err
		}

		w.bits <<= 8
		w.nBits -= 8
	}

	return nil
}

func (w *Writer) WriteBits(c uint32, width uint) error {
	if w.err == nil {
		w.err = w.write(c, width)
	}

	return w.err
}

var errClosed = errors.New("bit reader/writer is closed.")

func (w *Writer) Close() error {
	if w.err != nil {
		if w.err == errClosed {
			return nil
		}

		return w.err
	}

	// Write final bits (zero padded)
	if w.nBits > 0 {
		if w.order == MSB {
			w.bits >>= 24
		}
		if w.err = w.w.WriteByte(uint8(w.bits)); w.err != nil {
			return w.err
		}
	}

	w.err = w.w.Flush()
	if w.err != nil {
		return w.err
	}

	w.err = errClosed
	return nil
}

func NewWriter(w io.Writer, order Order) *Writer {
	bw := &Writer{order: order}

	switch order {
	case LSB:
		bw.write = bw.writeLSB
	case MSB:
		bw.write = bw.writeMSB
	default:
		bw.err = errors.New("bit writer: unknown order")
		return bw
	}

	if byteWriter, ok := w.(writer); ok {
		bw.w = byteWriter
	} else {
		bw.w = bufio.NewWriter(w)
	}

	return bw
}

/* Reaging */

type Reader struct {
	r     io.ByteReader
	read  func(widht uint) (uint16, error)
	bits  uint32
	nBits uint
	err   error
}

func (r *Reader) readLSB(width uint) (uint16, error) {
	for r.nBits < width {
		x, err := r.r.ReadByte()
		if err != nil {
			return 0, err
		}

		r.bits |= uint32(x) << r.nBits
		r.nBits += 8
	}

	bits := uint16(r.bits & (1<<width - 1))
	r.bits >>= width
	r.nBits -= width

	return bits, nil
}

func (r *Reader) readMSB(width uint) (uint16, error) {
	for r.nBits < width {
		x, err := r.r.ReadByte()
		if err != nil {
			return 0, err
		}

		r.bits |= uint32(x) << (24 - r.nBits)
		r.nBits += 8
	}

	bits := uint16(r.bits >> (32 - width))
	r.bits <<= width
	r.nBits -= width

	return bits, nil
}

func (r *Reader) ReadBits(width uint) (uint16, error) {
	var bits uint16

	if r.err == nil {
		bits, r.err = r.read(width)
	}

	return bits, r.err
}

func (r *Reader) Close() error {
	if r.err != nil && r.err != errClosed {
		return r.err
	}

	r.err = errClosed
	return nil
}

func NewReader(r io.Reader, order Order) *Reader {
	br := new(Reader)
	switch order {
	case LSB:
		br.read = br.readLSB
	case MSB:
		br.read = br.readMSB
	default:
		br.err = errors.New("bit reader: unknown order.")
		return br
	}

	if byteReader, ok := r.(io.ByteReader); ok {
		br.r = byteReader
	} else {
		br.r = bufio.NewReader(r)
	}

	return br
}
