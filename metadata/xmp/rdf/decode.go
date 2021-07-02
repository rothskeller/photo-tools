package rdf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

const (
	NSrdf = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	NSx   = "adobe:ns:meta/"
	NSxml = "http://www.w3.org/XML/1998/namespace"
)

// NewPacket creates a new, empty RDF packet.
func NewPacket() *Packet {
	return &Packet{
		properties: make(Struct),
		nsprefixes: map[string]string{NSxml: "xml"},
		nsuris:     map[string]string{"xml": NSxml},
	}
}

// ReadPacket parses the supplied buffer and returns the resulting Packet, or an
// error if the packet is invalid.
func ReadPacket(buf []byte) (p *Packet, err error) {
	var (
		doc *etree.Document
		// root *etree.Element
	)
	p = NewPacket()
	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(buf); err != nil {
		return nil, err
	}
	if err = simplifyDoc(&doc.Element); err != nil {
		return nil, err
	}
	/*
		if root, err = expectOneElement(&doc.Element); err != nil {
			return nil, err
		}
		p.pushElementNamespaces(root)
		if p.nsuris[root.Space] == NSx && root.Tag == "xmpmeta" {
			if root, err = expectOneElement(root); err != nil {
				return nil, err
			}
			p.pushElementNamespaces(root)
		}
		if p.nsuris[root.Space] != NSrdf || root.Tag != "RDF" {
			return nil, errors.New("XMP root element must be rdf:RDF")
		}
		if err = p.readRDF(root); err != nil {
			return nil, err
		}
		return p, nil
	*/
	return nil, nil
}

// readRDF reads the root RDF element.
func (p *Packet) readRDF(elm *etree.Element) (err error) {
	for _, attr := range elm.Attr {
		if attr.Space != "xmlns" {
			return fmt.Errorf("%s: unexpected attribute %s", elm.FullTag(), attr.FullKey())
		}
	}
	var children []*etree.Element
	if children, err = expectElements(elm); err != nil {
		return err
	}
	for _, child := range children {
		p.pushElementNamespaces(child)
		if p.nsuris[child.Space] != NSrdf || child.Tag != "Description" {
			return fmt.Errorf("%s: unexpected child element %s", elm.FullTag(), child.FullTag())
		}
		if err = p.readPropertyDescription(elm); err != nil {
			return err
		}
		p.popElementNamespaces(child)
	}
	return nil
}

// readPropertyDescription handles a top-level rdf:Description element.
func (p *Packet) readPropertyDescription(elm *etree.Element) (err error) {
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			continue
		}
		var nsuri = p.nsuris[attr.Space]
		if nsuri == "" {
			return fmt.Errorf("%s: %s: unregistered namespace", elm.FullTag(), attr.FullKey())
		}
		if nsuri == NSrdf && attr.Key == "about" {
			if p.about != "" && attr.Value != "" && p.about != attr.Value {
				return errors.New("mismatched values for rdf:about")
			}
			if attr.Value != "" {
				p.about = attr.Value
			}
			continue
		}
		if nsuri == NSrdf || nsuri == NSxml {
			return fmt.Errorf("%s: unexpected attribute %s", elm.FullTag(), attr.FullKey())
		}
		var name = Name{nsuri, attr.Key}
		if _, ok := p.properties[name]; ok {
			return fmt.Errorf("multiple values for %s", attr.FullKey())
		}
		p.properties[name] = Value{Value: attr.Value}
	}
	var children []*etree.Element
	if children, err = expectElements(elm); err != nil {
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
		return fmt.Errorf("%s: unregistered namespace", elm.FullTag())
	}
	if nsuri == NSrdf || nsuri == NSxml {
		return fmt.Errorf("%s: unexpected element", elm.FullTag())
	}
	var name = Name{nsuri, elm.Tag}
	if _, ok := p.properties[name]; ok {
		return fmt.Errorf("multiple values for %s", elm.FullTag())
	}
	var value Value
	if value, err = p.readValueElm(elm, true); err != nil {
		return err
	}
	p.properties[name] = value
	return nil
}

