package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/mwgrs"
)

// MWGRSNames returns the values of the mwg-rs:Name tags for face regions.
func (p *XMP) MWGRSNames() []string { return p.mwgrsNames }

func (p *XMP) getMWGRS() {
	var model *mwgrs.MWGRegions

	if p != nil && p.doc != nil {
		model = mwgrs.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	for _, r := range model.Regions.RegionList {
		if r.Type == "Face" && r.Name != "" {
			p.mwgrsNames = append(p.mwgrsNames, r.Name)
		}
	}
}

// Note: there is no SetMWGRSNames, because this library cannot set those tags.
