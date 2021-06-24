package xmp

import (
	"reflect"

	"github.com/rothskeller/photo-tools/metadata/xmp/models/dc"
)

func (p *XMP) getDC() {
	var model *dc.DublinCore

	if p != nil && p.doc != nil {
		model = dc.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.DCCreator = xmpStringsToMetadata(model.Creator)
	p.DCDescription = xmpAltStringToMetadata(model.Description)
	p.DCSubject = xmpStringsToMetadata(model.Subject)
	p.DCTitle = xmpAltStringToMetadata(model.Title)
}

func (p *XMP) setDC() {
	var (
		model *dc.DublinCore
		err   error
	)
	if model, err = dc.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if creator := metadataToXMPStrings(p.DCCreator); !reflect.DeepEqual(creator, model.Creator) {
		model.Creator = creator
		p.dirty = true
	}
	if desc := metadataToXMPAltString(p.DCDescription); !reflect.DeepEqual(desc, model.Description) {
		model.Description = desc
		p.dirty = true
	}
	if subj := metadataToXMPStrings(p.DCSubject); !reflect.DeepEqual(subj, model.Subject) {
		model.Subject = subj
		p.dirty = true
	}
	if title := metadataToXMPAltString(p.DCTitle); !reflect.DeepEqual(title, model.Title) {
		model.Title = title
		p.dirty = true
	}
}
