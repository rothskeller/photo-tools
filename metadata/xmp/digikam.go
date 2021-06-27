package xmp

import (
	"strings"

	"github.com/rothskeller/photo-tools/metadata/xmp/models/digikam"
	"trimmer.io/go-xmp/xmp"
)

func (p *XMP) getDigiKam() {
	var model *digikam.DigiKam

	if p != nil && p.doc != nil {
		model = digikam.FindModel(p.doc)
	}
	if model == nil || len(model.TagsList) == 0 {
		return
	}
	for _, xkw := range model.TagsList {
		p.DigiKamTagsList = append(p.DigiKamTagsList, strings.Split(xkw, "/"))
	}
}

func (p *XMP) setDigiKam() {
	var (
		model *digikam.DigiKam
		tags  xmp.StringList
		err   error
	)
	if model, err = digikam.MakeModel(p.doc); err != nil {
		panic(err)
	}
	for _, mkw := range p.DigiKamTagsList {
		tags = append(tags, strings.Join(mkw, "/"))
	}
	if !stringSliceEqual(tags, model.TagsList) {
		model.TagsList = tags
		p.dirty = true
	}
}
