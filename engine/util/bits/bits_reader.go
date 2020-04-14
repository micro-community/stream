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
