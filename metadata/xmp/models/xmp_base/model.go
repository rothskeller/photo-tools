// STR Modifications:
// * Removed synchronization with other models.
// * Removed all except the fields I need.

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

// Package xmpbase implements the XMP namespace as defined by XMP Specification Part 2.
package xmpbase

import (
	"fmt"
	"strings"

	"trimmer.io/go-xmp/xmp"
)

var (
	NsXmp = xmp.NewNamespace("xmp", "http://ns.adobe.com/xap/1.0/", NewModel)
)

func init() {
	xmp.Register(NsXmp, xmp.XmpMetadata)
}

func NewModel(name string) xmp.Model {
	return &XmpBase{}
}

func MakeModel(d *xmp.Document) (*XmpBase, error) {
	m, err := d.MakeModel(NsXmp)
	if err != nil {
		return nil, err
	}
	x, _ := m.(*XmpBase)
	return x, nil
}

func FindModel(d *xmp.Document) *XmpBase {
	if m := d.FindModel(NsXmp); m != nil {
		return m.(*XmpBase)
	}
	return nil
}

type XmpBase struct {
	CreateDate   string `xmp:"xmp:CreateDate"`
	MetadataDate string `xmp:"xmp:MetadataDate"`
	ModifyDate   string `xmp:"xmp:ModifyDate"`
}

func (x XmpBase) Can(nsName string) bool {
	return NsXmp.GetName() == nsName
}

func (x XmpBase) Namespaces() xmp.NamespaceList {
	return xmp.NamespaceList{NsXmp}
}

func prefixer(prefix string) xmp.ConverterFunc {
	return func(val string) string {
		return strings.Join([]string{prefix, val}, ":")
	}
}

func (x *XmpBase) SyncModel(d *xmp.Document) error {
	return nil
}

func (x *XmpBase) SyncFromXMP(d *xmp.Document) error {
	return nil
}

func (x XmpBase) SyncToXMP(d *xmp.Document) error {
	return nil
}

func (x *XmpBase) CanTag(tag string) bool {
	_, err := xmp.GetNativeField(x, tag)
	return err == nil
}

func (x *XmpBase) GetTag(tag string) (string, error) {
	if v, err := xmp.GetNativeField(x, tag); err != nil {
		return "", fmt.Errorf("%s: %v", NsXmp.GetName(), err)
	} else {
		return v, nil
	}
}

func (x *XmpBase) SetTag(tag, value string) error {
	if err := xmp.SetNativeField(x, tag, value); err != nil {
		return fmt.Errorf("%s: %v", NsXmp.GetName(), err)
	}
	return nil
}
