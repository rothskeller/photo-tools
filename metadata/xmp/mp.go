package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/mp"
)

// MPRegPersonDisplayNames returns the values of the MPReg:PersonDisplayName tag.
func (p *XMP) MPRegPersonDisplayNames() []string { return p.mpRegPersonDisplayNames }

func (p *XMP) getMP() {
	var model *mp.MPInfo

	if p != nil && p.doc != nil {
		model = mp.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	for _, r := range model.RegionInfo.Regions {
		p.mpRegPersonDisplayNames = append(p.mpRegPersonDisplayNames, r.PersonDisplayName)
	}
}

// Note that there is no SetMPRegPersonDisplayNames, because this library can't
// change those tags.