// readValueElm reads a value from an element.  The element could be:
//   - A property element under a top-level rdf:Description
//   - A structure field element under a structure's rdf:Description
//   - An array element (rdf:li)
//   - A qualifier element, under a qualified value's rdf:Description
//   - A qualified value (rdf:value)
// The qualOK parameter will be true for the first four of those, and false for
// the last one.  When false, the passed element cannot have an rdf:value or
// xml:lang attribute, an rdf:ParseType="resource" attribute with an rdf:value
// child element, or an rdf:Description child element with an rdf:value
// attribute or child element.  The function returns a Value; if qualOK is
// false, the Value is guaranteed to have no qualifiers.
func (p *Packet) readValueElm(elm *etree.Element, qualOK bool) (value Value, err error) {
	var (
		attrs    []etree.Attr
		children []*etree.Element
		chardata string
		str      Struct
	)
	for _, attr := range elm.Attr {
		switch {
		case attr.Space == "xmlns":
			break
		case p.nsuris[attr.Space] == NSxml && attr.Key == "lang":
			if value.Qualifiers != nil {
				return Value{}, fmt.Errorf("%s: multiple values for %s", elm.FullTag(), attr.FullKey())
			}
			value.Qualifiers = make(map[Name]Value)
			value.Qualifiers[Name{NSxml, "lang"}] = Value{Value: attr.Value}
		default:
			attrs = append(attrs, attr)
		}
	}
	children, chardata = parseElement(elm)
	// Case 1: no attributes and no children
	if len(children) == 0 && len(attrs) == 0 {
		value.Value = chardata
		goto CHECKQUAL
	}
	if strings.TrimSpace(chardata) != "" {
		return Value{}, fmt.Errorf("%s: unexpected text data", elm.FullTag())
	}
	// Case 2: single rdf:resource attribute and no children
	if len(children) == 0 && len(attrs) == 1 && attrs[0].Key == "resource" && p.nsuris[attrs[0].Space] == NSrdf {
		value.Value = URI(attrs[0].Value)
		goto CHECKQUAL
	}
	// Case 3: single rdf:ParseType="Resource" attribute
	if len(attrs) == 1 && attrs[0].Key == "ParseType" && attrs[0].Value == "Resource" && p.nsuris[attrs[0].Space] == NSrdf {
		str = make(Struct)
		for _, child := range children {
			p.pushElementNamespaces(child)
			switch nsuri := p.nsuris[child.Space]; nsuri {
			case NSrdf, NSxml:
				return Value{}, fmt.Errorf("%s: unexpected child element %s", elm.FullTag(), child.FullTag())
			case "":
				return Value{}, fmt.Errorf("%s: unregistered namespace", child.FullTag())
			default:
				if _, ok := str[Name{nsuri, child.Tag}]; ok {
					return Value{}, fmt.Errorf("%s: multiple values for %s", elm.FullTag(), child.FullTag())
				}
				var cval Value
				if cval, err = p.readValueElm(child, true); err != nil {
					return Value{}, err
				}
				str[Name{nsuri, child.Tag}] = cval
			}
			p.popElementNamespaces(child)
		}
		value.Value = str
		goto CHECKQUAL
	}
	if len(attrs) == 0 && len(children) == 1 {
		var child = children[0]
		p.pushElementNamespaces(child)
		// Case 4: single rdf:Description child
		if child.Tag == "Description" && p.nsuris[child.Space] == NSrdf {
			if value.Value, err = p.readStructDescription(child); err != nil {
				return Value{}, err
			}
			p.popElementNamespaces(child)
			goto CHECKQUAL
		}
		// Case 5: single rdf:Alt, rdf:Bag, or rdf:Seq child
		if (child.Tag == "Alt" || child.Tag == "Bag" || child.Tag == "Seq") && p.nsuris[child.Space] == NSrdf {
			var vals []Value
			if vals, err = p.readArray(child); err != nil {
				return Value{}, err
			}
			switch child.Tag {
			case "Alt":
				value.Value = Alt(vals)
			case "Bag":
				value.Value = Bag(vals)
			case "Seq":
				value.Value = Seq(vals)
			}
			p.popElementNamespaces(child)
			goto CHECKQUAL
		}
		p.popElementNamespaces(child)
	}
	// Case 6: structure type
	str = make(Struct)
	for _, attr := range attrs {
		var nsuri = p.nsuris[attr.Space]
		if nsuri == "" {
			return Value{}, fmt.Errorf("%s: %s: unregistered namespace", elm.FullTag(), attr.FullKey())
		}
		if nsuri == NSrdf || nsuri == NSxml {
			return Value{}, fmt.Errorf("%s: unrecognized attribute %s", elm.FullTag(), attr.FullKey())
		}
		if _, ok := str[Name{nsuri, attr.Key}]; ok {
			return Value{}, fmt.Errorf("%s: multiple values for %s", elm.FullTag(), attr.FullKey())
		}
		str[Name{nsuri, attr.Key}] = Value{Value: attr.Value}
	}
	for _, child := range children {
		p.pushElementNamespaces(child)
		var nsuri = p.nsuris[child.Space]
		if nsuri == "" {
			return Value{}, fmt.Errorf("%s: %s: unregistered namespace", elm.FullTag(), child.FullTag())
		}
		if nsuri == NSrdf {
			return Value{}, fmt.Errorf("%s: unrecognized child element %s", elm.FullTag(), child.FullTag())
		}
		if _, ok := str[Name{nsuri, child.Tag}]; ok {
			return Value{}, fmt.Errorf("%s: multiple values for %s", elm.FullTag(), child.FullTag())
		}
		var qualOK = nsuri != NSrdf || child.Tag != "value"
		if str[Name{nsuri, child.Tag}], err = p.readValueElm(child, qualOK); err != nil {
			return Value{}, err
		}
		p.popElementNamespaces(child)
	}
CHECKQUAL:
	if str, ok := value.Value.(Struct); ok {
		if qval, ok := str[Name{NSrdf, "value"}]; ok {
			// This is not a structure, it's a qualified value.
			if value.Qualifiers == nil {
				value.Qualifiers = make(map[Name]Value)
			}
			for key, qual := range str {
				if key.Namespace != NSrdf || key.Name != "value" {
					if _, ok := value.Qualifiers[key]; ok {
						return Value{}, fmt.Errorf("%s: multiple values for %s:%s", elm.FullTag(), p.nsprefixes[key.Namespace], key.Name)
					}
					value.Qualifiers[key] = qual
				}
			}
			value.Value = qval
		}
	}
	if len(value.Qualifiers) != 0 && !qualOK {
		return Value{}, fmt.Errorf("%s: qualifiers not allowed here", elm.FullTag())
	}
	return value, nil
}

