package xmp

import (
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
		p.IPTCLocationsShown = append(p.IPTCLocationsShown, p.xmpIPTCLocationToMetadata(xl))
	}
}

func (p *XMP) setIPTC() {
	var (
		model *iptc4ext.Iptc4xmpExt
		err   error
	)
	if model, err = iptc4ext.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if p.IPTCLocationCreated.Empty() {
		if model.LocationCreated != nil {
			model.LocationCreated = nil
			p.dirty = true
		}
	} else {
		if model.LocationCreated == nil {
			model.LocationCreated = new(iptc4ext.Location)
		}
		p.metadataToXMPIPTCLocation(p.IPTCLocationCreated, model.LocationCreated)
	}
	if len(model.LocationShown) > len(p.IPTCLocationsShown) {
		model.LocationShown = model.LocationShown[:len(p.IPTCLocationsShown)]
		p.dirty = true
	}
	for len(model.LocationShown) < len(p.IPTCLocationsShown) {
		model.LocationShown = append(model.LocationShown, &iptc4ext.Location{})
	}
	for i := range p.IPTCLocationsShown {
		p.metadataToXMPIPTCLocation(p.IPTCLocationsShown[i], model.LocationShown[i])
	}
}

func (p *XMP) xmpIPTCLocationToMetadata(xl *iptc4ext.Location) (ml metadata.Location) {
	if xl == nil {
		return ml
	}
	ml.CountryCode = xl.CountryCode
	ml.CountryName = xl.CountryName
	ml.State = xl.ProvinceState
	ml.City = xl.City
	ml.Sublocation = xl.Sublocation
	return ml
}

func (p *XMP) metadataToXMPIPTCLocation(ml metadata.Location, xl *iptc4ext.Location) {
	if xl.CountryCode != ml.CountryCode {
		xl.CountryCode = ml.CountryCode
		p.dirty = true
	}
	if !metadata.EqualAltStrings(ml.CountryName, xl.CountryName) {
		xl.CountryName = ml.CountryName
		p.dirty = true
	}
	if !metadata.EqualAltStrings(ml.State, xl.ProvinceState) {
		xl.ProvinceState = ml.State
		p.dirty = true
	}
	if !metadata.EqualAltStrings(ml.City, xl.City) {
		xl.City = ml.City
		p.dirty = true
	}
	if !metadata.EqualAltStrings(ml.Sublocation, xl.Sublocation) {
		xl.Sublocation = ml.Sublocation
		p.dirty = true
	}
}
