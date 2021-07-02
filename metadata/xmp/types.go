package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/rdf"
)

var xmlLang = rdf.Name{Namespace: rdf.NSxml, Name: "lang"}

// getAlt returns the value of a Language Alternative value from the XMP.
func (p *XMP) getAlt(from rdf.Struct, prefix, ns, name string) (as metadata.AltString) {
	if val, ok := from[rdf.Name{Namespace: ns, Name: name}]; ok {
		switch val := val.Value.(type) {
		case rdf.Alt:
			as = make(metadata.AltString, 0, len(val))
			for _, str := range val {
				var lang string
				if lt, ok := str.Qualifiers[xmlLang]; ok {
					switch lt := lt.Value.(type) {
					case string:
						lang = lt
					}
				}
				switch str := str.Value.(type) {
				case string:
					as = append(as, metadata.AltItem{Value: str, Lang: lang})
				default:
					p.log("%s:%s has wrong data type", ns, name)
				}
			}
		default:
			p.log("%s:%s has wrong data type", ns, name)
		}
	}
	return as
}

// setAlt sets a language alternative value in the XMP.
func (p *XMP) setAlt(in rdf.Struct, ns, name string, as metadata.AltString) {
	if len(as) == 0 {
		delete(in, rdf.Name{Namespace: ns, Name: name})
	} else {
		var values = make([]rdf.Value, len(as))
		for i := range as {
			values[i] = rdf.Value{
				Qualifiers: rdf.Struct{{Namespace: rdf.NSxml, Name: "lang"}: {Value: as[i].Lang}},
				Value:      rdf.Value{Value: as[i].Value},
			}
		}
		in[rdf.Name{Namespace: ns, Name: name}] = rdf.Value{Value: rdf.Alt(values)}
	}
	p.dirty = true
}

// getBag returns the value of an unordered array of text values from the XMP.
func (p *XMP) getBag(from rdf.Struct, prefix, ns, name string) (bag []string) {
	if val, ok := from[rdf.Name{Namespace: ns, Name: name}]; ok {
		switch val := val.Value.(type) {
		case rdf.Bag:
			bag = make([]string, 0, len(val))
			for _, str := range val {
				switch str := str.Value.(type) {
				case string:
					bag = append(bag, str)
				default:
					p.log("%s:%s has wrong data type", ns, name)
				}
			}
		default:
			p.log("%s:%s has wrong data type", ns, name)
		}
	}
	return bag
}

// setBag sets an unordered array of text values in the XMP.
func (p *XMP) setBag(in rdf.Struct, ns, name string, seq []string) {
	if len(seq) == 0 {
		delete(in, rdf.Name{Namespace: ns, Name: name})
	} else {
		var values = make([]rdf.Value, len(seq))
		for i := range seq {
			values[i] = rdf.Value{Value: seq[i]}
		}
		in[rdf.Name{Namespace: ns, Name: name}] = rdf.Value{Value: rdf.Bag(values)}
	}
	p.dirty = true
}

// getSeq returns the value of an ordered array of text values from the XMP.
func (p *XMP) getSeq(from rdf.Struct, prefix, ns, name string) (seq []string) {
	if val, ok := from[rdf.Name{Namespace: ns, Name: name}]; ok {
		switch val := val.Value.(type) {
		case rdf.Seq:
			seq = make([]string, 0, len(val))
			for _, str := range val {
				switch str := str.Value.(type) {
				case string:
					seq = append(seq, str)
				default:
					p.log("%s:%s has wrong data type", ns, name)
				}
			}
		default:
			p.log("%s:%s has wrong data type", ns, name)
		}
	}
	return seq
}

// setSeq sets an ordered array of text values in the XMP.
func (p *XMP) setSeq(in rdf.Struct, ns, name string, seq []string) {
	if len(seq) == 0 {
		delete(in, rdf.Name{Namespace: ns, Name: name})
	} else {
		var values = make([]rdf.Value, len(seq))
		for i := range seq {
			values[i] = rdf.Value{Value: seq[i]}
		}
		in[rdf.Name{Namespace: ns, Name: name}] = rdf.Value{Value: rdf.Seq(values)}
	}
	p.dirty = true
}

// getString returns the value of a simple string from the XMP.
func (p *XMP) getString(from rdf.Struct, prefix, ns, name string) (str string) {
	if val, ok := from[rdf.Name{Namespace: ns, Name: name}]; ok {
		switch val := val.Value.(type) {
		case string:
			str = val
		default:
			p.log("%s:%s has wrong data type", ns, name)
		}
	}
	return str
}

// setString sets a string value in the XMP.
func (p *XMP) setString(in rdf.Struct, ns, name string, str string) {
	if str == "" {
		delete(in, rdf.Name{Namespace: ns, Name: name})
	} else {
		in[rdf.Name{Namespace: ns, Name: name}] = rdf.Value{Value: str}
	}
	p.dirty = true
}

// getStrings returns the value of an array of text values from the XMP.  It
// accepts either Bag or Seq, or a single string value.
func (p *XMP) getStrings(from rdf.Struct, prefix, ns, name string) (list []string) {
	var vals []rdf.Value
	if val, ok := from[rdf.Name{Namespace: ns, Name: name}]; ok {
		switch val := val.Value.(type) {
		case rdf.Seq:
			vals = val
		case rdf.Bag:
			vals = val
		case string:
			return []string{val}
		default:
			p.log("%s:%s has wrong data type", ns, name)
		}
	}
	list = make([]string, 0, len(vals))
	for _, str := range vals {
		switch str := str.Value.(type) {
		case string:
			list = append(list, str)
		default:
			p.log("%s:%s has wrong data type", ns, name)
		}
	}
	return list
}