// readStructDescription reads an rdf:Description element for a structure or a
// qualified value.  The value is returned as a structure in either case.
func (p *Packet) readStructDescription(elm *etree.Element) (str Struct, err error) {
	str = make(Struct)
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			continue
		}
		var nsuri = p.nsuris[attr.Space]
		if nsuri == "" {
			return nil, fmt.Errorf("%s: %s: unregistered namespace", elm.FullTag(), attr.FullKey())
		}
		if _, ok := str[Name{nsuri, attr.Key}]; ok {
			return nil, fmt.Errorf("%s: multiple values for %s", elm.FullTag(), attr.FullKey())
		}
		str[Name{nsuri, attr.Key}] = Value{Value: attr.Value}
	}
	var children []*etree.Element
	if children, err = expectElements(elm); err != nil {
		return nil, err
	}
	for _, child := range children {
		p.pushElementNamespaces(child)
		var nsuri = p.nsuris[child.Space]
		if nsuri == "" {
			return nil, fmt.Errorf("%s: %s: unregistered namespace", elm.FullTag(), child.FullTag())
		}
		if _, ok := str[Name{nsuri, child.Tag}]; ok {
			return nil, fmt.Errorf("%s: multiple values for %s", elm.FullTag(), child.FullTag())
		}
		var qualOK = nsuri != NSrdf || child.Tag != "value"
		if str[Name{nsuri, child.Tag}], err = p.readValueElm(child, qualOK); err != nil {
			return nil, err
		}
		p.popElementNamespaces(child)
	}
	return str, nil
}

