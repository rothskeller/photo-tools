package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/ps"
)

// PSDateCreated returns the value of the photoshop:DateCreated tag.
func (p *XMP) PSDateCreated() metadata.DateTime { return p.psDateCreated }

func (p *XMP) getPS() {
	var model *ps.PhotoshopInfo

	if p != nil && p.doc != nil {
		model = ps.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.xmpDateTimeToMetadata(model.DateCreated, &p.psDateCreated)
}

// SetPSDateCreated sets the value of the photoshop:DateCreated tag.
func (p *XMP) SetPSDateCreated(v metadata.DateTime) (err error) {
	var model *ps.PhotoshopInfo

	if model, err = ps.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add photoshop model to XMP: %s", err)
	}
	if v.Equivalent(p.psDateCreated) {
		return nil
	}
	p.psDateCreated = v
	model.DateCreated = v.String()
	p.dirty = true
	return nil
}
