package rdf

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

const (
	nsRDF = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	nsX   = "adobe:ns:meta/"
	nsXML = "http://www.w3.org/XML/1998/namespace"
)

var startPIRE = regexp.MustCompile(`^begin=['"]\x{FEFF}?['"] id=['"]W5M0MpCehiHzreSzNTczkc9d['"]`)
var endPIRE = regexp.MustCompile(`^end=['"][wr]['"]$`)

// NewPacket creates a new, empty RDF packet.
func NewPacket() *Packet {
	return &Packet{
		properties: make(Struct),
		nsprefixes: map[string]string{nsXML: "xml"},
		nsuris:     map[string]string{"xml": nsXML},
	}
}

// ReadPacket parses the supplied buffer and returns the resulting Packet, or an
// error if the packet is invalid.
func ReadPacket(buf []byte) (p *Packet, err error) {
	var (
		doc      *etree.Document
		children []*etree.Element
		root     *etree.Element
	)
	p = NewPacket()
	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(buf); err != nil {
		return nil, err
	}
	if children, _, err = parseElement(&doc.Element); err != nil {
		return nil, err
	}
	switch len(children) {
	case 0:
		return nil, errors.New("XMP packet has no root element")
	case 1:
		break
	default:
		return nil, errors.New("XMP packet has multiple root elements")
	}
	root = children[0]
	p.pushElementNamespaces(root)
	if p.nsuris[root.Space] == nsX && root.Tag == "xmpmeta" {
		if children, _, err = parseElement(root); err != nil {
			return nil, err
		}
		switch len(children) {
		case 0:
			return nil, errors.New("x:xmpmeta has no child element")
		case 1:
			break
		default:
			return nil, errors.New("x:xmpmeta has multiple child elements")
		}
		root = children[0]
		p.pushElementNamespaces(root)
	}
	if p.nsuris[root.Space] != nsRDF || root.Tag != "RDF" {
		return nil, fmt.Errorf("XMP packet root element is %s:%s, not rdf:RDF", root.Space, root.Tag)
	}
	if children, err = elmChildren(root); err != nil {
		return nil, err
	}
	for _, child := range children {
		p.pushElementNamespaces(child)
		if p.nsuris[child.Space] != nsRDF || child.Tag != "Description" {
			return nil, fmt.Errorf("rdf:RDF contains unexpected %s:%s element", child.Space, child.Tag)
		}
		if err = p.readPropertyDescription(child); err != nil {
			return nil, err
		}
		p.popElementNamespaces(child)
	}
	return p, nil
}

// readPropertyDescription handles a top-level rdf:Description element.
func (p *Packet) readPropertyDescription(elm *etree.Element) (err error) {
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			continue
		}
		var nsuri = p.nsuris[attr.Space]
		if nsuri == "" {
			return fmt.Errorf("attribute %s on rdf:Description has unregistered namespace", attr.FullKey())
		}
		if nsuri == nsRDF && attr.Key == "about" {
			if p.about != "" && attr.Value != "" && p.about != attr.Value {
				return errors.New("mismatched values for rdf:about")
			}
			if attr.Value != "" {
				p.about = attr.Value
			}
			continue
		}
		if nsuri == nsRDF || nsuri == nsXML {
			return fmt.Errorf("unexpected attribute %s:%s on rdf:Description", attr.Space, attr.Key)
		}
		var name = Name{nsuri, attr.Key}
		if _, ok := p.properties[name]; ok {
			return fmt.Errorf("multiple values for %s:%s", attr.Space, attr.Key)
		}
		p.properties[name] = Value{Value: attr.Value}
	}
	var children []*etree.Element
	if children, err = elmChildren(elm); err != nil {
		return err
	}
	for _, child := range children {
		p.pushElementNamespaces(child)
		if err = p.readProperty(child); err != nil {
			return err
		}
		p.popElementNamespaces(child)
	}
	return nil
}

// readProperty handles a property definition element, as an immediate child of
// a top-level rdf:Description element.
func (p *Packet) readProperty(elm *etree.Element) (err error) {
	var nsuri = p.nsuris[elm.Space]
	if nsuri == "" {
		return fmt.Errorf("element %s has unregistered namespace", elm.FullTag())
	}
	if nsuri == nsRDF || nsuri == nsXML {
		return fmt.Errorf("unexpected element %s:%s", elm.Space, elm.Tag)
	}
	var name = Name{nsuri, elm.Tag}
	if _, ok := p.properties[name]; ok {
		return fmt.Errorf("multiple values for %s:%s", elm.Space, elm.Tag)
	}
	var value Value
	if value, err = p.readValueElm(elm, true); err != nil {
		return err
	}
	p.properties[name] = value
	return nil
}

