package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/mp"
)

func (p *XMP) getMP() {
	var model *mp.MPInfo

	if p != nil && p.doc != nil {
		model = mp.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	for _, r := range model.RegionInfo.Regions {
		p.MPFaces = append(p.MPFaces, r.PersonDisplayName)
	}
}

func (p *XMP) setMP() {
	var (
		model *mp.MPInfo
		err   error
	)
	if model, err = mp.MakeModel(p.doc); err != nil {
		panic(err)
	}
	j := 0
	for _, r := range model.RegionInfo.Regions {
		found := false
		for _, f := range p.MPFaces {
			if f == r.PersonDisplayName {
				found = true
				break
			}
		}
		if found {
			model.RegionInfo.Regions[j] = r
			j++
		}
	}
	if j < len(model.RegionInfo.Regions) {
		model.RegionInfo.Regions = model.RegionInfo.Regions[:j]
		p.dirty = true
	}
}
