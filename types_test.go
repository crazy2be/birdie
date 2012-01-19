package birdie

import (
	"testing"
	"bytes"
)

func TestReadInt(t* testing.T) {
	buf := bytes.NewBuffer([]byte{0x92, 0xA5, 0x7F, 0xFF, 0xFF, 0xFF})
	num, err := ReadInt(buf)
	if err != nil {
		t.Fatal(err)
	}
	// 76742655 == (((0x92 & 0x7F) << 7 | (0xA5 & 0x7F)) << 15) | (0x7FFF & 0x7FFF)
	if num != 76742655 {
		t.Fatalf("Got incorrect number. Expected %d, got %d.", 76742655, num)
	}
}