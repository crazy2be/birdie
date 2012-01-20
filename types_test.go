package birdie

import (
	"testing"
	"bytes"
	"log"
)

func TestRead4BitInt(t* testing.T) {
	buf := bytes.NewBuffer([]byte{0xFF, 0xFF, 0xFF, 0x7F})
	// (((((0x7F << 7) | 0x7F) << 7) | 0x7F) << 7) | 0x7F
	expected := uint64(268435455)
	
	num, err := ReadInt(buf)
	if err != nil {
		t.Fatal(err)
	}
	if num != expected {
		t.Fatalf("Got incorrect number. Expected %d, got %d.", expected, num)
	}
}

// func TestRead8BitInt(t* testing.T) {
// 	buf := bytes.NewBuffer([]byte{0x92, 0xA5, 0x8F, 0xFF, 0x7F, 0xFF, 0xFF, 0xFF})
// 	// (((((0x92 & 0x7F) << 7 | (0xA5 & 0x7F)) << 15) | (0x8FFF & 0x7FFF)) << 31) | (0x7FFFFFFF & 0x7FFFFFFF)
// 	expected := uint64(164742026212933631)
// 	
// 	num, err := ReadInt(buf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if num != expected {
// 		t.Fatalf("Got incorrect number. Expected %d, got %d.", expected, num)
// 	}
// }

func TestWrite4ByteInt(t* testing.T) {
	num := uint64(76742655)
	expected := []byte{0xA4, 0xCB, 0xFF, 0x7F}
	buf := &bytes.Buffer{}
	
	err := WriteInt(buf, num)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatal("Got incorrect bytes! Expected", expected, "got", buf.Bytes())
	}
}

func BenchmarkEncodeDecode(b *testing.B) {
	var err error
	var num uint64
	buf := &bytes.Buffer{}
	
	for i := uint64(0); i < uint64(b.N); i++ {
		err = WriteInt(buf, i)
		if err != nil {
			log.Fatal(err)
		}
		num, err = ReadInt(buf)
		if err != nil {
			log.Fatal(err)
		}
		if num != i {
			log.Fatal("Error encoding and decoding integer. Got", num, "expected", i)
		}
	}
}