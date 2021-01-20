package bits

import (
	"io"
)

//Reader for bits
type Reader struct {
	R    io.Reader
	n    int
	bits uint64
}

//ReadBits64 read a b664
func (rdr *Reader) ReadBits64(n int) (bits uint64, err error) {
	if rdr.n < n {
		var b [8]byte
		var got int
		want := (n - rdr.n + 7) / 8
		if got, err = rdr.R.Read(b[:want]); err != nil {
			return
		}
		if got < want {
			err = io.EOF
			return
		}
		for i := 0; i < got; i++ {
			rdr.bits <<= 8
			rdr.bits |= uint64(b[i])
		}
		rdr.n += got * 8
	}
	bits = rdr.bits >> uint(rdr.n-n)
	rdr.bits ^= bits << uint(rdr.n-n)
	rdr.n -= n
	return
}

//ReadBits read a bits
func (rdr *Reader) ReadBits(n int) (bits uint, err error) {
	var bits64 uint64
	if bits64, err = rdr.ReadBits64(n); err != nil {
		return
	}
	bits = uint(bits64)
	return
}

func (rdr *Reader) Read(p []byte) (n int, err error) {
	for n < len(p) {
		want := 8
		if len(p)-n < want {
			want = len(p) - n
		}
		var bits uint64
		if bits, err = rdr.ReadBits64(want * 8); err != nil {
			break
		}
		for i := 0; i < want; i++ {
			p[n+i] = byte(bits >> uint((want-i-1)*8))
		}
		n += want
	}
	return
}

//Writer of bits
type Writer struct {
	W    io.Writer
	n    int
	bits uint64
}

//WriteBits64 write a b64
func (wtr *Writer) WriteBits64(bits uint64, n int) (err error) {
	if wtr.n+n > 64 {
		move := uint(64 - wtr.n)
		mask := bits >> move
		wtr.bits = (wtr.bits << move) | mask
		wtr.n = 64
		if err = wtr.FlushBits(); err != nil {
			return
		}
		n -= int(move)
		bits ^= (mask << move)
	}
	wtr.bits = (wtr.bits << uint(n)) | bits
	wtr.n += n
	return
}

//WriteBits write n bits
func (wtr *Writer) WriteBits(bits uint, n int) (err error) {
	return wtr.WriteBits64(uint64(bits), n)
}

func (wtr *Writer) Write(p []byte) (n int, err error) {
	for n < len(p) {
		if err = wtr.WriteBits64(uint64(p[n]), 8); err != nil {
			return
		}
		n++
	}
	return
}

//FlushBits in buffer
func (wtr *Writer) FlushBits() (err error) {
	if wtr.n > 0 {
		var b [8]byte
		bits := wtr.bits
		if wtr.n%8 != 0 {
			bits <<= uint(8 - (wtr.n % 8))
		}
		want := (wtr.n + 7) / 8
		for i := 0; i < want; i++ {
			b[i] = byte(bits >> uint((want-i-1)*8))
		}
		if _, err = wtr.W.Write(b[:want]); err != nil {
			return
		}
		wtr.n = 0
	}
	return
}

//PutUInt64BE add a uint64 BE
func PutUInt64BE(b []byte, res uint64, n int) {
	n /= 8
	for i := 0; i < n; i++ {
		b[n-i-1] = byte(res)
		res >>= 8
	}
	return
}

//PutUIntBE add a uint BE
func PutUIntBE(b []byte, res uint, n int) {
	PutUInt64BE(b, uint64(res), n)
}

//WriteBytes write bytes in w
func WriteBytes(w io.Writer, b []byte, n int) (err error) {
	if len(b) < n {
		b = append(b, make([]byte, n-len(b))...)
	}
	_, err = w.Write(b[:n])
	return
}

func WriteUInt64BE(w io.Writer, val uint64, n int) (err error) {
	n /= 8
	var b [8]byte
	for i := n - 1; i >= 0; i-- {
		b[i] = byte(val)
		val >>= 8
	}
	return WriteBytes(w, b[:], n)
}

func WriteUIntBE(w io.Writer, val uint, n int) (err error) {
	return WriteUInt64BE(w, uint64(val), n)
}

func WriteInt64BE(w io.Writer, val int64, n int) (err error) {
	n /= 8
	var uval uint
	if val < 0 {
		uval = uint((1 << uint(n*8)) + val)
	} else {
		uval = uint(val)
	}
	return WriteUIntBE(w, uval, n)
}

func WriteIntBE(w io.Writer, val int, n int) (err error) {
	return WriteInt64BE(w, int64(val), n)
}

func WriteString(w io.Writer, val string, n int) (err error) {
	return WriteBytes(w, []byte(val), n)
}
