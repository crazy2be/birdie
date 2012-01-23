package birdie

import (
	"io"
)

func (c *Conn) NextPacket() (p *Packet, err error) {
	p = &Packet{}
	p.Type, err = ReadInt(c.con)
	if err != nil {
		return nil, err
	}
	p.Len, err = ReadInt(c.con)
	if err != nil {
		return nil, err
	}
	p.rd = c.con
	return
}

type Packet struct {
	Type uint64
	Len  uint64
	rd   io.Reader
}

func (p *Packet) Read(buf []byte) (n int, err error) {
	return p.rd.Read(buf)
}