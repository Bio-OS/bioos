//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package bom
// nolint
package bom

import (
	"io"
)

// Writer implements automatic BOM (Unicode Byte Order Mark) write
type Writer struct {
	wr  io.Writer // writer provided by the client
	n   int
	buf []byte // buffered data
	bom []byte // bom header
	err error  // last error
}

func (w *Writer) writeErr() error {
	err := w.err
	w.err = nil
	return err
}

// WriteWithBOM creates Writer which automatically write BOM (Unicode Byte Order Mark).
func WriteWithBOM(wt io.Writer, e Encoding) *Writer {
	// Is it already a Writer?
	b, ok := wt.(*Writer)
	if ok {
		return b
	}
	bom := genBOMHeader(e)
	return &Writer{
		wr:  wt,
		n:   0,
		bom: bom,
		buf: make([]byte, 0),
		err: nil,
	}
}

// Write is an implementation of io.Write interface.
// The bytes are taken from the underlying Write, but it will automatically write for BOMs.
func (w *Writer) Write(p []byte) (nn int, err error) {
	if len(w.bom) > 0 {
		newBuf := append(w.bom, w.buf...)
		w.buf = newBuf
		w.n += len(w.bom)
		w.bom = []byte{}
	}
	for len(p) > w.Available() && w.err == nil {
		var n int
		if w.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, w.err = w.wr.Write(p)
		} else {
			n = copy(w.buf[w.n:], p)
			w.n += n
			w.Flush()
		}
		nn += n
		p = p[n:]
	}
	if w.err != nil {
		return nn, w.err
	}
	n := copy(w.buf[w.n:], p)
	w.n += n
	nn += n
	return nn, nil
}

// Available returns how many bytes are unused in the buffer.
func (w *Writer) Available() int { return len(w.buf) - w.n }

// Buffered returns the number of bytes that have been written into the current buffer.
func (w *Writer) Buffered() int { return w.n }

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	if w.err != nil {
		return w.err
	}
	if w.n == 0 {
		return nil
	}
	n, err := w.wr.Write(w.buf[0:w.n])
	if n < w.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < w.n {
			copy(w.buf[0:w.n-n], w.buf[n:w.n])
		}
		w.n -= n
		w.err = err
		return err
	}
	w.n = 0
	return nil
}

func genBOMHeader(e Encoding) []byte {
	switch e {
	case UTF8:
		return []byte{0xEF, 0xBB, 0xBF}
	case UTF16BigEndian:
		return []byte{0xFE, 0xFF}
	case UTF16LittleEndian:
		return []byte{0xFF, 0xFE}
	case UTF32BigEndian:
		return []byte{0x00, 0x00, 0xFE, 0xFF}
	case UTF32LittleEndian:
		return []byte{0xFF, 0xFE, 0x00, 0x00}
	default:
		return []byte{}
	}
}