// readValueElm reads a property, structure field, or qualifier Value from an
// element.
func (p *Packet) readValueElm(elm *etree.Element, qualOK bool) (value Value, err error) {
	var (
		parseResource bool
		rdfResource   etree.Attr
		otherattrs    []etree.Attr
		children      []*etree.Element
		chardata      string
	)
	// Walk through the list of attributes and pull out the interesting ones.
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			continue
		}
		switch nsuri := p.nsuris[attr.Space]; nsuri {
		case nsRDF:
			switch attr.Key {
			case "resource":
				rdfResource = attr
			case "ParseType":
				if attr.Value != "Resource" {
					return Value{}, fmt.Errorf("%s:%s has unsupported rdf:ParseType attribute", elm.Space, elm.Tag)
				}
				parseResource = true
			case "value":
				otherattrs = append(otherattrs, attr) // handle below
			default:
				return Value{}, fmt.Errorf("%s:%s has unsupported %s:%s attribute", elm.Space, elm.Tag, attr.Space, attr.Key)
			}
		case nsXML:
			if attr.Key != "lang" {
				return Value{}, fmt.Errorf("%s:%s has unsupported %s:%s attribute", elm.Space, elm.Tag, attr.Space, attr.Key)
			}
			value.Qualifiers = append(value.Qualifiers, Qualifier{Name{nsXML, "lang"}, Value{Value: attr.Value}})
		case "":
			return Value{}, fmt.Errorf("%s:%s has attribute %s with unregistered namespace", elm.Space, elm.Tag, attr.FullKey())
		default:
			otherattrs = append(otherattrs, attr)
		}
	}
	if children, chardata, err = parseElement(elm); err != nil {
		return Value{}, err
	}
	if rdfResource.Key != "" {
		// Case 1: rdf:resource attribute, no data (URI value)
		if parseResource || len(otherattrs) != 0 || len(children) != 0 || strings.TrimSpace(chardata) != "" {
			return Value{}, fmt.Errorf("%s:%s has rdf:resource and conflicting other content", elm.Space, elm.Tag)
		}
		value.Value = URI(rdfResource.Value)
	} else if parseResource {
		// Case 2: rdf:ParseType="Resource" attribute, multiple child elements (structure)
		if len(otherattrs) != 0 || strings.TrimSpace(chardata) != "" {
			return Value{}, fmt.Errorf("%s:%s has rdf:ParseType=\"Resource\" and conflicting other content", elm.Space, elm.Tag)
		}
		if value.Value, err = p.readStructFieldsElm(elm); err != nil {
			return Value{}, err
		}
	} else if len(otherattrs) != 0 {
		// Case 4: other attributes (structure)
		if len(children) != 0 || strings.TrimSpace(chardata) != "" {
			return Value{}, fmt.Errorf("%s:%s has structure attributes and also children", elm.Space, elm.Tag)
		}
		var str Struct
		for _, attr := range otherattrs {
			str[Name{p.nsuris[attr.Space], attr.Key}] = Value{Value: attr.Value}
		}
		value.Value = str
	} else {
		switch len(children) {
		case 0:
			// Case 5: no attributes, textual data
			value.Value = chardata
		case 1:
			p.pushElementNamespaces(children[0])
			if p.nsuris[children[0].Space] != nsRDF {
				return Value{}, fmt.Errorf("%s:%s has unrecognized child %s:%s", elm.Space, elm.Tag, children[0].Space, children[0].Tag)
			}
			switch children[0].Tag {
			case "Description":
				// Case 6: no attributes, single rdf:Description child (structure or qualified value)
				var desc Value
				if desc, err = p.readValueDescElm(children[0]); err != nil {
					return Value{}, err
				}
				if _, ok := desc.Value.(Struct); !ok {
					// If it's not a struct, it must be a qualified value.  Merge the qualifiers.
					value.Qualifiers = append(value.Qualifiers, desc.Qualifiers...)
					value.Value = desc.Value
				} else {
					value = desc
				}
			case "Alt", "Bag", "Seq":
				// Case 7: no attributes, single rdf:Bag, rdf:Seq, or rdf:Alt child (array)
				if value.Value, err = p.readArray(children[0]); err != nil {
					return Value{}, err
				}
			default:
				return Value{}, fmt.Errorf("%s:%s has unrecognized child %s:%s", elm.Space, elm.Tag, children[0].Space, children[0].Tag)
			}
			p.popElementNamespaces(children[0])
		default:
			return Value{}, fmt.Errorf("%s:%s has multiple children", elm.Space, elm.Tag)
		}
	}
	// If the resulting value is a structure with an rdf:value field, then
	// it actually isn't a structure; it's a qualified value.
	if str, ok := value.Value.(Struct); ok {
		if val, ok := str[Name{nsRDF, "value"}]; ok {
			for key, qual := range str {
				if key.Namespace != nsRDF || key.Name != "value" {
					value.Qualifiers = append(value.Qualifiers, Qualifier{key, Value{Value: qual}})
				}
			}
			value.Value = val
		}
	}
	return value, nil
}

