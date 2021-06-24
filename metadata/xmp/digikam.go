package xmp

import (
	"reflect"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
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
		parts := strings.Split(xkw, "/")
		var comps = make(metadata.Keyword, len(parts))
		for i, c := range parts {
			comps[i] = metadata.KeywordComponent{Word: c}
		}
		p.DigiKamTagsList = append(p.DigiKamTagsList, comps)
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
		var comps = make([]string, len(mkw))
		for i, c := range mkw {
			comps[i] = c.Word
		}
		tags = append(tags, strings.Join(comps, "/"))
	}
	if !reflect.DeepEqual(tags, model.TagsList) {
		model.TagsList = tags
		p.dirty = true
	}
}
