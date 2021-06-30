package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/iptc4ext"
)

// IPTCLocationCreated returns the value of the Iptc4xmpExt:LocationCreated tag.
func (p *XMP) IPTCLocationCreated() metadata.Location { return p.iptcLocationCreated }

// IPTCLocationsShown returns the values of the Iptc4xmpExt:LocationShown tag.
func (p *XMP) IPTCLocationsShown() []metadata.Location { return p.iptcLocationsShown }

func (p *XMP) getIPTC() {
	var model *iptc4ext.Iptc4xmpExt

	if p != nil && p.doc != nil {
		model = iptc4ext.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.iptcLocationCreated = xmpIPTCLocationToMetadata(model.LocationCreated)
	for _, xl := range model.LocationShown {
		p.iptcLocationsShown = append(p.iptcLocationsShown, xmpIPTCLocationToMetadata(xl))
	}
}
func xmpIPTCLocationToMetadata(xl *iptc4ext.Location) (ml metadata.Location) {
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

// SetIPTCLocationCreated sets the value of the Iptc4xmpExt:LocationCreated tag.
func (p *XMP) SetIPTCLocationCreated(v metadata.Location) (err error) {
	var model *iptc4ext.Iptc4xmpExt

	if model, err = iptc4ext.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add iptc4ext model to XMP: %s", err)
	}
	if v.Equal(p.iptcLocationCreated) {
		return nil
	}
	p.iptcLocationCreated = v
	if v.Empty() {
		model.LocationCreated = nil
	} else {
		metadataToXMPIPTCLocation(v, model.LocationCreated)
	}
	p.dirty = true
	return nil
}

// SetIPTCLocationsShown sets the values of the Iptc4xmpExt:LocationShown tag.
func (p *XMP) SetIPTCLocationsShown(v []metadata.Location) (err error) {
	var model *iptc4ext.Iptc4xmpExt

	if model, err = iptc4ext.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add exif model to XMP: %s", err)
	}
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
	model.LocationShown = make(iptc4ext.LocationArray, len(v))
	for i := range v {
		model.LocationShown[i] = new(iptc4ext.Location)
		metadataToXMPIPTCLocation(v[i], model.LocationShown[i])
	}
	p.dirty = true
	return nil
}

func metadataToXMPIPTCLocation(ml metadata.Location, xl *iptc4ext.Location) {
	xl.CountryCode = ml.CountryCode
	xl.CountryName = ml.CountryName
	xl.ProvinceState = ml.State
	xl.City = ml.City
	xl.Sublocation = ml.Sublocation
}
