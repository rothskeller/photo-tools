package rdf

import (
	"fmt"
	"sort"

	"github.com/beevik/etree"
)

// Render renders the packet as encoded XML.
func (p *Packet) Render() (out []byte, err error) {
	p.nsprefixes[NSrdf] = "rdf"
	p.nsprefixes[NSxml] = "xml"
	var doc = etree.NewDocument()
	doc.CreateProcInst("xpacket", `begin="" id="W5M0MpCehiHzreSzNTczkc9d"`)
	xmpmeta := doc.CreateElement("x:xmpmeta")
	xmpmeta.CreateAttr("xmlns:x", NSx)
	root := xmpmeta.CreateElement("rdf:RDF")
	if err = p.renderNamespaces(root); err != nil {
		return nil, fmt.Errorf("RDF: %s", err)
	}
	desc := p.renderStruct(root, p.Properties, true)
	desc.CreateAttr("rdf:about", p.about)
	doc.CreateProcInst("xpacket", `end="w"`)
	if out, err = doc.WriteToBytes(); err != nil {
		return nil, err
	}
	return out, nil
}

// renderNamespaces adds xmlns attributes to the root element for each namespace
// actually used in the RDF block.  It returns an error if any namespace lacks a
// prefix or if multiple namespaces use the same prefix.
func (p *Packet) renderNamespaces(root *etree.Element) error {
	var nsuris = map[string]string{"rdf": NSrdf}
	if err := p.renderNamespacesStruct(p.Properties, nsuris); err != nil {
		return err
	}
	delete(nsuris, "xml") // shouldn't emit an xmlns for it
	var prefixes = make([]string, 0, len(nsuris))
	for prefix := range nsuris {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)
	for _, prefix := range prefixes {
		root.CreateAttr("xmlns:"+prefix, nsuris[prefix])
	}
	return nil
}
func (p *Packet) renderNamespacesStruct(str Struct, nsuris map[string]string) error {
	for key, value := range str {
		prefix := p.nsprefixes[key.Namespace]
		if prefix == "" {
			return fmt.Errorf("no XML prefix for %s", key.Namespace)
		}
		if exist := nsuris[prefix]; exist != "" && exist != key.Namespace {
			return fmt.Errorf("multiple namespaces with prefix %q", prefix)
		}
		nsuris[prefix] = key.Namespace
		if err := p.renderNamespacesValue(value, nsuris); err != nil {
			return err
		}
	}
	return nil
}
func (p *Packet) renderNamespacesValue(value Value, nsuris map[string]string) error {
	if err := p.renderNamespacesStruct(value.Qualifiers, nsuris); err != nil {
		return err
	}
	switch value := value.Value.(type) {
	case Alt:
		for _, v := range value {
			if err := p.renderNamespacesValue(v, nsuris); err != nil {
				return err
			}
		}
	case Bag:
		for _, v := range value {
			if err := p.renderNamespacesValue(v, nsuris); err != nil {
				return err
			}
		}
	case Seq:
		for _, v := range value {
			if err := p.renderNamespacesValue(v, nsuris); err != nil {
				return err
			}
		}
	case Struct:
		if err := p.renderNamespacesStruct(value, nsuris); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packet) renderStruct(elm *etree.Element, str Struct, forceDesc bool) *etree.Element {
	var names = make([]Name, 0, len(str))
	var hasSimple, hasComplex bool
	for name, value := range str {
		names = append(names, name)
		if _, ok := value.Value.(string); !ok || len(value.Qualifiers) != 0 {
			hasComplex = true
		} else {
			hasSimple = true
		}
	}
	sort.Slice(names, func(i, j int) bool {
		ipfx := p.nsprefixes[names[i].Namespace]
		jpfx := p.nsprefixes[names[j].Namespace]
		if ipfx != jpfx {
			return ipfx < jpfx
		}
		return names[i].Name < names[j].Name
	})
	if (hasSimple && hasComplex) || forceDesc {
		elm = elm.CreateElement("rdf:Description")
	} else if hasComplex && !hasSimple {
		elm.CreateAttr("rdf:parseType", "Resource")
	}
	for _, name := range names {
		p.renderValue(elm, name, str[name], true)
	}
	return elm
}

func (p *Packet) renderValue(elm *etree.Element, name Name, value Value, canAttr bool) {
	var quals = value.Qualifiers

	if len(quals) == 0 && canAttr {
		if value, ok := value.Value.(string); ok {
			elm.CreateAttr(p.prefixedName(name), value)
			return
		}
	}
	elm = elm.CreateElement(p.prefixedName(name))
	if langv, ok := quals[Name{NSxml, "lang"}]; ok && len(quals) == 1 {
		if lang, ok := langv.Value.(string); ok {
			elm.CreateAttr("xml:lang", lang)
			quals = nil
		}
	}
	if len(quals) != 0 {
		p.renderStruct(elm, quals, false)
		p.renderValue(elm, Name{NSrdf, "value"}, Value{Value: value.Value}, true)
		return
	}
	switch value := value.Value.(type) {
	case string:
		elm.CreateText(value)
	case URI:
		elm.CreateAttr("rdf:resource", string(value))
	case Alt:
		p.renderArray(elm.CreateElement("rdf:Alt"), value)
	case Bag:
		p.renderArray(elm.CreateElement("rdf:Bag"), value)
	case Seq:
		p.renderArray(elm.CreateElement("rdf:Seq"), value)
	case Struct:
		p.renderStruct(elm, value, false)
	}
}

func (p *Packet) renderArray(elm *etree.Element, values []Value) {
	for _, value := range values {
		p.renderValue(elm, Name{NSrdf, "li"}, value, false)
	}
}

func (p *Packet) prefixedName(name Name) string {
	return fmt.Sprintf("%s:%s", p.nsprefixes[name.Namespace], name.Name)
}
