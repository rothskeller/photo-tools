package tiff

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var testInput1 = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // header, big-endian
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // pointer to IFD0
	/* 0008 */ 0x00, 0x02, // 2 tags in IFD0
	/* 000A */ 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x31, 0x32, 0x33, 0x34, // tag 1, bytes '1234'
	/* 0016 */ 0x00, 0x02, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x26, // tag 2, ptr to IFD1
	/* 0022 */ 0x00, 0x00, 0x00, 0x46, // ptr to IFD2
	/* 0026 */ 0x00, 0x01, // 1 tag in IFD1
	/* 0028 */ 0x00, 0x03, 0x00, 0x02, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x38, // tag 3, string, "Hello, world"
	/* 0034 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0038 */ 'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', 0x00, // tag 3 data
	/* 0045 */ 0x00, // padding
	/* 0046 */ 0x00, 0x01, // 1 tag in IFD2
	/* 0048 */ 0x00, 0x04, 0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x58, // tag 4, rational, 2/3
	/* 0054 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0058 */ 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, // tag 4 data
}

func TestInput1(t *testing.T) {
	var tl TIFF
	err := tl.Read(bytes.NewReader(testInput1))
	if err != nil {
		t.Fatal(err)
	}
	tag := tl.IFD0().Tag(1)
	if tag == nil {
		t.Fatal("no tag 1")
	}
	if by, err := tag.AsBytes(); err != nil {
		t.Fatalf("tag 1 %s", err)
	} else if !bytes.Equal(by, []byte("1234")) {
		t.Errorf("tag 1 wrong value")
	}
	if tag = tl.IFD0().Tag(2); tag == nil {
		t.Fatal("no tag 2")
	}
	ifd, err := tag.AsIFD()
	if err != nil {
		t.Fatalf("tag 2 %s", err)
	}
	if tag = ifd.Tag(3); tag == nil {
		t.Fatal("no tag 3")
	}
	if s, err := tag.AsString(); err != nil {
		t.Fatalf("tag 3 %s", err)
	} else if s != "Hello, world" {
		t.Errorf("tag 3 wrong value")
	}
	if ifd, err = tl.IFD0().NextIFD(); err != nil {
		t.Fatalf("next ifd %s", err)
	}
	if tag = ifd.Tag(4); tag == nil {
		t.Fatal("no tag 4")
	}
	if r, err := tag.AsRationals(); err != nil {
		t.Fatalf("tag 4 %s", err)
	} else if len(r) != 2 || r[0] != 2 || r[1] != 3 {
		t.Errorf("tag 4 wrong values")
	}
}

// Identical data except for adding a ! at the end of the string.  But since we
// have redone the layout with the largest IFD first, the order is now IFD1,
// IFD0, IFD2.
var testOutput1 = []byte{
	0x4d, 0x4d, 0x00, 0x2a, // header, big-endian
	0x00, 0x00, 0x00, 0x28, // pointer to IFD0
	0x00, 0x01, // 1 tag in IFD1
	0x00, 0x03, 0x00, 0x02, 0x00, 0x00, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x1a, // tag 3, string, "Hello, world!"
	0x00, 0x00, 0x00, 0x00, // no next pointer
	'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!', 0x00, // tag 3 data
	0x00, 0x02, // 2 tags in IFD0
	0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x31, 0x32, 0x33, 0x34, // tag 1, bytes '1234'
	0x00, 0x02, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, // tag 2, ptr to IFD1
	0x00, 0x00, 0x00, 0x46, // ptr to IFD2
	0x00, 0x01, // 1 tag in IFD2
	0x00, 0x04, 0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x58, // tag 4, rational, 2/3
	0x00, 0x00, 0x00, 0x00, // no next pointer
	0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, // tag 4 data
}

func TestWrite1(t *testing.T) {
	var tl TIFF
	tl.Read(bytes.NewReader(testInput1))
	ifd0 := tl.IFD0()
	ifd1, _ := ifd0.Tag(2).AsIFD()
	tag3 := ifd1.Tag(3)
	tag3.SetString("Hello, world!")
	var buf bytes.Buffer
	tl.Layout()
	if _, err := tl.Write(&buf); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), testOutput1) {
		spew.Dump(buf.Bytes())
		t.Fatal(0)
	}
}

