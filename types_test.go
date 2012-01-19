package birdie

import (
	"testing"
	"bytes"
)

func TestRead4BitInt(t* testing.T) {
	buf := bytes.NewBuffer([]byte{0x92, 0xA5, 0x7F, 0xFF, 0xFF, 0xFF})
	// (((0x92 & 0x7F) << 7 | (0xA5 & 0x7F)) << 15) | (0x7FFF & 0x7FFF)
	expected := 76742655
	
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
	expected := 164742026212933631
	
	num, err := ReadInt(buf)
	if err != nil {
		t.Fatal(err)
	}
	if num != expected {
		t.Fatalf("Got incorrect number. Expected %d, got %d.", expected, num)
	}
}