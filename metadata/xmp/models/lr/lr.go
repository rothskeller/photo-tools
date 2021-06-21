// Package lr defines the Lightroom XMP data model, since it is not provided by
// go-xmp.
package lr

import (
	"fmt"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsLr = xmp.NewNamespace("lr", "http://ns.adobe.com/lightroom/1.0/", NewModel)
)

func init() {
	xmp.Register(NsLr, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &Lightroom{}
}

func MakeModel(d *xmp.Document) (*Lightroom, error) {
	m, err := d.MakeModel(NsLr)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*Lightroom)
	return x, nil
}

func FindModel(d *xmp.Document) *Lightroom {
	if m := d.FindModel(NsLr); m != nil {
		return m.(*Lightroom)
	}
	return nil
}

type Lightroom struct {
	HierarchicalSubject xmp.StringArray `xmp:"lr:hierarchicalSubject"`
}

func (x Lightroom) Can(nsName string) bool {
	return NsLr.GetName() == nsName
}

func (x Lightroom) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsLr}
}

func (x *Lightroom) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *Lightroom) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x Lightroom) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *Lightroom) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *Lightroom) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsLr.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *Lightroom) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsLr.GetName(), err)
	}
	return nil
}
