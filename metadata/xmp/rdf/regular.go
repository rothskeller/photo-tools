package rdf

import (
	"fmt"

	"github.com/beevik/etree"
)

func regularize(elm *etree.Element) error {
	if err := rejectRDFType(elm); err != nil {
		return err
	}
	if err := expandParseTypeResource(elm); err != nil {
		return err
	}
	if err := expandAttributes(elm); err != nil {
		return err
	}
	return nil
}

// rejectRDFType searches for rdf:type attributes, and raises an error if it
// finds any.  They are allowed by the XMP/RDF spec, but they are too complex to
// be worth handling (especially since I've never seen one in the real world).
func rejectRDFType(elm *etree.Element) error {
	for _, attr := range elm.Attr {
		if attr.Space == NSrdf && attr.Key == "type" {
			return fmt.Errorf("%s: rdf:type: unsupported attribute", elm.FullTag())
		}
	}
	for _, child := range elm.Child {
		if child, ok := child.(*etree.Element); ok {
			if err := rejectRDFType(child); err != nil {
				return err
			}
		}
	}
	return nil
}

// expandParseTypeResource finds all rdf:parseType="Resource" attributes.
// Whenever it finds one, it inserts an rdf:Description element under the
// element with that attribute, and moves all of that elements children under
// the rdf:Description.  It also removes the rdf:parseType attribute from the
// parent.  It returns an error if it finds an rdf:parseType with any value
// other than "Resource".
func expandParseTypeResource(elm *etree.Element) error {
	var found = -1
	for idx, attr := range elm.Attr {
		if attr.Space == NSrdf && attr.Key == "parseType" {
			if attr.Value != "Resource" {
				return fmt.Errorf("%s: rdf:%s: unsupported value %q", elm.FullTag(), attr.Key, attr.Value)
			}
			found = idx
			break
		}
	}
	if found >= 0 {
		var desc = etree.NewElement("Description")
		desc.Space = NSrdf
		for _, child := range elm.Child {
			desc.AddChild(child)
		}
		elm.AddChild(desc)
		elm.Attr = append(elm.Attr[0:found], elm.Attr[found+1:]...)
	}
	for _, child := range elm.Child {
		if child, ok := child.(*etree.Element); ok {
			if err := expandParseTypeResource(child); err != nil {
				return err
			}
		}
	}
	return nil
}

// expandAttributes finds all attributes and changes them into child elements
// with text values.  If the target element is not an rdf:Description, one is
// created underneath the target, and the child elements are placed under it.
//
// xml:lang is handled specially.  If the target element contains an
// rdf:Description, which in turn contains an rdf:value, the element for
// xml:lang is added under the rdf:Description.  Otherwise, a new
// rdf:Description is created under the target, with the element for the
// xml:lang plus an rdf:value element, and the previous children of the target
// are put under the rdf:value.
func expandAttributes(elm *etree.Element) error {
	var lang etree.Attr
	var desc *etree.Element
	if elm.Space == NSrdf && elm.Tag == "Description" {
		desc = elm
	}
	for _, attr := range elm.Attr {
		if attr.Space == NSxml && attr.Key == "lang" {
			if lang.Key != "" {
				return fmt.Errorf("%s: multiple values for xml:lang", elm.FullTag())
			}
			lang = attr
			continue
		}
		if desc == nil {
			desc = etree.NewElement("Description")
			desc.Space = NSrdf
			for _, child := range elm.Child {
				desc.AddChild(child)
			}
			elm.AddChild(desc)
		}
		var aelm = etree.NewElement(attr.Key)
		aelm.Space = attr.Space
		aelm.CreateText(attr.Value)
		desc.AddChild(aelm)
	}
	elm.Attr = nil
	if lang.Key != "" {
		var desc *etree.Element
		if len(elm.Child) == 1 {
			if child, ok := elm.Child[0].(*etree.Element); ok && child.Space == NSrdf && child.Tag == "Description" {
				desc = child
			}
		}
		if desc != nil {
			var found = false
			for _, child := range elm.Child {
				if child, ok := child.(*etree.Element); ok && child.Space == NSrdf && child.Tag == "value" {
					found = true
					break
				}
			}
			if !found {
				desc = nil
			}
		}
		if desc == nil {
			var value = etree.NewElement("value")
			value.Space = NSrdf
			for _, child := range elm.Child {
				value.AddChild(child)
			}
			desc = etree.NewElement("Description")
			desc.Space = NSrdf
			desc.AddChild(value)
			elm.AddChild(desc)
		}
		var aelm = etree.NewElement("lang")
		aelm.Space = NSxml
		aelm.CreateText(lang.Value)
		desc.AddChild(aelm)
	}
	for _, child := range elm.Child {
		if child, ok := child.(*etree.Element); ok {
			if err := expandAttributes(child); err != nil {
				return err
			}
		}
	}
	return nil
}
