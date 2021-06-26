package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
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
	p.DCDescription = model.Description
	p.DCSubject = model.Subject
	p.DCTitle = model.Title
}

func (p *XMP) setDC() {
	var (
		model *dc.DublinCore
		err   error
	)
	if model, err = dc.MakeModel(p.doc); err != nil {
		panic(err)
	}
	if !stringSliceEqual(xmp.StringList(p.DCCreator), model.Creator) {
		model.Creator = p.DCCreator
		p.dirty = true
	}
	if !metadata.EqualAltStrings(p.DCDescription, model.Description) {
		model.Description = p.DCDescription
		p.dirty = true
	}
	if !stringSliceEqual(xmp.StringArray(p.DCSubject), model.Subject) {
		model.Subject = p.DCSubject
		p.dirty = true
	}
	if !metadata.EqualAltStrings(p.DCTitle, model.Title) {
		model.Title = p.DCTitle
		p.dirty = true
	}
}