// readArray reads an rdf:Alt, rdf:Bag, or rdf:Seq element.
func (p *Packet) readArray(elm *etree.Element) (ary []Value, err error) {
	for _, attr := range elm.Attr {
		if attr.Space == "xmlns" {
			continue
		}
		return nil, fmt.Errorf("%s: unexpected attribute %s", elm.FullTag(), attr.FullKey())
	}
	var children []*etree.Element
	if children, err = expectElements(elm); err != nil {
		return nil, err
	}
	for _, child := range children {
		p.pushElementNamespaces(child)
		if child.Tag != "li" || p.nsuris[child.Space] != NSrdf {
			return nil, fmt.Errorf("%s: unexpected child element %s", elm.FullTag(), child.FullTag())
		}
		var val Value
		if val, err = p.readValueElm(child, true); err != nil {
			return nil, err
		}
		ary = append(ary, val)
		p.popElementNamespaces(child)
	}
	return ary, nil
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

// simplifyDoc removes nodes from the document that will distract us: comments,
// processing instructions, and directives.  It also removes whitespace-only
// text nodes that are siblings of element nodes.  It raises an error if there
// are any non-whitespace text nodes that are siblings of element nodes.  It
// merges multiple text nodes (not siblings of element nodes) into a single one.
func simplifyDoc(elm *etree.Element) error {
	var (
		seenElement  bool
		seenCharData *etree.CharData
	)
	for i := 0; i < len(elm.Child); {
		switch child := elm.Child[i].(type) {
		case *etree.Element:
			seenElement = true
			if err := simplifyDoc(child); err != nil {
				return err
			}
			i++
		case *etree.CharData:
			if seenCharData != nil {
				seenCharData.Data += child.Data
				elm.RemoveChild(child)
			} else {
				seenCharData = child
				i++
			}
		default:
			elm.RemoveChild(child)
		}
	}
	if seenElement && seenCharData != nil {
		if strings.TrimSpace(seenCharData.Data) != "" {
			return fmt.Errorf("invalid XML/RDF document: %s element has both element and text content", elm.FullTag())
		}
		elm.RemoveChild(seenCharData)
	}
	return nil
}

// expectElements returns the list of child elements from the supplied parent;
// it raises an error if the parent has text content.
func expectElements(parent *etree.Element) (children []*etree.Element, err error) {
	children = make([]*etree.Element, 0, len(parent.Child))
	for _, child := range parent.Child {
		switch child := child.(type) {
		case *etree.Element:
			children = append(children, child)
		default:
			return nil, fmt.Errorf("%s:%s cannot have text content", parent.Space, parent.Tag)
		}
	}
	return children, nil
}

// expectOneElement returns the child element from the supplied parent; it
// raises an error if the parent has text content, no child, or multiple
// children.
func expectOneElement(parent *etree.Element) (child *etree.Element, err error) {
	switch len(parent.Child) {
	case 0:
		return nil, fmt.Errorf("%s:%s is missing its required child element", parent.Space, parent.Tag)
	case 1:
		switch child := parent.Child[0].(type) {
		case *etree.Element:
			return child, nil
		default:
			return nil, fmt.Errorf("%s:%s cannot have text content", parent.Space, parent.Tag)
		}
	default:
		return nil, fmt.Errorf("%s:%s cannot have multiple child elements", parent.Space, parent.Tag)
	}
}

// parseElement returns the children or the character data from an element.
// (simplifyDoc ensured that it will be one or the other.)
func parseElement(elm *etree.Element) (children []*etree.Element, chardata string) {
	for _, child := range elm.Child {
		switch child := child.(type) {
		case *etree.Element:
			children = append(children, child)
		case *etree.CharData:
			chardata += child.Data
		}
	}
	return children, chardata
}
