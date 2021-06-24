// Package digikam defines the digiKam XMP data model, since it is not provided
// by go-xmp.
package digikam

import (
	"fmt"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsDigiKam = xmp.NewNamespace("digiKam", "http://www.digikam.org/ns/1.0/", NewModel)
)

func init() {
	xmp.Register(NsDigiKam, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &DigiKam{}
}

func MakeModel(d *xmp.Document) (*DigiKam, error) {
	m, err := d.MakeModel(NsDigiKam)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*DigiKam)
	return x, nil
}

func FindModel(d *xmp.Document) *DigiKam {
	if m := d.FindModel(NsDigiKam); m != nil {
		return m.(*DigiKam)
	}
	return nil
}

type DigiKam struct {
	TagsList xmp.StringList `xmp:"digiKam:TagsList"`
}

func (x DigiKam) Can(nsName string) bool {
	return NsDigiKam.GetName() == nsName
}

func (x DigiKam) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsDigiKam}
}

func (x *DigiKam) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *DigiKam) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x DigiKam) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *DigiKam) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *DigiKam) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsDigiKam.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *DigiKam) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsDigiKam.GetName(), err)
	}
	return nil
}