// readValueDescElm reads an rdf:Description element for a value (either a
// structure or a qualified value).
func (p *Packet) readValueDescElm(elm *etree.Element) (value Value, err error) {
	var (
		str        Struct
		otherattrs []etree.Attr
		children   []*etree.Element
		chardata   string
	)
	// Walk through the list of attributes and pull out the interesting ones.
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			continue
		}
		switch nsuri := p.nsuris[attr.Space]; nsuri {
		case nsRDF:
			switch attr.Key {
			case "value":
				if _, ok := str[Name{nsuri, attr.Key}]; ok {
					return Value{}, fmt.Errorf("multiple values for %s:%s", attr.Space, attr.Key)
				}
				str[Name{nsuri, attr.Key}] = Value{Value: attr.Value}
			default:
				return Value{}, fmt.Errorf("%s:%s has unsupported %s:%s attribute", elm.Space, elm.Tag, attr.Space, attr.Key)
			}
		case nsXML:
			return Value{}, fmt.Errorf("%s:%s has unsupported %s:%s attribute", elm.Space, elm.Tag, attr.Space, attr.Key)
		case "":
			return Value{}, fmt.Errorf("%s:%s has attribute %s with unregistered namespace", elm.Space, elm.Tag, attr.FullKey())
		default:
			if _, ok := str[Name{nsuri, attr.Key}]; ok {
				return Value{}, fmt.Errorf("multiple values for %s:%s", attr.Space, attr.Key)
			}
			str[Name{nsuri, attr.Key}] = Value{Value: attr.Value}
		}
	}
	if children, err = elmChildren(elm); err != nil {
		return Value{}, err
	}
	for _, child := range children {
		p.pushElementNamespaces(child)
		switch nsuri := p.nsuris[child.Space]; nsuri {
		case nsRDF, nsXML:
			return Value{}, fmt.Errorf("rdf:Description contains unexpected %s:%s element", child.Space, child.Tag)
		case "":
			return Value{}, fmt.Errorf("%s element has unregistered namespace", child.FullTag())
		default:
			if _, ok := str[Name{nsuri, child.Tag}]; ok {
				return Value{}, fmt.Errorf("multiple values for %s:%s", child.Space, child.Tag)
			}
			if str[Name{nsuri, child.Tag}], err = p.readValueElm(child, false); err != nil {
				return Value{}, err
			}
		}
		p.popElementNamespaces(child)
	}
	// If the resulting value is a structure with an rdf:value field, then
	// it actually isn't a structure; it's a qualified value.
	if val, ok := str[Name{nsRDF, "value"}]; ok {
		for key, qual := range str {
			if key.Namespace != nsRDF || key.Name != "value" {
				value.Qualifiers = append(value.Qualifiers, Qualifier{key, Value{Value: qual}})
			}
		}
		value.Value = val
	} else {
		value.Value = str
	}
	return value, nil
}

// getAttribute gets the value of an attribute of an element.
func (p *Packet) getAttribute(elm *etree.Element, nsuri, key string) (value string, found bool) {
	for _, attr := range elm.Attr {
		if p.nsuris[attr.Space] == nsuri && attr.Key == key {
			return attr.Value, true
		}
	}
	return "", false
}

// pushElementNamespaces registers the namespaces named in xmlns attributes on
// the element.  Any previous registrations of them are pushed on a stack, to be
// restored by popElementNamespaces.
func (p *Packet) pushElementNamespaces(elm *etree.Element) {
	var was = make(map[string]string)

	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			if old, ok := p.nsuris[attr.Key]; ok {
				was[attr.Key] = old
			}
			p.nsuris[attr.Key] = attr.Value
			p.nsprefixes[attr.Value] = attr.Key
		}
	}
	p.nsstack = append(p.nsstack, was)
}

// popElementNamespace unregisters the namespaces named in xmlns attributes on
// the element, restoring their previous registrations if any.
func (p *Packet) popElementNamespaces(elm *etree.Element) {
	var was map[string]string

	was, p.nsstack = p.nsstack[len(p.nsstack)-1], p.nsstack[:len(p.nsstack)-1]
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			if old, ok := was[attr.Key]; ok {
				p.nsuris[attr.Key] = old
			} else {
				delete(p.nsuris, attr.Key)
			}
		}
	}
}

// elmChildren ensures that an element has only element children, and returns
// them.
func elmChildren(elm *etree.Element) (children []*etree.Element, err error) {
	var chardata string
	if children, chardata, err = parseElement(elm); err != nil {
		return nil, err
	}
	if chardata != "" {
		return nil, fmt.Errorf("text found instead of child elements in %s:%s element", elm.Space, elm.Tag)
	}
	return children, nil
}

// parseElement ensures that an element has either element children or
// non-whitespace cdata children, but not both; it returns whichever it finds.
func parseElement(elm *etree.Element) (children []*etree.Element, chardata string, err error) {
	for _, child := range elm.Child {
		switch child := child.(type) {
		case *etree.Element:
			children = append(children, child)
		case *etree.CharData:
			chardata += child.Data
		default:
			// Comments, processing instructions, and directives are ignored.
		}
	}
	if len(children) != 0 {
		if strings.TrimSpace(chardata) != "" {
			return nil, "", fmt.Errorf("mixed text and element contents of %s:%s", elm.Space, elm.Tag)
		}
		return children, "", nil
	}
	return nil, chardata, nil
}
