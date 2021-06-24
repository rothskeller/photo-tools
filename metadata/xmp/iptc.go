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
	for _, xl := range model.LocationShown {
		var ml = p.xmpIPTCLocationToMetadata(xl)
		if len(ml) != 0 {
			p.IPTCLocationsShown = append(p.IPTCLocationsShown, ml)
		}
	}
}

func (p *XMP) setIPTC() {
	var (
		model *iptc4ext.Iptc4xmpExt
		shown iptc4ext.LocationArray
		err   error
	)
	if model, err = iptc4ext.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if loc := metadataToXMPIPTCLocation(p.IPTCLocationCreated); !reflect.DeepEqual(loc, model.LocationCreated) {
		model.LocationCreated = loc
		p.dirty = true
	}
	for _, ml := range p.IPTCLocationsShown {
		if len(ml) != 0 {
			shown = append(shown, metadataToXMPIPTCLocation(ml))
		}
	}
	if !reflect.DeepEqual(shown, model.LocationShown) {
		model.LocationShown = shown
		p.dirty = true
	}
}

func (p *XMP) xmpIPTCLocationToMetadata(xl *iptc4ext.Location) (mls []metadata.Location) {
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
	mls = make([]metadata.Location, len(langs))
	for i, lang := range langs {
		mls[i].Lang = lang
		var cname, state, city, subloc string
		if s := xl.CountryName.Get(lang); s != "" {
			cname = s
		} else {
			cname = xl.CountryName.Default()
		}
		if s := xl.ProvinceState.Get(lang); s != "" {
			state = s
		} else {
			state = xl.ProvinceState.Default()
		}
		if s := xl.City.Get(lang); s != "" {
			city = s
		} else {
			city = xl.City.Default()
		}
		if s := xl.Sublocation.Get(lang); s != "" {
			subloc = s
		} else {
			subloc = xl.Sublocation.Default()
		}
		if err := mls[i].ParseComponents(xl.CountryCode, cname, state, city, subloc); err != nil {
			p.log("invalid XMP IPTC location")
			return nil
		}
	}
	return mls
}

func metadataToXMPIPTCLocation(mls []metadata.Location) (xl *iptc4ext.Location) {
	if len(mls) == 0 {
		return nil
	}
	xl = new(iptc4ext.Location)
	for i, ml := range mls {
		if i == 0 {
			xl.CountryCode = ml.CountryCode
			xl.CountryName.AddDefault(ml.Lang, ml.CountryName)
			xl.ProvinceState.AddDefault(ml.Lang, ml.State)
			xl.City.AddDefault(ml.Lang, ml.City)
			xl.Sublocation.AddDefault(ml.Lang, ml.Sublocation)
		} else {
			if ml.CountryName != xl.CountryName.Default() {
				xl.CountryName.Add(ml.Lang, ml.CountryName)
			}
			if ml.State != xl.ProvinceState.Default() {
				xl.ProvinceState.Add(ml.Lang, ml.State)
			}
			if ml.City != xl.City.Default() {
				xl.City.Add(ml.Lang, ml.City)
			}
			if ml.Sublocation != xl.Sublocation.Default() {
				xl.Sublocation.Add(ml.Lang, ml.Sublocation)
			}
		}
	}
	return xl
}
