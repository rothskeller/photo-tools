package xmp

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata/xmp/rdf"
)

const nsMWGRS = "http://www.metadataworkinggroup.com/schemas/regions/"
const pfxMWGRS = "mwg-rs"

// MWGRSNames returns the values of the mwg-rs:Name tags for face regions.
func (p *XMP) MWGRSNames() []string { return p.mwgrsNames }

func (p *XMP) getMWGRS() {
	if val, ok := p.rdf.Properties[rdf.Name{Namespace: nsMWGRS, Name: "Regions"}]; ok {
		switch val := val.Value.(type) {
		case rdf.Struct:
			if bag, ok := val[rdf.Name{Namespace: nsMWGRS, Name: "RegionList"}]; ok {
				switch bag := bag.Value.(type) {
				case rdf.Bag:
					p.mwgrsNames = make([]string, 0, len(bag))
					for _, reg := range bag {
						switch reg := reg.Value.(type) {
						case rdf.Struct:
							if typ, ok := reg[rdf.Name{Namespace: nsMWGRS, Name: "Type"}]; ok {
								switch typ := typ.Value.(type) {
								case string:
									if typ != "Face" {
										continue
									}
								default:
									p.log("mwg-rs:Type has wrong data type")
								}
							}
							if name, ok := reg[rdf.Name{Namespace: nsMWGRS, Name: "Name"}]; ok {
								switch name := name.Value.(type) {
								case string:
									p.mwgrsNames = append(p.mwgrsNames, name)
								default:
									p.log("mwg-rs:Name has wrong data type")
								}
							}
						default:
							p.log("mwg-rs:RegionList has wrong data type")
						}
					}
				default:
					p.log("mwg-rs:RegionList has wrong data type")
				}
			}
		default:
			p.log("mwg-rs:Regions has wrong data type")
		}
	}
	p.rdf.RegisterNamespace(pfxMWGRS, nsMWGRS)
}

// SetMWGRSNames sets the values of the mwg-rs:Name tag.  Note however, that it
// cannot add any tags (because it doesn't have the information to do so
// completely); it can only remove them.
func (p *XMP) SetMWGRSNames(v []string) (err error) {
	var (
		regions   rdf.Struct
		bag       rdf.Bag
		nextInBag int
		nextInV   int
	)
	if val, ok := p.rdf.Properties[rdf.Name{Namespace: nsMWGRS, Name: "Regions"}]; ok {
		regions = val.Value.(rdf.Struct)
		if val, ok = regions[rdf.Name{Namespace: nsMWGRS, Name: "RegionList"}]; ok {
			bag = val.Value.(rdf.Bag)
		}
	}
	if len(bag) == 0 && len(v) == 0 {
		return nil
	}
	if len(bag) == 0 {
		goto NOADD
	}
	if len(v) == 0 {
		delete(p.rdf.Properties, rdf.Name{Namespace: nsMWGRS, Name: "Regions"})
		p.mwgrsNames = v
		p.dirty = true
		return nil
	}
	for _, oldv := range bag {
		if t, ok := oldv.Value.(rdf.Struct)[rdf.Name{Namespace: nsMWGRS, Name: "Type"}]; !ok || t.Value.(string) != "Face" {
			bag[nextInBag] = oldv
			nextInBag++
		} else if n, ok := oldv.Value.(rdf.Struct)[rdf.Name{Namespace: nsMWGRS, Name: "Name"}]; ok {
			if nextInV < len(v) && n.Value.(string) == v[nextInV] {
				bag[nextInBag] = oldv
				nextInBag, nextInV = nextInBag+1, nextInV+1
			}
		} else {
			bag[nextInBag] = oldv
			nextInBag++
		}
	}
	if nextInV < len(v) {
		goto NOADD
	}
	if nextInBag == len(bag) {
		return nil // nothing removed
	}
	regions[rdf.Name{Namespace: nsMWGRS, Name: "RegionList"}] = rdf.Value{Value: bag[:nextInBag]}
	p.mwgrsNames = v
	p.dirty = true
	return nil
NOADD:
	return errors.New("mwg-rs:RegionList: cannot add face region")
}
