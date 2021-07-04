package xmp

import (
	"errors"

	"github.com/rothskeller/photo-tools/metadata/xmp/rdf"
)

const nsMP = "http://ns.microsoft.com/photo/1.2/"
const pfxMP = "MP"
const nsMPRI = "http://ns.microsoft.com/photo/1.2/t/RegionInfo#"
const pfxMPRI = "MPRI"
const nsMPReg = "http://ns.microsoft.com/photo/1.2/t/Region#"
const pfxMPReg = "MPReg"

// MPRegPersonDisplayNames returns the values of the MPReg:PersonDisplayName tag.
func (p *XMP) MPRegPersonDisplayNames() []string { return p.mpRegPersonDisplayNames }

func (p *XMP) getMP() {
	if val, ok := p.rdf.Properties[rdf.Name{Namespace: nsMP, Name: "RegionInfo"}]; ok {
		switch val := val.Value.(type) {
		case rdf.Struct:
			if bag, ok := val[rdf.Name{Namespace: nsMPRI, Name: "Regions"}]; ok {
				switch bag := bag.Value.(type) {
				case rdf.Bag:
					p.mpRegPersonDisplayNames = make([]string, 0, len(bag))
					for _, reg := range bag {
						switch reg := reg.Value.(type) {
						case rdf.Struct:
							if _, ok := reg[rdf.Name{Namespace: nsMPReg, Name: "Rectangle"}]; !ok {
								// digiKam tends to create region entries for anyone who has a face
								// in *any* photo, but omits the rectangle if they don't have a face
								// in *this* photo.  For our purposes, if there's no rectangle, it
								// doesn't count.
								continue
							}
							if name, ok := reg[rdf.Name{Namespace: nsMPReg, Name: "PersonDisplayName"}]; ok {
								switch name := name.Value.(type) {
								case string:
									p.mpRegPersonDisplayNames = append(p.mpRegPersonDisplayNames, name)
								default:
									p.log("MPReg:PersonDisplayName has wrong data type")
								}
							}
						default:
							p.log("MPRI:Regions has wrong data type")
						}
					}
				default:
					p.log("MPRI:Regions has wrong data type")
				}
			}
		default:
			p.log("MP:RegionInfo has wrong data type")
		}
	}
	p.rdf.RegisterNamespace(pfxMP, nsMP)
	p.rdf.RegisterNamespace(pfxMPRI, nsMPRI)
	p.rdf.RegisterNamespace(pfxMPReg, nsMPReg)
}

// SetMPRegPersonDisplayNames sets the values of the MPReg:PersonDisplayName
// tag.  Note however, that it cannot add any tags (because it doesn't have the
// information to do so completely); it can only remove them.
func (p *XMP) SetMPRegPersonDisplayNames(v []string) (err error) {
	var (
		regionInfo rdf.Struct
		bag        rdf.Bag
		nextInBag  int
		nextInV    int
	)
	if val, ok := p.rdf.Properties[rdf.Name{Namespace: nsMP, Name: "RegionInfo"}]; ok {
		regionInfo = val.Value.(rdf.Struct)
		if val, ok = regionInfo[rdf.Name{Namespace: nsMPRI, Name: "Regions"}]; ok {
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
		delete(p.rdf.Properties, rdf.Name{Namespace: nsMP, Name: "RegionInfo"})
		p.mpRegPersonDisplayNames = v
		p.dirty = true
		return nil
	}
	for _, oldv := range bag {
		if n, ok := oldv.Value.(rdf.Struct)[rdf.Name{Namespace: nsMPReg, Name: "PersonDisplayName"}]; !ok {
			bag[nextInBag] = oldv
			nextInBag++
		} else if nextInV < len(v) && n.Value.(string) == v[nextInV] {
			bag[nextInBag] = oldv
			nextInBag, nextInV = nextInBag+1, nextInV+1
		}
	}
	if nextInV < len(v) {
		goto NOADD
	}
	if nextInBag == len(bag) {
		return nil // nothing removed
	}
	regionInfo[rdf.Name{Namespace: nsMPRI, Name: "Regions"}] = rdf.Value{Value: bag[:nextInBag]}
	p.mpRegPersonDisplayNames = v
	p.dirty = true
	return nil
NOADD:
	return errors.New("MPRI:Regions: cannot add face region")
}
