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

// next7BitByte reads a single byte from rd, shifts *num up by 7, and stores the low 7 bytes of the newly read byte in the now-freed low 7 bits of *num. Returns if the high bit was set (and we should continue reading), and the error, if any.
func next7BitByte(rd io.Reader, num *uint64) (bool, error) {
	b, err := readByte(rd)
	if err != nil {
		return false, err
	}
	
	digits := uint64(b & 0x7F)
	*num = *num << 7;
	*num |= digits;
	
	if (b & 0x80 == 0x80) {
		return true, nil
	}
	
	return false, nil
}

// ReadInt reads a variable-length integer from rd and returns it and an error, if any. Encodes integers bytewise in groups sized by powers of two, with the first bit in each group signifiying if there are additional byte groups to be read.
// I.e. in the first byte of a group, 0b1xxxxxxx signifies that the next byte should be read, while 0b0xxxxxxx signifies that this is the end of the integer.
// An example sequence of bytes might be:
//	[0b1xxxxxxx, 0b1xxxxxxx, 0b0xxxxxxx, 0b0xxxxxxx]
// Note that 2^0 == 1, so the first group is one byte. 2^1 == 2, so the second group involves reading a total of two bytes (including the one already read, so one more). 2^2 == 4, so the third group involves reading a total of four bytes. Any additional bytes after this point will not be read by this method.
func ReadInt(rd io.Reader) (uint64, error) {
	var num uint64
	
	// First two chunks are bytes
	for i := 0; i < 2; i++ {
		highBit, err := next7BitByte(rd, &num)
		if err != nil {
			return 0, err
		}
		if !highBit {
			return num, nil
		}
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

// Returns the number of bytes required to write the given integer to a network stream using WriteInt().
func bytesNeeded(num uint64) int {
	if (num < 0x7F) {
		return 1
	} else if (num < 0x7F7F) {
		return 2
	} else if (num < 0x7F7F7FFF) {
		return 4
	} else if (num < 0x7F7F7FFF7FFFFFFF) {
		return 8
	}
	return -1
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

func hex(num uint64) uint64 {
	fmt.Printf("%d or %x\n", num, num)
	return num
}

func writeIntBytewise(wr io.Writer, num uint64) error {
	writeByte(wr, byte(hex((num & (0xFF << 24)) >> 24)))
	writeByte(wr, byte(hex((num & (0xFF << 16)) >> 16)))
	writeByte(wr, byte(hex((num & (0xFF << 8)) >> 8)))
	writeByte(wr, byte(hex((num & (0xFF)))))
	return nil
}

func WriteInt(wr io.Writer, num uint64) error {
	fmt.Printf("%d or %x\n", num, num)
	if (num & 0xFF << 8 * 7 == num) {
		println("Only high byte is used!")
	}
	println("Requires", bytesNeeded(num), "bytes to encode.")
	
	numbytes := bytesNeeded(num)
	if (numbytes == 1) {
		err := writeByte(wr, byte(num))
		return err
	} else if (numbytes == 2) {
		panic("Not implemented!")
	} else if (numbytes == 4) {
		hex(num)
		if (num & 0x80) == 0x80 {
			num += (0x80 & num) << 1
			num = num & 0xFFFFFF7F
		}
		hex(num)
		if (num & (0x80 << 8)) == 0x80 << 8 {
			num += (0x08000 & num) << 1
			num = num & 0xFFFF7F7F
		}
		hex(num)
// 		if (num & 0x80000000) == 0x80000000 {
// 			num = num & 0xFF7F7F7F
// 			num += 0x080000000 << 1
// 		}
		// TODO: Check errors
		writeIntBytewise(wr, num)
	}
	if (num & 0x0000007F) == num {
		err := writeByte(wr, byte(num & 0x0000007F))
		if err != nil {
			return err
		}
		return nil
	}
	if (num & 0x00003FFF) == num {
		
	}
	return nil
}