package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/dc"
)

// DCCreator returns the values of the dc:creator tag.
func (p *XMP) DCCreator() []string { return p.dcCreator }

// DCDescription returns the values of the dc:description tag.
func (p *XMP) DCDescription() metadata.AltString { return p.dcDescription }

// DCSubject returns the values of the dc:subject tag.
func (p *XMP) DCSubject() []string { return p.dcSubject }

// DCTitle returns the values of the dc:title tag.
func (p *XMP) DCTitle() metadata.AltString { return p.dcTitle }

func (p *XMP) getDC() {
	var model *dc.DublinCore

	if p != nil && p.doc != nil {
		model = dc.FindModel(p.doc)
	}
	if model == nil {
		return
	}
	p.dcCreator = model.Creator
	p.dcDescription = model.Description
	p.dcSubject = model.Subject
	p.dcTitle = model.Title
}

// SetDCCreator sets the values of the dc:creator tag.
func (p *XMP) SetDCCreator(v []string) error {
	var (
		model *dc.DublinCore
		err   error
	)
	if model, err = dc.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add dc model to XMP: %s", err)
	}
	if !stringSliceEqual(v, p.dcCreator) {
		p.dcCreator = v
		model.Creator = v
		p.dirty = true
	}
	return nil
}

// SetDCDescription sets the values of the dc:description tag.
func (p *XMP) SetDCDescription(v metadata.AltString) error {
	var (
		model *dc.DublinCore
		err   error
	)
	if model, err = dc.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add model to XMP: %s", err)
	}
	if !metadata.EqualAltStrings(v, p.dcDescription) {
		p.dcDescription = v
		model.Description = v
		p.dirty = true
	}
	return nil
}

// SetDCSubject sets the values of the dc:subject tag.
func (p *XMP) SetDCSubject(v []string) error {
	var (
		model *dc.DublinCore
		err   error
	)
	if model, err = dc.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add model to XMP: %s", err)
	}
	if !stringSliceEqual(v, p.dcSubject) {
		p.dcSubject = v
		model.Subject = v
		p.dirty = true
	}
	return nil
}

// SetDCTitle sets the values of the dc:title tag.
func (p *XMP) SetDCTitle(v metadata.AltString) error {
	var (
		model *dc.DublinCore
		err   error
	)
	if model, err = dc.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add model to XMP: %s", err)
	}
	if !metadata.EqualAltStrings(v, p.dcTitle) {
		p.dcTitle = v
		model.Title = v
		p.dirty = true
	}
	return nil
}
