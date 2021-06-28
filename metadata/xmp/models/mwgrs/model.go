// Package mwgrs implements Metadata Working Group Region Structure metadata.
package mwgrs

import (
	"fmt"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsMwgRs = xmp.NewNamespace("mwg-rs", "http://www.metadataworkinggroup.com/schemas/regions/", NewModel)
)

func init() {
	xmp.Register(NsMwgRs, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &MWGRegions{}
}

func MakeModel(d *xmp.Document) (*MWGRegions, error) {
	m, err := d.MakeModel(NsMwgRs)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*MWGRegions)
	return x, nil
}

func FindModel(d *xmp.Document) *MWGRegions {
	if m := d.FindModel(NsMwgRs); m != nil {
		return m.(*MWGRegions)
	}
	return nil
}

type MWGRegions struct {
	Regions RegionInfo `xmp:"mwg-rs:Regions"`
}

func (x MWGRegions) Can(nsName string) bool {
	return NsMwgRs.GetName() == nsName
}

func (x MWGRegions) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsMwgRs}
}

func (x *MWGRegions) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *MWGRegions) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x MWGRegions) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *MWGRegions) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *MWGRegions) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsMwgRs.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *MWGRegions) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsMwgRs.GetName(), err)
	}
	return nil
}
