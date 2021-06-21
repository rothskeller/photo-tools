// STR Modifications:
// * Eliminated synchronization with dc and xmp models.
// * Kept only the fields I need.
// * Changed DateTime to a string.

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

// Package tiff implements metadata for TIFF and JPEG image files as defined
// by XMP Specification Part 1.
package tiff

import (
	"fmt"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsTiff = xmp.NewNamespace("tiff", "http://ns.adobe.com/tiff/1.0/", NewModel)
)

func init() {
	xmp.Register(NsTiff, xmp.ImageMetadata)
}

func NewModel(name string) xmp.Model {
	return &TiffInfo{}
}

func MakeModel(d *xmp.Document) (*TiffInfo, error) {
	m, err := d.MakeModel(NsTiff)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*TiffInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *TiffInfo {
	if m := d.FindModel(NsTiff); m != nil {
		return m.(*TiffInfo)
	}
	return nil
}

type TiffInfo struct {
	Artist           string        `xmp:"tiff:Artist"`
	DateTime         string        `xmp:"tiff:DateTime"`
	ImageDescription xmp.AltString `xmp:"tiff:ImageDescription"`
}

func (x TiffInfo) Can(nsName string) bool {
	return NsTiff.GetName() == nsName
}

func (x TiffInfo) Namespaces() xmp.NamespaceList {
	return []*xmp.Namespace{NsTiff}
}

func (x *TiffInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *TiffInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

// also remap X_* attributes to the correct standard positions, but don't overwrite
func (x TiffInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *TiffInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *TiffInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsTiff.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *TiffInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsTiff.GetName(), err)
	}
	return nil
}
