package mp

import "trimmer.io/go-xmp/xmp"

type RegionInfo struct {
	Regions RegionList `xmp:"MPRI:Regions"`
}

type RegionList []Region

func (x *RegionList) UnmarshalXMP(d *xmp.Decoder, n *xmp.Node, model xmp.Model) error {
	return xmp.UnmarshalArray(d, n, xmp.ArrayTypeUnordered, x)
}

type Region struct {
	PersonDisplayName string `xmp:"MPReg:PersonDisplayName"`
	Rectangle         string `xmp:"MPReg:Rectangle"`
}
