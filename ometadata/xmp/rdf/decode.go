package rdf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

// Constants for namespace URIs.
const (
	NSrdf = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	NSx   = "adobe:ns:meta/"
	NSxml = "http://www.w3.org/XML/1998/namespace"
)

// NewPacket creates a new, empty RDF packet.
func NewPacket() *Packet {
	return &Packet{
		Properties: make(Struct),
		nsprefixes: map[string]string{NSxml: "xml"},
	}
}

// ReadPacket parses the supplied buffer and returns the resulting Packet, or an
// error if the packet is invalid.
func ReadPacket(buf []byte) (p *Packet, err error) {
	var (
		doc    *etree.Document
		root   *element
		nsuris = map[string]string{"xml": NSxml}
	)
	p = NewPacket()
	doc = etree.NewDocument()
	if err = doc.ReadFromBytes(buf); err != nil {
		return nil, err
	}
	if root, err = simplifyElement(doc.Root(), nsuris, p.nsprefixes); err != nil {
		return nil, err
	}
	if root.name.is(NSx, "xmpmeta") {
		if len(root.children) != 1 {
			return nil, fmt.Errorf("%s: expected one rdf:RDF child element", root.path())
		}
		root = root.children[0]
	}
	if !root.name.is(NSrdf, "RDF") {
		return nil, errors.New("XMP root element must be rdf:RDF")
	}
	if err = p.readRDF(root); err != nil {
		return nil, err
	}
	return p, nil
}

// readRDF reads the root RDF element.
func (p *Packet) readRDF(elm *element) (err error) {
	if elm.text != "" {
		return fmt.Errorf("%s: unexpected text content", elm.path())
	}
	if len(elm.attrs) != 0 {
		return fmt.Errorf("%s: cannot have attributes", elm.path())
	}
	for _, child := range elm.children {
		if !child.name.is(NSrdf, "Description") {
			return fmt.Errorf("%s: unexpected child element %s", elm.path(), child.name)
		}
		if err = p.readPropertyDescription(child); err != nil {
			return err
		}
	}
	return nil
}

// readPropertyDescription handles a top-level rdf:Description element.
func (p *Packet) readPropertyDescription(elm *element) (err error) {
	if elm.text != "" {
		return fmt.Errorf("%s: unexpected text content", elm.path())
	}
	for key, val := range elm.attrs {
		if key.is(NSrdf, "about") {
			if p.about != "" && val != "" && p.about != val {
				return errors.New("rdf:about: inconsistent values")
			}
			if val != "" {
				p.about = val
			}
			continue
		}
		if key.Namespace == NSrdf || key.Namespace == NSxml {
			return fmt.Errorf("%s: unexpected attribute %s", elm.path(), key)
		}
		if _, ok := p.Properties[key]; ok {
			return fmt.Errorf("multiple values for %s", key)
		}
		p.Properties[key] = Value{Value: val}
	}
	for _, child := range elm.children {
		if err = p.readProperty(child); err != nil {
			return err
		}
	}
	return nil
}

// readProperty handles a property definition element, as an immediate child of
// a top-level rdf:Description element.
func (p *Packet) readProperty(elm *element) (err error) {
	if elm.name.Namespace == NSrdf || elm.name.Namespace == NSxml {
		return fmt.Errorf("%s: unexpected element", elm.path())
	}
	if _, ok := p.Properties[elm.name]; ok {
		return fmt.Errorf("multiple values for %s", elm.path())
	}
	var value Value
	if err = p.readValueElm(elm, &value); err != nil {
		return err
	}
	p.Properties[elm.name] = value
	return nil
}

// readNQValueElm reads a value from an rdf:value element.  It's the same as
// readValueElm except that it disallows a qualified value.
func (p *Packet) readNQValueElm(elm *element, value *Value) (err error) {
	err = p.readValueElm(elm, value)
	if err == nil && len(value.Qualifiers) != 0 {
		*value = Value{}
		err = fmt.Errorf("%s: nested qualifiers not allowed", elm.path())
	}
	return
}

// readValueElm reads a value from an element.  The element could be:
//   - A property element under a top-level rdf:Description
//   - A structure field element under a structure's rdf:Description
//   - An array element (rdf:li)
//   - A qualifier element, under a qualified value's rdf:Description
//   - A qualified value (rdf:value)
func (p *Packet) readValueElm(elm *element, value *Value) (err error) {
	if lang, ok := elm.attrs[Name{NSxml, "lang"}]; ok {
		delete(elm.attrs, Name{NSxml, "lang"})
		if value.Qualifiers == nil {
			value.Qualifiers = make(map[Name]Value)
		}
		value.Qualifiers[Name{NSxml, "lang"}] = Value{Value: lang}
	}
	err = p.readTextValueElm(elm, value)
	if err == nil && value.Value == nil {
		err = p.readURIValueElm(elm, value)
	}
	if err == nil && value.Value == nil {
		err = p.readArrayValueElm(elm, value)
	}
	if err == nil && value.Value == nil {
		err = p.readDescribedValueElm(elm, value)
	}
	if err == nil && value.Value == nil {
		err = p.readPTResourceValueElm(elm, value)
	}
	if err == nil && value.Value == nil {
		err = p.readQualifiedValueElm(elm, value)
	}
	if err == nil && value.Value == nil {
		err = p.readStructureValueElm(elm, value)
	}
	if err == nil && value.Value == nil {
		err = fmt.Errorf("%s: invalid value element", elm.path())
	}
	return err
}

