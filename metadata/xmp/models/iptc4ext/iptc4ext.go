// Package iptc4ext defines (parts of) the IPTC Extension XMP data model, since
// it is not provided by go-xmp.
package iptc4ext

import (
	"fmt"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsIptc4xmpExt = xmp.NewNamespace("Iptc4xmpExt", "http://iptc.org/std/Iptc4xmpExt/2008-02-29/", NewModel)
)

func init() {
	xmp.Register(NsIptc4xmpExt, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &Iptc4xmpExt{}
}

func MakeModel(d *xmp.Document) (*Iptc4xmpExt, error) {
	m, err := d.MakeModel(NsIptc4xmpExt)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*Iptc4xmpExt)
	return x, nil
}

func FindModel(d *xmp.Document) *Iptc4xmpExt {
	if m := d.FindModel(NsIptc4xmpExt); m != nil {
		return m.(*Iptc4xmpExt)
	}
	return nil
}

type Iptc4xmpExt struct {
	LocationCreated *Location     `xmp:"Iptc4xmpExt:LocationCreated"`
	LocationShown   LocationArray `xmp:"Iptc4xmpExt:LocationShown"`
}

func (x Iptc4xmpExt) Can(nsName string) bool {
	return NsIptc4xmpExt.GetName() == nsName
}

func (x Iptc4xmpExt) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsIptc4xmpExt}
}

func (x *Iptc4xmpExt) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *Iptc4xmpExt) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x Iptc4xmpExt) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *Iptc4xmpExt) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *Iptc4xmpExt) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsIptc4xmpExt.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *Iptc4xmpExt) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsIptc4xmpExt.GetName(), err)
	}
	return nil
}

type LocationArray []*Location

func (x LocationArray) Typ() xmp.ArrayType {
	return xmp.ArrayTypeOrdered
}

func (x LocationArray) MarshalXMP(e *xmp.Encoder, node *xmp.Node, m xmp.Model) error {
	return xmp.MarshalArray(e, node, x.Typ(), x)
}

func (x *LocationArray) UnmarshalXMP(d *xmp.Decoder, node *xmp.Node, m xmp.Model) error {
	return xmp.UnmarshalArray(d, node, x.Typ(), x)
}

type Location struct {
	City          xmp.AltString `xmp:"Iptc4xmpExt:City"`
	CountryCode   string        `xmp:"Iptc4xmpExt:CountryCode"`
	CountryName   xmp.AltString `xmp:"Iptc4xmpExt:CountryName"`
	ProvinceState xmp.AltString `xmp:"Iptc4xmpExt:ProvinceState"`
	Sublocation   xmp.AltString `xmp:"Iptc4xmpExt:Sublocation"`
}
