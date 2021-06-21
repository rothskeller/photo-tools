// STR Modifications:
// * Reduced ExifInfo down to only the fields I need.
// * Changed date-typed fields to strings.
// * Changed GPS fields to strings.
// * Remove ExifEx and ExifAux, not needed.

// Copyright (c) 2017-2018 Alexander Eichhorn
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// XMP EXIF Mapping for Exif 2.3 metadata CIPA DC-010-2012
//
// Exif Spec
// Exif 2.3 http://www.cipa.jp/std/documents/e/DC-010-2012_E.pdf
// Exif 2.3 http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf
// Exif 2.3.1 http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
//
// see https://www.media.mit.edu/pia/Research/deepview/exif.html
// for a very good explanation of tags and ifd's

// Package exif implements the Exif 2.3.1 metadata standard as defined in CIPA DC-008-2016.
package exif

import (
	"fmt"
	"strings"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsExif *xmp.Namespace    = xmp.NewNamespace("exif", "http://ns.adobe.com/exif/1.0/", NewModel)
	nslist xmp.NamespaceList = xmp.NamespaceList{NsExif}
)

func init() {
	for _, v := range nslist {
		xmp.Register(v, xmp.ImageMetadata)
	}
}

func NewModel(name string) xmp.Model {
	switch name {
	case "exif":
		return &ExifInfo{}
	}
	return nil
}

func MakeModel(d *xmp.Document) (*ExifInfo, error) {
	m, err := d.MakeModel(NsExif)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*ExifInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *ExifInfo {
	if m := d.FindModel(NsExif); m != nil {
		return m.(*ExifInfo)
	}
	return nil
}

type ExifInfo struct {
	PixelXDimension   int             `exif:"0xa002" xmp:"exif:PixelXDimension"`
	PixelYDimension   int             `exif:"0xa003" xmp:"exif:PixelYDimension"`
	UserComment       xmp.StringArray `exif:"0x9286" xmp:"exif:UserComment"`
	DateTimeOriginal  string          `exif:"0x9003" xmp:"exif:DateTimeOriginal"`
	DateTimeDigitized string          `exif:"0x9004" xmp:"exif:DateTimeDigitized"`
	GPSLatitude       string          `exif:"-"      xmp:"exif:GPSLatitude"`
	GPSLongitude      string          `exif:"-"      xmp:"exif:GPSLongitude"`
	GPSAltitudeRef    string          `exif:"0x0005" xmp:"exif:GPSAltitudeRef"`
	GPSAltitude       string          `exif:"0x0006" xmp:"exif:GPSAltitude"`
}

func (m *ExifInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsExif}
}

func (m *ExifInfo) Can(nsName string) bool {
	return nsName == NsExif.GetName()
}

func (x *ExifInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *ExifInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x *ExifInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *ExifInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *ExifInfo) GetTag(tag string) (string, error) {
	tag = strings.ToLower(tag)
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("exif: %v", err)
	} else {
		return v, nil
	}
}

func (x *ExifInfo) SetTag(tag, value string) error {
	tag = strings.ToLower(tag)
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("exif: %v", err)
	}
	return nil
}

func (x *ExifInfo) GetLocaleTag(lang string, tag string) (string, error) {
	tag = strings.ToLower(tag)
	if val, err := xmp.GetLocaleField(x, lang, tag); err != nil {
		return "", fmt.Errorf("exif: %v", err)
	} else {
		return val, nil
	}
}

func (x *ExifInfo) SetLocaleTag(lang string, tag, value string) error {
	tag = strings.ToLower(tag)
	if err := xmp.SetLocaleField(x, lang, tag, value); err != nil {
		return fmt.Errorf("exif: %v", err)
	}
	return nil
}
