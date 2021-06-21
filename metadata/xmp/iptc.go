package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/iptc4ext"
	"trimmer.io/go-xmp/xmp"
)

// IPTCLocationsCreated returns the canonical IPTC location created.
func (p *XMP) IPTCLocationsCreated() []*metadata.Location {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := iptc4ext.FindModel(p.doc); model != nil {
		if model.LocationCreated != nil {
			return fromIPTCLocationLanguages(model.LocationCreated)
		}
	}
	return nil
}

// IPTCLocationsShown returns the canonical IPTC location shown.
func (p *XMP) IPTCLocationsShown() (locs []*metadata.Location) {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := iptc4ext.FindModel(p.doc); model != nil {
		for _, il := range model.LocationShown {
			locs = append(locs, fromIPTCLocationLanguages(il)...)
		}
		return locs
	}
	return nil
}

// fromIPTCLocationLanguages returns a set of metadata locations expressing all
// of the language variants in the IPTC location.
func fromIPTCLocationLanguages(il *iptc4ext.Location) (locs []*metadata.Location) {
	var languages []string
	if il == nil {
		return nil
	}
	for _, alt := range il.CountryName {
		languages = addUnique(languages, alt.Lang)
	}
	for _, alt := range il.ProvinceState {
		languages = addUnique(languages, alt.Lang)
	}
	for _, alt := range il.City {
		languages = addUnique(languages, alt.Lang)
	}
	for _, alt := range il.Sublocation {
		languages = addUnique(languages, alt.Lang)
	}
	locs = make([]*metadata.Location, len(languages))
	for idx, lang := range languages {
		var loc metadata.Location
		loc.CountryCode = il.CountryCode
		loc.CountryName = langOrDefault(il.CountryName, lang)
		loc.State = langOrDefault(il.ProvinceState, lang)
		loc.City = langOrDefault(il.City, lang)
		loc.Sublocation = langOrDefault(il.Sublocation, lang)
		locs[idx] = &loc
	}
	return locs
}
func addUnique(list []string, s string) []string {
	for _, item := range list {
		if item == s {
			return list
		}
	}
	return append(list, s)
}
func langOrDefault(alts xmp.AltString, lang string) string {
	if s := alts.Get(lang); s != "" {
		return s
	}
	return alts.Default()
}

// SetIPTCLocationCreated returns the canonical IPTC location created.
func (p *XMP) SetIPTCLocationCreated(loc *metadata.Location) {
	model, err := iptc4ext.MakeModel(p.doc)
	if err != nil {
		p.log("XMP iptc4ext.MakeModel: %s", err)
		return
	}
	model.LocationCreated = toIPTCLocation(loc)
}

// SetIPTCLocationShown returns the canonical IPTC location shown.
func (p *XMP) SetIPTCLocationShown(loc *metadata.Location) {
	model, err := iptc4ext.MakeModel(p.doc)
	if err != nil {
		p.log("XMP iptc4ext.MakeModel: %s", err)
		return
	}
	if loc.Valid() {
		model.LocationShown = []*iptc4ext.Location{toIPTCLocation(loc)}
	} else {
		model.LocationShown = nil
	}
}

// toIPTCLocation returns the metadata location converted to IPTC location form.
func toIPTCLocation(loc *metadata.Location) (il *iptc4ext.Location) {
	if !loc.Valid() {
		return nil
	}
	il = new(iptc4ext.Location)
	il.CountryCode = loc.CountryCode
	il.CountryName = xmp.NewAltString(loc.CountryName)
	il.ProvinceState = xmp.NewAltString(loc.State)
	il.City = xmp.NewAltString(loc.City)
	il.Sublocation = xmp.NewAltString(loc.Sublocation)
	return il
}
