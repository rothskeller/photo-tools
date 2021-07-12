// Package rdf handles the marshaling and unmarshaling of RDF documents, obeying
// (and limited to) the syntax described in the Adobe XMP Specification, Part 1.
package rdf

import "fmt"

// A Packet represents the entire RDF packet.
type Packet struct {
	properties Struct
	nsprefixes map[string]string
	about      string
	dirty      bool
}

// A Value represents a value in an RDF file.  It comprises zero or more
// qualifiers, plus a simple value.
type Value struct {
	Qualifiers map[Name]Value
	Value      SimpleValue
}

// A Name is the name for a property, a structure field, or a qualifier.  It has
// a namespace URI and a local name.  Note that the prefix used to represent the
// namespace is not formally part of the name.
type Name struct {
	Namespace string
	Name      string
}

// A SimpleValue is an unqualified value of a property, structure field, or
// qualifier.  The semantically allowed types are string, URI, Seq, Bag, Alt,
// and Struct.
type SimpleValue interface{}

// A URI is a string containing a URI.  This is semantically equivalent to a
// regular string, but encoded differently.
type URI string

// A Seq is an ordered list of values.
type Seq []Value

// A Bag is an unordered list of values.
type Bag []Value

// An Alt is an ordered set of alternative values, with the first one being
// considered the default.
type Alt []Value

// A Struct is an unordered set of name/value pairs.
type Struct map[Name]Value

// RegisterNamespace sets the namespace prefix to use for the specified URI.
func (p *Packet) RegisterNamespace(prefix, uri string) {
	p.nsprefixes[uri] = prefix
}

// Properties returns a list of all defined property names (in unspecified
// order).
func (p *Packet) Properties() (names []Name) {
	names = make([]Name, 0, len(p.properties))
	for name := range p.properties {
		names = append(names, name)
	}
	return names
}

// Property returns the Value of the property with the specified Name, or nil if
// that property does not exist.
func (p *Packet) Property(name Name) Value { return p.properties[name] }

// SetProperty sets the Value of the property with the specified Name to the
// specified Value.  It also marks the RDF packet as being dirty.  Callers
// should ensure that they only call SetProperty when the value of the property
// has actually changed.
func (p *Packet) SetProperty(name Name, value Value) {
	p.properties[name] = value
	p.dirty = true
}

// RemoveProperty removes the property with the specified Name.  It marks the
// RDF packet as being dirty if the property previously existed.
func (p *Packet) RemoveProperty(name Name) {
	if _, ok := p.properties[name]; ok {
		delete(p.properties, name)
		p.dirty = true
	}
}

// Dirty returns whether the RDF packet has been changed since it was read.
func (p *Packet) Dirty() bool { return p.dirty }

func (n Name) String() string {
	return fmt.Sprintf("[%s]%s", n.Namespace, n.Name)
}

// is tests a name for equality.
func (n Name) is(space, local string) bool {
	return n.Namespace == space && n.Name == local
}
