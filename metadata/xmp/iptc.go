package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/rdf"
)

const nsIPTC = "http://iptc.org/std/Iptc4xmpExt/2008-02-29/"
const pfxIPTC = "Iptc4xmpExt"

// IPTCLocationCreated returns the value of the Iptc4xmpExt:LocationCreated tag.
func (p *XMP) IPTCLocationCreated() metadata.Location { return p.iptcLocationCreated }

// IPTCLocationsShown returns the values of the Iptc4xmpExt:LocationShown tag.
func (p *XMP) IPTCLocationsShown() []metadata.Location { return p.iptcLocationsShown }

func (p *XMP) getIPTC() {
	if val, ok := p.rdf.Properties[rdf.Name{Namespace: nsIPTC, Name: "LocationCreated"}]; ok {
		switch val := val.Value.(type) {
		case rdf.Struct:
			p.iptcLocationCreated = p.getLocationFromStruct(val)
		default:
			p.log("Iptc4xmpExt:LocationCreated has wrong data type")
		}
	}
	if val, ok := p.rdf.Properties[rdf.Name{Namespace: nsIPTC, Name: "LocationShown"}]; ok {
		switch val := val.Value.(type) {
		case rdf.Seq:
			p.iptcLocationsShown = make([]metadata.Location, 0, len(val))
			for _, loc := range val {
				switch loc := loc.Value.(type) {
				case rdf.Struct:
					p.iptcLocationsShown = append(p.iptcLocationsShown, p.getLocationFromStruct(loc))
				default:
					p.log("Iptc4xmpExt:LocationCreated has wrong data type")
				}
			}
		default:
			p.log("Iptc4xmpExt:LocationShown has wrong data type")
		}
	}
}

func (p *XMP) getLocationFromStruct(str rdf.Struct) (loc metadata.Location) {
	loc.CountryCode = p.getString(str, pfxIPTC, nsIPTC, "CountryCode")
	loc.CountryName = p.getAlt(str, pfxIPTC, nsIPTC, "CountryName")
	loc.State = p.getAlt(str, pfxIPTC, nsIPTC, "ProvinceState")
	loc.City = p.getAlt(str, pfxIPTC, nsIPTC, "City")
	loc.Sublocation = p.getAlt(str, pfxIPTC, nsIPTC, "Sublocation")
	return loc
}

// SetIPTCLocationCreated sets the value of the Iptc4xmpExt:LocationCreated tag.
func (p *XMP) SetIPTCLocationCreated(v metadata.Location) (err error) {
	if v.Equal(p.iptcLocationCreated) {
		return nil
	}
	p.iptcLocationCreated = v
	if v.Empty() {
		delete(p.rdf.Properties, rdf.Name{Namespace: nsIPTC, Name: "LocationCreated"})
	} else {
		p.rdf.Properties[rdf.Name{Namespace: nsIPTC, Name: "LocationCreated"}] = p.getStructFromLocation(v)
	}
	p.dirty = true
	return nil
}

// SetIPTCLocationsShown sets the values of the Iptc4xmpExt:LocationShown tag.
func (p *XMP) SetIPTCLocationsShown(v []metadata.Location) (err error) {
	if len(v) == len(p.iptcLocationsShown) {
		mismatch := false
		for i := range v {
			if !v[i].Equal(p.iptcLocationsShown[i]) {
				mismatch = true
				break
			}
		}
		if !mismatch {
			return nil
		}
	}
	p.iptcLocationsShown = v
	if len(v) == 0 {
		delete(p.rdf.Properties, rdf.Name{Namespace: nsIPTC, Name: "LocationShown"})
	} else {
		var lis = make([]rdf.Value, len(v))
		for i := range v {
			lis[i] = p.getStructFromLocation(v[i])
		}
		p.rdf.Properties[rdf.Name{Namespace: nsIPTC, Name: "LocationShown"}] = rdf.Value{Value: rdf.Seq(lis)}
	}
	p.dirty = true
	return nil
}

func (p *XMP) getStructFromLocation(loc metadata.Location) rdf.Value {
	var str rdf.Struct

	p.setString(str, nsIPTC, "CountryCode", loc.CountryCode)
	p.setAlt(str, nsIPTC, "CountryName", loc.CountryName)
	p.setAlt(str, nsIPTC, "ProvinceState", loc.State)
	p.setAlt(str, nsIPTC, "City", loc.City)
	p.setAlt(str, nsIPTC, "Sublocation", loc.Sublocation)
	return rdf.Value{Value: str}
}
