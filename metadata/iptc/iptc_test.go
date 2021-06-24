package iptc

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/rothskeller/photo-tools/metadata"
)

var start = []byte{
	/* 0000 */ 0x38, 0x42, 0x49, 0x4D, // PSIR resource type
	/* 0004 */ 0x88, 0x99, // PSIR ID (fake)
	/* 0006 */ 0x04, 'B', 'l', 'a', 'h', 0x00, // resource name and pad
	/* 000C */ 0x00, 0x00, 0x00, 0x07, // length
	/* 0010 */ 'T', 'e', 's', 't', 'i', 'n', 'g', 0x00, // resource data and pad

	/* 0018 */ 0x38, 0x42, 0x49, 0x4D, // PSIR resource type
	/* 001C */ 0x04, 0x04, // PSIR ID for IPTC
	/* 001E */ 0x00, 0x00, // Empty name and pad
	/* 0020 */ 0x00, 0x00, 0x00, 0x10, // length
	/* 0024 */ 0x1C, 0x01, 0x5A, // data set for coded character set
	/* 0027 */ 0x00, 0x03, 0x1B, 0x25, 0x47, // UTF-8 character set
	/* 002B */ 0x1C, 0x02, 0x19, // data set for Keyword
	/* 002E */ 0x00, 0x03, 'k', 'w', '1', // Keyword

	/* 0034 */ 0x38, 0x42, 0x49, 0x4D, // PSIR resource type
	/* 0038 */ 0x88, 0xAA, // PSIR ID (fake)
	/* 003A */ 0x05, 'B', 'l', 'a', 'h', '2', // resource name
	/* 0040 */ 0x00, 0x00, 0x00, 0x08, // length
	/* 0044 */ 'T', 'e', 's', 't', 'i', 'n', 'g', '2', // resource data
}

func TestRewriteIPTC(t *testing.T) {
	iptc := Parse(start, 0)
	iptc.Keywords = []*metadata.String{metadata.NewString("new1"), metadata.NewString("new2")}
	iptc.ObjectName = metadata.NewString("me")
	out := iptc.Render()
	if !bytes.Equal(out, rewriteIPTCExpected) {
		t.Error("wrong output")
		spew.Dump(rewriteIPTCExpected)
		spew.Dump(out)
	}
}

var rewriteIPTCExpected = []byte{
	0x38, 0x42, 0x49, 0x4D, // PSIR resource type
	0x88, 0x99, // PSIR ID (fake)
	0x04, 'B', 'l', 'a', 'h', 0x00, // resource name and pad
	0x00, 0x00, 0x00, 0x07, // length
	'T', 'e', 's', 't', 'i', 'n', 'g', 0x00, // resource data and pad

	0x38, 0x42, 0x49, 0x4D, // PSIR resource type
	0x04, 0x04, // PSIR ID for IPTC
	0x00, 0x00, // Empty name and pad
	0x00, 0x00, 0x00, 0x21, // length
	0x1C, 0x01, 0x5A, // data set for coded character set
	0x00, 0x03, 0x1B, 0x25, 0x47, // UTF-8 character set
	0x1C, 0x02, 0x05, // data set for object name
	0x00, 0x02, 'm', 'e', // object name
	0x1C, 0x02, 0x19, // data set for Keyword
	0x00, 0x04, 'n', 'e', 'w', '1', // Keyword
	0x1C, 0x02, 0x19, // data set for Keyword
	0x00, 0x04, 'n', 'e', 'w', '2', // Keyword
	0x00, // padding

	0x38, 0x42, 0x49, 0x4D, // PSIR resource type
	0x04, 0x25, // PSIR ID (hash of IPTC)
	0x00, 0x00, // Empty name and pad
	0x00, 0x00, 0x00, 0x10, // length
	0x24, 0xb3, 0x09, 0x1a, 0xcd, 0x57, 0x4b, 0x06,
	0x57, 0xb1, 0xa5, 0xcc, 0xb0, 0xab, 0xf4, 0xa6, // hash

	0x38, 0x42, 0x49, 0x4D, // PSIR resource type
	0x88, 0xAA, // PSIR ID (fake)
	0x05, 'B', 'l', 'a', 'h', '2', // resource name
	0x00, 0x00, 0x00, 0x08, // length
	'T', 'e', 's', 't', 'i', 'n', 'g', '2', // resource data
}
