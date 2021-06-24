package xmp

import (
	"reflect"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/iptc4ext"
)

func (p *XMP) getIPTC() {
	var model *iptc4ext.Iptc4xmpExt

	if p != nil && p.doc != nil {
		model = iptc4ext.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.IPTCLocationCreated = p.xmpIPTCLocationToMetadata(model.LocationCreated)
	p.IPTCLocationsShown = p.xmpIPTCLocationsToMetadata(model.LocationShown)
}

func (p *XMP) setIPTC() {
	var (
		model *iptc4ext.Iptc4xmpExt
		err   error
	)
	if model, err = iptc4ext.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if loc := metadataToXMPIPTCLocation(p.IPTCLocationCreated); !reflect.DeepEqual(loc, model.LocationCreated) {
		model.LocationCreated = loc
		p.dirty = true
	}
	if locs := metadataToXMPIPTCLocations(p.IPTCLocationsShown); !reflect.DeepEqual(locs, model.LocationShown) {
		model.LocationShown = locs
		p.dirty = true
	}
}

func (p *XMP) xmpIPTCLocationToMetadata(xl *iptc4ext.Location) (ml metadata.Multilingual) {
	if xl == nil {
		return nil
	}
	// What languages are used in the xl?
	var langs []string
	var seen = map[string]bool{}
	for _, alt := range xl.CountryName {
		if !seen[alt.Lang] {
			langs = append(langs, alt.Lang)
			seen[alt.Lang] = true
		}
	}
	for _, alt := range xl.ProvinceState {
		if !seen[alt.Lang] {
			langs = append(langs, alt.Lang)
			seen[alt.Lang] = true
		}
	}
	for _, alt := range xl.City {
		if !seen[alt.Lang] {
			langs = append(langs, alt.Lang)
			seen[alt.Lang] = true
		}
	}
	for _, alt := range xl.Sublocation {
		if !seen[alt.Lang] {
			langs = append(langs, alt.Lang)
			seen[alt.Lang] = true
		}
	}
	ml = make(metadata.Multilingual, len(langs))
	for i, lang := range langs {
		var ccode, cname, state, city, subloc *metadata.String
		if xl.CountryCode != "" {
			ccode = metadata.NewString(xl.CountryCode)
		}
		if s := xl.CountryName.Get(lang); s != "" {
			cname = metadata.NewString(s)
		} else if s := xl.CountryName.Default(); s != "" {
			cname = metadata.NewString(s)
		}
		if s := xl.ProvinceState.Get(lang); s != "" {
			state = metadata.NewString(s)
		} else if s := xl.ProvinceState.Default(); s != "" {
			state = metadata.NewString(s)
		}
		if s := xl.City.Get(lang); s != "" {
			city = metadata.NewString(s)
		} else if s := xl.City.Default(); s != "" {
			city = metadata.NewString(s)
		}
		if s := xl.Sublocation.Get(lang); s != "" {
			subloc = metadata.NewString(s)
		} else if s := xl.Sublocation.Default(); s != "" {
			subloc = metadata.NewString(s)
		}
		var loc metadata.Location
		if err := loc.ParseComponents(ccode, cname, state, city, subloc); err != nil {
			p.log("invalid XMP IPTC location")
			return nil
		}
		ml[i] = &metadata.LangDatum{Lang: lang, Metadatum: &loc}
	}
	return ml
}

func metadataToXMPIPTCLocation(ml metadata.Multilingual) (xl *iptc4ext.Location) {
	if len(ml) == 0 {
		return nil
	}
	xl = new(iptc4ext.Location)
	for i, ld := range ml {
		ccode, cname, state, city, subloc := ld.Metadatum.(*metadata.Location).AsComponents()
		if i == 0 {
			xl.CountryCode = ccode.String()
			xl.CountryName.AddDefault(ld.Lang, cname.String())
			xl.ProvinceState.AddDefault(ld.Lang, state.String())
			xl.City.AddDefault(ld.Lang, city.String())
			xl.Sublocation.AddDefault(ld.Lang, subloc.String())
		} else {
			if cname.String() != xl.CountryName.Default() {
				xl.CountryName.Add(ld.Lang, cname.String())
			}
			if state.String() != xl.ProvinceState.Default() {
				xl.ProvinceState.Add(ld.Lang, state.String())
			}
			if city.String() != xl.City.Default() {
				xl.City.Add(ld.Lang, city.String())
			}
			if subloc.String() != xl.Sublocation.Default() {
				xl.Sublocation.Add(ld.Lang, subloc.String())
			}
		}
	}
	return xl
}

func (p *XMP) xmpIPTCLocationsToMetadata(xl iptc4ext.LocationArray) (ml []metadata.Multilingual) {
	for _, loc := range xl {
		var mloc = p.xmpIPTCLocationToMetadata(loc)
		if len(mloc) != 0 {
			ml = append(ml, mloc)
		}
	}
	return ml
}

func metadataToXMPIPTCLocations(ml []metadata.Multilingual) (xl iptc4ext.LocationArray) {
	for _, loc := range ml {
		if len(loc) != 0 {
			xl = append(xl, metadataToXMPIPTCLocation(loc))
		}
	}
	return xl
}
