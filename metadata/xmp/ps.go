package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/ps"
)

// PSDateCreated returns the creation date from the XMP.
func (p *XMP) PSDateCreated() string {
	if p == nil || p.doc == nil {
		return ""
	}
	if model := ps.FindModel(p.doc); model != nil {
		if model.DateCreated != "" && !dateRE.MatchString(model.DateCreated) {
			p.log("PSDateCreated: invalid value: %q", model.DateCreated)
			return ""
		}
		return canonicalDate(model.DateCreated)
	}
	return ""
}

// SetPSDateCreated sets the creation date in the XMP.
func (p *XMP) SetPSDateCreated(v string) {
	model, err := ps.MakeModel(p.doc)
	if err != nil {
		p.log("XMP ps.MakeModel: %s", err)
		return
	}
	model.DateCreated = v
}
