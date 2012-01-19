package birdie

import (
	"io"
)

type Packet struct {
	Type int16
	Len  int32
	rd   io.Reader
}

func (p *Packet) Read(buf []byte) (n int, err error) {
	return p.rd.Read(buf)
}