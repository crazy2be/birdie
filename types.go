package birdie

import (
	"io"
	"fmt"
)

func readByte(rd io.Reader) (byte, error) {
	buf := make([]byte, 1)
	n, err := rd.Read(buf)
	if n != 1 {
		return 0, err
	}
	return buf[0], err
}

// ReadInt reads a variable-length integer from rd and returns it and an error, if any. Encodes integers bytewise, with the first bit in each group signifiying if there are additional bytes to be read.
// I.e. in the first byte of a group, 0b1xxxxxxx signifies that the next byte should be read, while 0b0xxxxxxx signifies that this is the end of the integer.
// An example sequence of bytes might be:
//	[0b1xxxxxxx, 0b1xxxxxxx, 0b0xxxxxxx, 0b0xxxxxxx]
// Here, the first three bytes would be read, the uppermost bit in each byte discarded, and all of them combined together as
// 	(((b1 << 7) | b2) << 7) | b3
// (that is, in big-endian order)
func ReadInt(rd io.Reader) (uint64, error) {
	var num uint64
	var err error
	var b byte = 0x80
	
	for (b & 0x80) == 0x80 {
		b, err = readByte(rd)
		if err != nil {
			return num, err
		}
		
		digits := uint64(b & 0x7F)
		num = num << 7;
		num |= digits;
	}
	
	return num, nil
}

func writeByte(wr io.Writer, b byte) error {
	n, err := wr.Write([]byte{b})
	if n != 1 {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func WriteInt(wr io.Writer, num uint64) error {
	var shifty uint = 64 - 8 // 64 bit integer mines 8 bit integer (byte)
	
	// Our algorithm simply doesn't write anything for a value of zero, but we still need to write something as per the protocol. Handle it as a special case.
	if num == 0 {
		return writeByte(wr, 0x00)
	}
	
	sentFirst := false
	
	for ; shifty < 64; shifty -= 7 {
		b := byte((num & (0x7F << shifty)) >> shifty)
		if b == 0 && !sentFirst {
			continue
		}
		sentFirst = true
		if shifty > 0 {
			b = b | 0x80
		}
		err := writeByte(wr, b)
		if err != nil {
			return err
		}
	}
	return nil
}