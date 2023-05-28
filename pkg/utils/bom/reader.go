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

package bom

import (
	"errors"
	"io"
)

// Encoding is type alias for detected UTF encoding.
type Encoding int

// Constants to identify detected UTF encodings.
const (
	// Unknown encoding, returned when no BOM was detected
	Unknown Encoding = iota

	// UTF8 , BOM bytes: EF BB BF
	UTF8

	// UTF16BigEndian UTF-16, big-endian, BOM bytes: FE FF
	UTF16BigEndian

	// UTF16LittleEndian UTF-16, little-endian, BOM bytes: FF FE
	UTF16LittleEndian

	// UTF32BigEndian UTF-32, big-endian, BOM bytes: 00 00 FE FF
	UTF32BigEndian

	// UTF32LittleEndian UTF-32, little-endian, BOM bytes: FF FE 00 00
	UTF32LittleEndian
)

// String returns a user-friendly string representation of the encoding. Satisfies fmt.Stringer interface.
func (e Encoding) String() string {
	switch e {
	case UTF8:
		return "UTF8"
	case UTF16BigEndian:
		return "UTF16BigEndian"
	case UTF16LittleEndian:
		return "UTF16LittleEndian"
	case UTF32BigEndian:
		return "UTF32BigEndian"
	case UTF32LittleEndian:
		return "UTF32LittleEndian"
	default:
		return "Unknown"
	}
}

const maxConsecutiveEmptyReads = 100

// ReadSkipBOM creates Reader which automatically detects BOM (Unicode Byte Order Mark) and removes it as necessary.
// It also returns the encoding detected by the BOM.
// If the detected encoding is not needed, you can call the ReadSkipBOM function.
func ReadSkipBOM(rd io.Reader) (*Reader, Encoding) {
	// Is it already a Reader?
	b, ok := rd.(*Reader)
	if ok {
		return b, Unknown
	}

	enc, left, err := detectUtf(rd)
	return &Reader{
		rd:  rd,
		buf: left,
		err: err,
	}, enc
}

// Reader implements automatic BOM (Unicode Byte Order Mark) checking and
// removing as necessary for an io.Reader object.
type Reader struct {
	rd  io.Reader // reader provided by the client
	buf []byte    // buffered data
	err error     // last error
}

// Read is an implementation of io.Reader interface.
// The bytes are taken from the underlying Reader, but it checks for BOMs, removing them as necessary.
func (r *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if r.buf == nil {
		if r.err != nil {
			return 0, r.readErr()
		}

		return r.rd.Read(p)
	}

	// copy as much as we can
	n = copy(p, r.buf)
	r.buf = nilIfEmpty(r.buf[n:])
	return n, nil
}

func (r *Reader) readErr() error {
	err := r.err
	r.err = nil
	return err
}

var errNegativeRead = errors.New("utfbom: reader returned negative count from Read")

func detectUtf(rd io.Reader) (enc Encoding, buf []byte, err error) {
	buf, err = readBOM(rd)

	if len(buf) >= 4 {
		if isUTF32BigEndianBOM4(buf) {
			return UTF32BigEndian, nilIfEmpty(buf[4:]), err
		}
		if isUTF32LittleEndianBOM4(buf) {
			return UTF32LittleEndian, nilIfEmpty(buf[4:]), err
		}
	}

	if len(buf) > 2 && isUTF8BOM3(buf) {
		return UTF8, nilIfEmpty(buf[3:]), err
	}

	if (err != nil && err != io.EOF) || (len(buf) < 2) {
		return Unknown, nilIfEmpty(buf), err
	}

	if isUTF16BigEndianBOM2(buf) {
		return UTF16BigEndian, nilIfEmpty(buf[2:]), err
	}
	if isUTF16LittleEndianBOM2(buf) {
		return UTF16LittleEndian, nilIfEmpty(buf[2:]), err
	}

	return Unknown, nilIfEmpty(buf), err
}

func readBOM(rd io.Reader) (buf []byte, err error) {
	const maxBOMSize = 4
	var bom [maxBOMSize]byte // used to read BOM

	// read as many bytes as possible
	for nEmpty, n := 0, 0; err == nil && len(buf) < maxBOMSize; buf = bom[:len(buf)+n] {
		if n, err = rd.Read(bom[len(buf):]); n < 0 {
			panic(errNegativeRead)
		}
		if n > 0 {
			nEmpty = 0
		} else {
			nEmpty++
			if nEmpty >= maxConsecutiveEmptyReads {
				err = io.ErrNoProgress
			}
		}
	}
	return
}

func isUTF32BigEndianBOM4(buf []byte) bool {
	return buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0xFE && buf[3] == 0xFF
}

func isUTF32LittleEndianBOM4(buf []byte) bool {
	return buf[0] == 0xFF && buf[1] == 0xFE && buf[2] == 0x00 && buf[3] == 0x00
}

func isUTF8BOM3(buf []byte) bool {
	return buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF
}

func isUTF16BigEndianBOM2(buf []byte) bool {
	return buf[0] == 0xFE && buf[1] == 0xFF
}

func isUTF16LittleEndianBOM2(buf []byte) bool {
	return buf[0] == 0xFF && buf[1] == 0xFE
}

func nilIfEmpty(buf []byte) (res []byte) {
	if len(buf) > 0 {
		res = buf
	}
	return
}
