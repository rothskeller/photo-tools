package xmp

import (
	"reflect"

	"github.com/rothskeller/photo-tools/metadata/xmp/models/dc"
	"trimmer.io/go-xmp/xmp"
)

func (p *XMP) getDC() {
	var model *dc.DublinCore

	if p != nil && p.doc != nil {
		model = dc.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.DCCreator = model.Creator
	p.DCDescription = xmpAltStringToMetadata(model.Description)
	p.DCSubject = model.Subject
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
	if !reflect.DeepEqual(xmp.StringList(p.DCCreator), model.Creator) {
		model.Creator = p.DCCreator
		p.dirty = true
	}
	if desc := metadataToXMPAltString(p.DCDescription); !reflect.DeepEqual(desc, model.Description) {
		model.Description = desc
		p.dirty = true
	}
	if !reflect.DeepEqual(xmp.StringArray(p.DCSubject), model.Subject) {
		model.Subject = p.DCSubject
		p.dirty = true
	}
	if title := metadataToXMPAltString(p.DCTitle); !reflect.DeepEqual(title, model.Title) {
		model.Title = title
		p.dirty = true
	}
}
