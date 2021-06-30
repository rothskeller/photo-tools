package xmp

import (
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/digikam"
	"trimmer.io/go-xmp/xmp"
)

// DigiKamTagsList returns the values of the digiKam:TagsList tag.
func (p *XMP) DigiKamTagsList() []metadata.Keyword { return p.digiKamTagsList }

func (p *XMP) getDigiKam() {
	var model *digikam.DigiKam

	if p != nil && p.doc != nil {
		model = digikam.FindModel(p.doc)
	}
	if model == nil || len(model.TagsList) == 0 {
		return
	}
	for _, xkw := range model.TagsList {
		p.digiKamTagsList = append(p.digiKamTagsList, strings.Split(xkw, "/"))
	}
}

// SetDigiKamTagsList sets the values of the digiKam:TagsList tag.
func (p *XMP) SetDigiKamTagsList(v []metadata.Keyword) error {
	var (
		model *digikam.DigiKam
		tags  xmp.StringList
		err   error
	)
	if model, err = digikam.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add digiKam model to XMP: %s", err)
	}
	for _, mkw := range v {
		tags = append(tags, strings.Join(mkw, "/"))
	}
	if !stringSliceEqual(tags, model.TagsList) {
		p.digiKamTagsList = v
		model.TagsList = tags
		p.dirty = true
	}
	return nil
}