// Input 2 is the same as input 1, except little endian.  Also, IFD1 is missing
// its next IFD pointer, and there's some extra data (with an odd length) at the
// end.
var testInput2 = []byte{
	/* 0000 */ 0x49, 0x49, 0x2A, 0x00, // header, little-endian
	/* 0004 */ 0x08, 0x00, 0x00, 0x00, // pointer to IFD0
	/* 0008 */ 0x02, 0x00, // 2 tags in IFD0
	/* 000A */ 0x01, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x31, 0x32, 0x33, 0x34, // tag 1, bytes '1234'
	/* 0016 */ 0x02, 0x00, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x26, 0x00, 0x00, 0x00, // tag 2, ptr to IFD1
	/* 0022 */ 0x42, 0x00, 0x00, 0x00, // ptr to IFD2
	/* 0026 */ 0x01, 0x00, // 1 tag in IFD1
	/* 0028 */ 0x03, 0x00, 0x02, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x34, 0x00, 0x00, 0x00, // tag 3, string, "Hello, world"
	// missing next IFD pointer here
	/* 0034 */ 'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', 0x00, // tag 3 data
	/* 0041 */ 0x00, // padding
	/* 0042 */ 0x01, 0x00, // 1 tag in IFD2
	/* 0044 */ 0x04, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x00, 0x54, 0x00, 0x00, 0x00, // tag 4, rational, 2/3
	/* 0050 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0054 */ 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, // tag 4 data
	/* 005C */ 0xF1, 0xF2, 0xF3, // excess data at end
}

func TestInput2(t *testing.T) {
	var tl TIFF
	err := tl.Read(bytes.NewReader(testInput1))
	if err != nil {
		t.Fatal(err)
	}
	tag := tl.IFD0().Tag(1)
	if tag == nil {
		t.Fatal("no tag 1")
	}
	if by, err := tag.AsBytes(); err != nil {
		t.Fatalf("tag 1 %s", err)
	} else if !bytes.Equal(by, []byte("1234")) {
		t.Errorf("tag 1 wrong value")
	}
	if tag = tl.IFD0().Tag(2); tag == nil {
		t.Fatal("no tag 2")
	}
	ifd, err := tag.AsIFD()
	if err != nil {
		t.Fatalf("tag 2 %s", err)
	}
	if tag = ifd.Tag(3); tag == nil {
		t.Fatal("no tag 3")
	}
	if s, err := tag.AsString(); err != nil {
		t.Fatalf("tag 3 %s", err)
	} else if s != "Hello, world" {
		t.Errorf("tag 3 wrong value")
	}
	if ifd, err = tl.IFD0().NextIFD(); err != nil {
		t.Fatalf("next ifd %s", err)
	}
	if tag = ifd.Tag(4); tag == nil {
		t.Fatal("no tag 4")
	}
	if r, err := tag.AsRationals(); err != nil {
		t.Fatalf("tag 4 %s", err)
	} else if len(r) != 2 || r[0] != 2 || r[1] != 3 {
		t.Errorf("tag 4 wrong values")
	}
}

// The missing IFD pointer will be fixed, and the IFDs will be reordered by
// size.  The extraneous bytes won't have been moved, and a null will be added
// after them to pad to an even byte.
var testOutput2 = []byte{
	/* 0000 */ 0x49, 0x49, 0x2A, 0x00, // header, big-endian
	/* 0004 */ 0x60, 0x00, 0x00, 0x00, // pointer to IFD0
	/* 0008 */ 0x01, 0x00, // 1 tag in IFD1
	/* 000A */ 0x03, 0x00, 0x02, 0x00, 0x0F, 0x00, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, // tag 3, string, "Hello, world!!"
	/* 0016 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 001A */ 'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', '!', '!', 0x00, // tag 3 data
	/* 0029 */ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // unused space
	/* 0042 */ 0x01, 0x00, // 1 tag in IFD2
	/* 0044 */ 0x04, 0x00, 0x05, 0x00, 0x01, 0x00, 0x00, 0x00, 0x54, 0x00, 0x00, 0x00, // tag 4, rational, 2/3
	/* 0050 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0054 */ 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, // tag 4 data
	/* 005C */ 0xF1, 0xF2, 0xF3, // excess data at end
	/* 005F */ 0x00, // padding
	/* 0060 */ 0x02, 0x00, // 2 tags in IFD0
	/* 0062 */ 0x01, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x31, 0x32, 0x33, 0x34, // tag 1, bytes '1234'
	/* 006E */ 0x02, 0x00, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, // tag 2, ptr to IFD1
	/* 0072 */ 0x42, 0x00, 0x00, 0x00, // ptr to IFD2
}

func TestWrite2(t *testing.T) {
	var tl TIFF
	tl.Read(bytes.NewReader(testInput2))
	ifd0 := tl.IFD0()
	ifd1, _ := ifd0.Tag(2).AsIFD()
	tag3 := ifd1.Tag(3)
	tag3.SetString("Hello, world!!")
	var buf bytes.Buffer
	tl.Layout()
	if _, err := tl.Write(&buf); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), testOutput2) {
		spew.Dump(buf.Bytes())
		t.Fatal(0)
	}
}

