// Package mp implements Microsoft Photo region metadata.
package mp

import (
	"fmt"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsMP    *xmp.Namespace = xmp.NewNamespace("MP", "http://ns.microsoft.com/photo/1.2/", NewModel)
	NsMPRI  *xmp.Namespace = xmp.NewNamespace("MPRI", "http://ns.microsoft.com/photo/1.2/t/RegionInfo#", NewModel)
	NsMPReg *xmp.Namespace = xmp.NewNamespace("MPReg", "http://ns.microsoft.com/photo/1.2/t/Region#", NewModel)
)

func init() {
	xmp.Register(NsMP, xmp.XmpMetadata)
	xmp.Register(NsMPRI, xmp.XmpMetadata)
	xmp.Register(NsMPReg, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	println(name)
	return &MPInfo{}
}

func MakeModel(d *xmp.Document) (*MPInfo, error) {
	m, err := d.MakeModel(NsMP)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*MPInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *MPInfo {
	if m := d.FindModel(NsMP); m != nil {
		return m.(*MPInfo)
	}
	return nil
}

type MPInfo struct {
	RegionInfo RegionInfo `xmp:"MP:RegionInfo"`
}

func (x MPInfo) Can(nsName string) bool {
	return NsMP.GetName() == nsName
}

func (x MPInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsMP}
}

func (x *MPInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *MPInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x MPInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *MPInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *MPInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsMP.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *MPInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsMP.GetName(), err)
	}
	return nil
}
