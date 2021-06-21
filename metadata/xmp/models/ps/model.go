// STR Modifications:
// * Eliminated everything except the DateCreated tag, and made that a string.

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

// see also Photoshop Specification at
// http://www.adobe.com/devnet-apps/photoshop/fileformatashtml

// Package ps implements Adobe Photoshop metadata as defined by XMP Specification Part 2 Chapter 3.2.
package ps

import (
	"fmt"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsPhotoshop = xmp.NewNamespace("photoshop", "http://ns.adobe.com/photoshop/1.0/", NewModel)
)

func init() {
	xmp.Register(NsPhotoshop, xmp.ImageMetadata)
}

func NewModel(name string) xmp.Model {
	return &PhotoshopInfo{}
}

func MakeModel(d *xmp.Document) (*PhotoshopInfo, error) {
	m, err := d.MakeModel(NsPhotoshop)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*PhotoshopInfo)
	return x, nil
}

func FindModel(d *xmp.Document) *PhotoshopInfo {
	if m := d.FindModel(NsPhotoshop); m != nil {
		return m.(*PhotoshopInfo)
	}
	return nil
}

type PhotoshopInfo struct {
	DateCreated string `xmp:"photoshop:DateCreated"`
}

func (x PhotoshopInfo) Can(nsName string) bool {
	return NsPhotoshop.GetName() == nsName
}

func (x PhotoshopInfo) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsPhotoshop}
}

func (x *PhotoshopInfo) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *PhotoshopInfo) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x PhotoshopInfo) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *PhotoshopInfo) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *PhotoshopInfo) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsPhotoshop.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *PhotoshopInfo) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsPhotoshop.GetName(), err)
	}
	return nil
}
