package main

import (
	"encoding/binary"
	"io"
)

type demostream struct {
	r   io.ReadSeeker
	pos int
}

func (d *demostream) GetVarInt() uint64 {
	var x uint64
	var s uint
	buf := make([]byte, 1)
	for i := 0; ; i++ {
		_, err := d.r.Read(buf)
		d.pos++
		if err != nil {
			panic(err)
		}
		if buf[0] < 0x80 {
			if i > 9 || i == 9 && buf[0] > 1 {
				panic("overflow")
			}
			return x | uint64(buf[0])<<s
		}
		x |= uint64(buf[0]&0x7f) << s
		s += 7
	}
}
func (d *demostream) GetCurrentOffset() int {
	return d.pos
}
func (d *demostream) GetByte() byte {
	buf := make([]byte, 1)
	n, err := d.r.Read(buf)
	if err != nil {
		panic(err)
	}
	d.pos += n
	return buf[0]
}
func (d *demostream) GetInt() int32 {
	var x int32
	err := binary.Read(d.r, binary.LittleEndian, &x)
	if err != nil {
		panic(err)
	}
	d.pos += 4
	return x
}
func NewDemoStream(reader io.ReadSeeker) *demostream {
	stream := demostream{r: reader, pos: 0}
	return &stream
}
func (d *demostream) Read(out []byte) (int, error) {
	n, err := d.r.Read(out)
	d.pos += n
	return n, err
}
func (d *demostream) Skip(n int64) {
	d.pos += int(n)
	d.r.Seek(n, 1)
}
