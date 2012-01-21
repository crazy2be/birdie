package birdie

import (
	"errors"
	"io"
)

var (
	ErrNotAllWritten = errors.New("Not all bytes written to stream")
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
	var digits uint64
	
	for (b & 0x80) == 0x80 {
		b, err = readByte(rd)
		if err != nil {
			return num, err
		}
		
		digits = uint64(b & 0x7F)
		num = num << 7;
		num |= digits;
	}
	
	return num, nil
}

// WriteInt writes num to the io.Writer given by wr, using the same encoding as the ReadInt() method above. Returns an error, if any.
func WriteInt(wr io.Writer, num uint64) error {
	var shifty uint = 64 - 8 // 64 bit integer mines 8 bit integer (byte)
	
	sentFirst := false
	var b byte
	var err error
	// Maximum length is 9 bytes to encode a uint64
	buf := make([]byte, 0, 9)
	
	for ; shifty < 64; shifty -= 7 {
		b = byte((num & (0x7F << shifty)) >> shifty)
		if b == 0 && !sentFirst {
			continue
		}
		sentFirst = true
		if shifty > 0 {
			b = b | 0x80
		}
		buf = buf[:len(buf)+1]
		buf[len(buf)-1] = b
	}
	
	// We always want to send at least one byte
	if !sentFirst {
		buf = append(buf, 0x00)
	}
	
	n, err := wr.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return ErrNotAllWritten
	}
	
	return nil
}

func ReadString(rd io.Reader) (string, error) {
	length, err := ReadInt(rd)
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(rd, buf)
	return string(buf), err
}

func WriteString(wr io.Writer, data string) (error) {
	err := WriteInt(wr, uint64(len(data)))
	if err != nil {
		return err
	}
	
	_, err = io.WriteString(wr, data)
	if err != nil {
		return err
	}
	return nil
}