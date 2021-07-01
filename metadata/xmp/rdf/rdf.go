// Package rdf handles the marshaling and unmarshaling of RDF documents, obeying
// (and limited to) the syntax described in the Adobe XMP Specification, Part 1.
package rdf

// A Packet represents the entire RDF packet.
type Packet struct {
	properties Struct
	nsprefixes map[string]string
	nsuris     map[string]string
	nsstack    []map[string]string
	about      string
}

// A Value represents a value in an RDF file.  It comprises zero or more
// qualifiers, plus a simple value.
type Value struct {
	Qualifiers []Qualifier
	Value      SimpleValue
}

// A Qualifier represents a qualifier in an RDF file.  It has a name and a
// value.
type Qualifier struct {
	Name  Name
	Value Value
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
