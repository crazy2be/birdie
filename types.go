package birdie

import (
	"io"
)

func readByte(rd io.Reader) (byte, error) {
	buf := make([]byte, 1)
	n, err := rd.Read(buf)
	if n != 1 {
		return 0, err
	}
	return buf[0], err
}

// ReadUint16 reads two bytes and ors them together, with b1 << 8 | b2. Returns an error, if any. NOT VARIABLE LENGTH.
func readUint16(rd io.Reader) (uint16, error) {
	b1, err := readByte(rd)
	if err != nil {
		return 0, err
	}
	b2, err := readByte(rd)
	if err != nil {
		return uint16(b1), err
	}
	return uint16(b1) << 8 | uint16(b2), nil
}

func readUint32(rd io.Reader) (uint32, error) {
	s1, err := readUint16(rd)
	if err != nil {
		return 0, err
	}
	s2, err := readUint16(rd)
	if err != nil {
		return uint32(s1), err
	}
	return uint32(s1) << 16 | uint32(s2), nil
}

// ReadInt reads a variable-length integer from rd and returns it and an error, if any. Encodes integers bytewise in groups sized by powers of two, with the first bit in each group signifiying if there are additional byte groups to be read.
// I.e. in the first byte of a group, 0b1xxxxxxx signifies that the next byte should be read, while 0b0xxxxxxx signifies that this is the end of the integer.
// An example sequence of bytes might be:
//	[0b1xxxxxxx, 0b1xxxxxxx, 0b0xxxxxxx, 0b0xxxxxxx]
// Note that 2^0 == 1, so the first group is one byte. 2^1 == 2, so the second group involves reading a total of two bytes (including the one already read, so one more). 2^2 == 4, so the third group involves reading a total of four bytes. Any additional bytes after this point will not be read by this method.
func ReadInt(rd io.Reader) (uint64, error) {
	var num uint64
	
	b1, err := readByte(rd)
	if (err != nil) {
		return 0, err
	}
	
	num = uint64(b1 & 0x7F)
	if (b1 & 0x80 != 0x80) {
		return num, nil
	}
	
	
	b2, err := readByte(rd)
	if err != nil {
		return num, err
	}
	
	num = num << 7
	num |= uint64(b2 & 0x7F)
	if (b2 & 0x80) != 0x80 {
		return num, nil
	}
	
	
	s1, err := readUint16(rd)
	if err != nil {
		return num, err
	}
	
	num = num << 15
	num |= uint64(s1 & 0x7FFF)
	if (s1 & 0x8000) != 0x8000 {
		return num, nil
	}
	
	
	i1, err := readUint32(rd)
	if err != nil {
		return num, err
	}
	
	num = num << 31
	num |= uint64(i1 & 0x7FFFFFFF)
	if (i1 & 0x80000000) == 0x80000000 {
		panic("Integers longer than 8 bytes (uint64) not implemented in ReadInt()!")
	}
	return num, nil
}