package birdie

import (
	"testing"
	"bytes"
)

func TestRead4BitInt(t* testing.T) {
	buf := bytes.NewBuffer([]byte{0x92, 0xA5, 0x7F, 0xFF, 0xFF, 0xFF})
	// (((0x92 & 0x7F) << 7 | (0xA5 & 0x7F)) << 15) | (0x7FFF & 0x7FFF)
	expected := uint64(76742655)
	
	num, err := ReadInt(buf)
	if err != nil {
		t.Fatal(err)
	}
	if num != expected {
		t.Fatalf("Got incorrect number. Expected %d, got %d.", expected, num)
	}
}

func TestRead8BitInt(t* testing.T) {
	buf := bytes.NewBuffer([]byte{0x92, 0xA5, 0x8F, 0xFF, 0x7F, 0xFF, 0xFF, 0xFF})
	// (((((0x92 & 0x7F) << 7 | (0xA5 & 0x7F)) << 15) | (0x8FFF & 0x7FFF)) << 31) | (0x7FFFFFFF & 0x7FFFFFFF)
	expected := uint64(164742026212933631)
	
	num, err := ReadInt(buf)
	if err != nil {
		t.Fatal(err)
	}
	if num != expected {
		t.Fatalf("Got incorrect number. Expected %d, got %d.", expected, num)
	}
}

func TestWriteIntBytewise(t *testing.T) {
	num := uint64(1908421)
	buf := &bytes.Buffer{}
	
	writeIntBytewise(buf, num)
	b := buf.Bytes()
	result := (((((((uint64(b[0]) << 8) | uint64(b[1])) << 8) | uint64(b[2])) << 8) | uint64(b[3])))
	if result != num {
		t.Fatal("Got", result, "expected", num)
	}
}

func TestWrite4ByteInt(t* testing.T) {
	num := uint64(76742655)
	expected := []byte{0x92, 0xA5, 0x7F, 0xFF}
	buf := &bytes.Buffer{}
	
	err := WriteInt(buf, num)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatal("Got incorrect bytes! Expected", expected, "got", buf.Bytes())
	}
}