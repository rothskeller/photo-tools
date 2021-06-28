package mwgrs

import (
	"trimmer.io/go-xmp/xmp"
)

type RegionInfo struct {
	RegionList          RegionList `xmp:"mwg-rs:RegionList"`
	AppliedToDimensions string     `xmp:"mwg-rs:AppliedToDimensions,omit"`
}

type RegionList []RegionStruct

func (x *RegionList) UnmarshalXMP(d *xmp.Decoder, n *xmp.Node, model xmp.Model) error {
	return xmp.UnmarshalArray(d, n, xmp.ArrayTypeUnordered, x)
}

type RegionStruct struct {
	Type string `xmp:"mwg-rs:Type"`
	Name string `xmp:"mwg-rs:Name"`
	Area string `xmp:"mwg-rs:Area,omit"`
}
