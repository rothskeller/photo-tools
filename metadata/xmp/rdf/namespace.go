package rdf

import (
	"fmt"

	"github.com/beevik/etree"
)

// expandNamespaces walks through the element tree, changing the Space field in
// all elements and attributes to the full URI of the namespace referenced
// there.  It removes all xmlns attributes from the tree.  It returns a map from
// namespace URI to the prefix used to reference that namespace.  (If more than
// one prefix is used to reference the same namespace, the one used latest in
// the document is returned.)  It returns an error if it finds any unregistered
// namespace prefix or any element or attribute without a namespace.
func expandNamespaces(elm *etree.Element) (prefixes map[string]string, err error) {
	var nsuris = map[string]string{"xml": NSxml}
	var nsprefixes = map[string]string{NSxml: "xml"}
	for _, child := range elm.Child {
		if child, ok := child.(*etree.Element); ok {
			if err := expandNamespaces1(child, nsuris, nsprefixes); err != nil {
				return nil, err
			}
		}
	}
	return nsprefixes, nil
}
func expandNamespaces1(elm *etree.Element, nsuris, nsprefixes map[string]string) error {
	var save = make(map[string]string)
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			if nsuris[attr.Key] != "" {
				save[attr.Key] = nsuris[attr.Key]
			}
			nsuris[attr.Key] = attr.Value
			nsprefixes[attr.Value] = attr.Key
		}
	}
	for i, attr := range elm.Attr {
		if attr.Space != "xmlns" {
			if elm.Attr[i].Space = nsuris[elm.Attr[i].Space]; elm.Attr[i].Space == "" {
				return fmt.Errorf("%s: %s: unregistered namespace", elm.FullTag(), attr.FullKey())
			}
		}
	}
	if nsuri := nsuris[elm.Space]; nsuri != "" {
		elm.Space = nsuri
	} else {
		return fmt.Errorf("%s: unregistered namespace", elm.FullTag())
	}
	for _, child := range elm.Child {
		if child, ok := child.(*etree.Element); ok {
			if err := expandNamespaces1(child, nsuris, nsprefixes); err != nil {
				return err
			}
		}
	}
	for i := 0; i < len(elm.Attr); {
		var attr = elm.Attr[i]
		if attr.Space == "xmlns" {
			delete(nsuris, attr.Key)
			elm.RemoveAttr(attr.FullKey())
		} else {
			i++
		}
	}
	for key, uri := range save {
		nsuris[key] = uri
	}
	return nil
}

// prefixNamespaces walks through the element tree, replacing the namespace URIs
// in the Space entries in elements and attributes with the assigned prefixes
// for those namespaces.  It returns an error if any namespaces are used that
// don't have an assigned prefix.  It adds xmlns attributes to the root element
// for all namespaces that were actually used.
func prefixNamespaces(elm *etree.Element, nsprefixes map[string]string) error {
	var used = make(map[string]bool)
	for _, prefix := range nsprefixes {
		if used[prefix] {
			return fmt.Errorf("prefix %q assigned to multiple namespaces", prefix)
		}
		used[prefix] = true
	}
	used = make(map[string]bool)
	if err := prefixNamespaces1(elm, nsprefixes, used); err != nil {
		return err
	}
	for uri := range used {
		prefix := nsprefixes[uri]
		if uri == NSxml && prefix == "xml" {
			continue
		}
		elm.CreateAttr("xmlns:"+prefix, uri)
	}
	return nil
}
func prefixNamespaces1(elm *etree.Element, nsprefixes map[string]string, used map[string]bool) error {
	used[elm.Space] = true
	prefix := nsprefixes[elm.Space]
	if prefix == "" {
		return fmt.Errorf("no assigned prefix for namespace %s", elm.Space)
	}
	elm.Space = prefix
	for i, attr := range elm.Attr {
		used[attr.Space] = true
		prefix = nsprefixes[attr.Space]
		if prefix == "" {
			return fmt.Errorf("no assigned prefix for namespace %s", attr.Space)
		}
		elm.Attr[i].Space = prefix
	}
	for _, child := range elm.Child {
		if child, ok := child.(*etree.Element); ok {
			if err := prefixNamespaces1(child, nsprefixes, used); err != nil {
				return err
			}
		}
	}
	return nil
}