// readTextValueElm examines elm to see if it is a text value.
func (p *Packet) readTextValueElm(elm *element, value *Value) (err error) {
	if elm.text == "" && (len(elm.attrs) != 0 || len(elm.children) != 0) {
		return nil
	}
	if len(elm.attrs) != 0 {
		return fmt.Errorf("%s: text value element cannot have attributes", elm.path())
	}
	value.Value = elm.text
	return nil
}

// readURIValueElm examines elm to see if it is a URI value.
func (p *Packet) readURIValueElm(elm *element, value *Value) (err error) {
	uri, ok := elm.attrs[Name{NSrdf, "resource"}]
	if !ok {
		return nil
	}
	// Don't need to check elm.text; that was done by readTextValueElm.
	if len(elm.attrs) != 1 {
		return fmt.Errorf("%s: element with rdf:resource attribute cannot have other attributes", elm.path())
	}
	if len(elm.children) != 0 {
		return fmt.Errorf("%s: element with rdf:resource attribute cannot have content", elm.path())
	}
	value.Value = URI(uri)
	return nil
}

// readArrayValueElm examines elm to see if it is an array value
func (p *Packet) readArrayValueElm(elm *element, value *Value) (err error) {
	var (
		aryelm *element
		ary    []Value
	)
	if len(elm.children) != 1 {
		return nil
	}
	aryelm = elm.children[0]
	if !aryelm.name.is(NSrdf, "Alt") && !aryelm.name.is(NSrdf, "Bag") && !aryelm.name.is(NSrdf, "Seq") {
		return nil
	}
	// Don't need to check elm.text; that was done by readTextValueElm.
	if len(elm.attrs) != 0 {
		return fmt.Errorf("%s: element with array child cannot have attributes", elm.path())
	}
	if ary, err = p.readArrayListElm(aryelm); err != nil {
		return err
	}
	switch aryelm.name.Name {
	case "Alt":
		value.Value = Alt(ary)
	case "Bag":
		value.Value = Bag(ary)
	case "Seq":
		value.Value = Seq(ary)
	}
	return nil
}

// readArrayListElm reads an rdf:Alt, rdf:Bag, or rdf:Seq element.
func (p *Packet) readArrayListElm(elm *element) (ary []Value, err error) {
	if elm.text != "" {
		return nil, fmt.Errorf("%s: unexpected text content", elm.path())
	}
	if len(elm.attrs) != 0 {
		return nil, fmt.Errorf("%s: element cannot have attributes", elm.path())
	}
	for _, lielm := range elm.children {
		if !lielm.name.is(NSrdf, "li") {
			return nil, fmt.Errorf("%s: unexpected child element %s", elm.path(), lielm.name)
		}
		var val Value
		if err = p.readValueElm(lielm, &val); err != nil {
			return nil, err
		}
		ary = append(ary, val)
	}
	return ary, nil
}

// readDescribedValueElm examines elm to see if it contains a value description
// (i.e., a single rdf:Description child element).
func (p *Packet) readDescribedValueElm(elm *element, value *Value) (err error) {
	var desc *element

	if len(elm.children) != 1 || !elm.children[0].name.is(NSrdf, "Description") {
		return nil
	}
	desc = elm.children[0]
	// Don't need to check elm.text; that was done by readTextValueElm.
	if len(elm.attrs) != 0 {
		return fmt.Errorf("%s: element cannot have attributes", elm.path())
	}
	err = p.readQualifiedValueElm(desc, value)
	if err == nil && value.Value == nil {
		err = p.readStructureValueElm(desc, value)
	}
	return err
}

// readPTResourceValueElm examines elm to see if it contains an rdf:parseType
// Resource.
func (p *Packet) readPTResourceValueElm(elm *element, value *Value) (err error) {
	if val := elm.attrs[Name{NSrdf, "parseType"}]; val != "Resource" {
		return nil
	}
	delete(elm.attrs, Name{NSrdf, "parseType"})
	if len(elm.attrs) != 0 {
		return fmt.Errorf("%s: element with rdf:parseType attribute cannot have other attributes", elm.path())
	}
	err = p.readQualifiedValueElm(elm, value)
	if err == nil && value.Value == nil {
		err = p.readStructureValueElm(elm, value)
	}
	return err
}