var testInput3 = []byte{ // Same as test input 1, copied for convenience.
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // header, big-endian
	/* 0004 */ 0x00, 0x00, 0x00, 0x08, // pointer to IFD0
	/* 0008 */ 0x00, 0x02, // 2 tags in IFD0
	/* 000A */ 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x31, 0x32, 0x33, 0x34, // tag 1, bytes '1234'
	/* 0016 */ 0x00, 0x02, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x26, // tag 2, ptr to IFD1
	/* 0022 */ 0x00, 0x00, 0x00, 0x46, // ptr to IFD2
	/* 0026 */ 0x00, 0x01, // 1 tag in IFD1
	/* 0028 */ 0x00, 0x03, 0x00, 0x02, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x38, // tag 3, string, "Hello, world"
	/* 0034 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0038 */ 'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', 0x00, // tag 3 data
	/* 0045 */ 0x00, // padding
	/* 0046 */ 0x00, 0x01, // 1 tag in IFD2
	/* 0048 */ 0x00, 0x04, 0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x58, // tag 4, rational, 2/3
	/* 0054 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0058 */ 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, // tag 4 data
}

var testOutput3 = []byte{
	/* 0000 */ 0x4D, 0x4D, 0x00, 0x2A, // header, big-endian
	/* 0004 */ 0x00, 0x00, 0x00, 0x8E, // pointer to IFD0
	/* 0008 */ 0x00, 0x01, // 1 tag in IFD3
	/* 000A */ 0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x1A, // tag 6, 5 bytes
	/* 0016 */ 0x00, 0x00, 0x00, 0x6C, // ptr to IFD4
	/* 001A */ 0x44, 0x55, 0x66, 0x77, 0x88, // tag 6 data
	/* 001F */ 0x00, // padding
	/* 0020 */ 0, 0, 0, 0, 0, 0, // unused space
	/* 0026 */ 0x00, 0x01, // 1 tag in IFD1
	/* 0028 */ 0x00, 0x03, 0x00, 0x02, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x38, // tag 3, string, "Hello, world"
	/* 0034 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0038 */ 'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', 0x00, // tag 3 data
	/* 0045 */ 0x00, // padding
	/* 0046 */ 0x00, 0x02, // 2 tags in IFD2
	/* 0048 */ 0x00, 0x04, 0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x64, // tag 4, rational, 2/3
	/* 0054 */ 0x00, 0x05, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, // tag 5, pointer to IFD3
	/* 0060 */ 0x00, 0x00, 0x00, 0x00, // no next pointer
	/* 0064 */ 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, // tag 4 data
	/* 006C */ 0x00, 0x01, // 1 tag in IFD4
	/* 006E */ 0x00, 0x07, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x7E, // tag 7, 2 rationals
	/* 007A */ 0x00, 0x00, 0x00, 0x00, // no next ptr
	/* 007E */ 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x05, // tag 7 data
	/* 0086 */ 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x07, // tag 7 data continued
	/* 008E */ 0x00, 0x01, // 1 tags in IFD0
	/* 0090 */ 0x00, 0x02, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x26, // tag 2, ptr to IFD1
	/* 009C */ 0x00, 0x00, 0x00, 0x46, // ptr to IFD2
}

func TestWrite3(t *testing.T) {
	var tl TIFF
	tl.Read(bytes.NewReader(testInput3))
	ifd0 := tl.IFD0()
	ifd0.DeleteTag(1)
	ifd2, _ := ifd0.NextIFD()
	tag5 := ifd2.AddTag(5, 4)
	ifd3, _ := tag5.AddIFD()
	tag6 := ifd3.AddTag(6, 1)
	tag6.SetBytes([]byte{0x44, 0x55, 0x66, 0x77, 0x88})
	ifd4, _ := ifd3.AddNextIFD()
	tag7 := ifd4.AddTag(7, 5)
	tag7.SetRationals([]uint32{4, 5, 6, 7})
	var buf bytes.Buffer
	tl.Layout()
	if _, err := tl.Write(&buf); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), testOutput3) {
		spew.Dump(buf.Bytes())
		t.Fatal(0)
	}
}

func TestRangeCoalesce(t *testing.T) {
	var r rangelist
	r.add(0, 2)
	r.add(4, 6)
	r.add(2, 4)
	if len(r.r) != 2 || r.r[0] != 0 || r.r[1] != 6 {
		t.Error("fail")
	}
}
