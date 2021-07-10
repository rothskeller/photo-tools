package xmp

import (
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
)

const nsDigiKam = "http://www.digikam.org/ns/1.0/"
const pfxDigiKam = "digiKam"

// DigiKamTagsList returns the values of the digiKam:TagsList tag.
func (p *XMP) DigiKamTagsList() []metadata.Keyword { return p.digiKamTagsList }

func (p *XMP) getDigiKam() {
	var list = p.getStrings(p.rdf.Properties, pfxDigiKam, nsDigiKam, "TagsList")
	for _, xkw := range list {
		p.digiKamTagsList = append(p.digiKamTagsList, strings.Split(xkw, "/"))
	}
	p.rdf.RegisterNamespace(pfxDigiKam, nsDigiKam)
}

// SetDigiKamTagsList sets the values of the digiKam:TagsList tag.
func (p *XMP) SetDigiKamTagsList(v []metadata.Keyword) error {
	var tags, old []string

	old = p.getStrings(p.rdf.Properties, pfxDigiKam, nsDigiKam, "TagsList")
	for _, mkw := range v {
		tags = append(tags, strings.Join(mkw, "/"))
	}
	if !stringSliceEqual(tags, old) {
		p.digiKamTagsList = v
		p.setSeq(p.rdf.Properties, nsDigiKam, "TagsList", tags)
		p.dirty = true
	}
	return nil
}