// readQualifiedValueElm examines elm to see if it is a qualified value (i.e.,
// an rdf:value attribute or an rdf:value child element).
func (p *Packet) readQualifiedValueElm(elm *element, value *Value) (err error) {
	if val, ok := elm.attrs[Name{NSrdf, "value"}]; ok {
		value.Value = val
	} else {
		for _, child := range elm.children {
			if child.name.is(NSrdf, "value") {
				var val Value
				if err = p.readNQValueElm(child, &val); err != nil {
					return err
				}
				value.Value = val.Value
				break
			}
		}
	}
	if value.Value == nil {
		return nil
	}
	if value.Qualifiers == nil {
		value.Qualifiers = make(map[Name]Value)
	}
	for key, val := range elm.attrs {
		if key.is(NSrdf, "value") {
			continue
		}
		if key.Namespace == NSrdf || key.Namespace == NSxml {
			return fmt.Errorf("%s: unexpected attribute %s", elm.path(), key)
		}
		value.Qualifiers[key] = Value{Value: val}
	}
	for _, child := range elm.children {
		if child.name.is(NSrdf, "value") {
			continue
		}
		if child.name.Namespace == NSrdf || child.name.Namespace == NSxml {
			return fmt.Errorf("%s: unexpected child element %s", elm.path(), child.name)
		}
		var val Value
		if err = p.readValueElm(child, &val); err != nil {
			return err
		}
		value.Qualifiers[child.name] = val
	}
	return nil
}

// readStructureValueElm reads elm as a structure value.  It always fills in
// value, unless it returns an error.
func (p *Packet) readStructureValueElm(elm *element, value *Value) (err error) {
	var str = make(Struct)

	for key, val := range elm.attrs {
		if key.Namespace == NSrdf || key.Namespace == NSxml {
			return fmt.Errorf("%s: unexpected attribute %s", elm.path(), key)
		}
		str[key] = Value{Value: val}
	}
	for _, child := range elm.children {
		if child.name.Namespace == NSrdf || child.name.Namespace == NSxml {
			return fmt.Errorf("%s: unexpected child element %s", elm.path(), child.name)
		}
		var val Value
		if err = p.readValueElm(child, &val); err != nil {
			return err
		}
		str[child.name] = val
	}
	value.Value = str
	return nil
}

// An element is a simplified version of an etree.Element.  The document tree is
// converted into this element form in order to make the parsing algorithms
// simpler.
type element struct {
	name     Name
	attrs    map[Name]string
	children []*element
	text     string
	parent   *element
}

// simplifyElement recursively converts an XML/RDF etree.Element into a
// package-local element structure.  In the process, it updates the nsprefixes
// map, which maps from namespace URI to the namespace prefix used for it in the
// document.  nsuris is a reverse map which can be seeded with namespace URIs
// for predefined prefixes (generally just "xml").
func simplifyElement(rdfe *etree.Element, nsuris, nsprefixes map[string]string) (se *element, err error) {
	var (
		nsuri      string
		saveNSuris = make(map[string]string)
	)
	for _, attr := range rdfe.Attr {
		if attr.Space != "xmlns" {
			continue
		}
		if old, ok := nsuris[attr.Key]; ok {
			saveNSuris[attr.Key] = old
		}
		nsuris[attr.Key] = attr.Value
		nsprefixes[attr.Value] = attr.Key
	}
	if nsuri = nsuris[rdfe.Space]; nsuri == "" {
		return nil, fmt.Errorf("%s: unregistered namespace", rdfe.FullTag())
	}
	se = &element{
		name:  Name{nsuri, rdfe.Tag},
		attrs: make(map[Name]string),
	}
	for _, child := range rdfe.Child {
		switch child := child.(type) {
		case *etree.Element:
			if ce, err := simplifyElement(child, nsuris, nsprefixes); err == nil {
				ce.parent = se
				se.children = append(se.children, ce)
			} else {
				return nil, err
			}
		case *etree.CharData:
			se.text += child.Data
		}
	}
	if len(se.children) != 0 {
		if strings.TrimSpace(se.text) != "" {
			return nil, fmt.Errorf("%s: cannot have both child elements and text content", rdfe.FullTag())
		}
		se.text = ""
	}
	for _, attr := range rdfe.Attr {
		if attr.Space == "xmlns" {
			continue
		}
		if nsuri = nsuris[attr.Space]; nsuri == "" {
			return nil, fmt.Errorf("%s: %s: unregistered namespace", rdfe.FullTag(), attr.FullKey())
		}
		se.attrs[Name{nsuri, attr.Key}] = attr.Value
	}
	for _, attr := range rdfe.Attr {
		if attr.Space != "xmlns" {
			continue
		}
		delete(nsuris, attr.Key)
	}
	for key, val := range saveNSuris {
		nsuris[key] = val
	}
	return se, nil
}

func (e *element) path() string {
	var list []string
	for ; e != nil; e = e.parent {
		list = append(list, e.name.Name)
	}
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return strings.Join(list, "/")
}
