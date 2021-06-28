package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/mwgrs"
)

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
			p.MWGRSFaces = append(p.MWGRSFaces, r.Name)
		}
	}
}

func (p *XMP) setMWGRS() {
	var (
		model *mwgrs.MWGRegions
		err   error
	)
	if model, err = mwgrs.MakeModel(p.doc); err != nil {
		panic(err)
	}
	j := 0
	for _, r := range model.Regions.RegionList {
		if r.Type != "Face" {
			model.Regions.RegionList[j] = r
			j++
			continue
		}
		found := false
		for _, f := range p.MWGRSFaces {
			if f == r.Name {
				found = true
				break
			}
		}
		if found {
			model.Regions.RegionList[j] = r
			j++
		}
	}
	if j < len(model.Regions.RegionList) {
		model.Regions.RegionList = model.Regions.RegionList[:j]
		p.dirty = true
	}
}
