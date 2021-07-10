package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

const nsDC = "http://purl.org/dc/elements/1.1/"
const pfxDC = "dc"

// DCCreator returns the values of the dc:creator tag.
func (p *XMP) DCCreator() []string { return p.dcCreator }

// DCDescription returns the values of the dc:description tag.
func (p *XMP) DCDescription() metadata.AltString { return p.dcDescription }

// DCSubject returns the values of the dc:subject tag.
func (p *XMP) DCSubject() []string { return p.dcSubject }

// DCTitle returns the values of the dc:title tag.
func (p *XMP) DCTitle() metadata.AltString { return p.dcTitle }

func (p *XMP) getDC() {
	p.dcCreator = p.getStrings(p.rdf.Properties, pfxDC, nsDC, "creator")
	p.dcDescription = p.getAlt(p.rdf.Properties, pfxDC, nsDC, "description")
	p.dcSubject = p.getStrings(p.rdf.Properties, pfxDC, nsDC, "subject")
	p.dcTitle = p.getAlt(p.rdf.Properties, pfxDC, nsDC, "title")
	p.rdf.RegisterNamespace(pfxDC, nsDC)
}

// SetDCCreator sets the values of the dc:creator tag.
func (p *XMP) SetDCCreator(v []string) error {
	if !stringSliceEqual(v, p.dcCreator) {
		p.dcCreator = v
		p.setSeq(p.rdf.Properties, nsDC, "creator", v)
	}
	return nil
}

// SetDCDescription sets the values of the dc:description tag.
func (p *XMP) SetDCDescription(v metadata.AltString) error {
	if !v.Equal(p.dcDescription) {
		p.dcDescription = v
		p.setAlt(p.rdf.Properties, nsDC, "description", v)
		p.dirty = true
	}
	return nil
}

// SetDCSubject sets the values of the dc:subject tag.
func (p *XMP) SetDCSubject(v []string) error {
	if !stringSliceEqual(v, p.dcSubject) {
		p.dcSubject = v
		p.setBag(p.rdf.Properties, nsDC, "subject", v)
	}
	return nil
}

// SetDCTitle sets the values of the dc:title tag.
func (p *XMP) SetDCTitle(v metadata.AltString) error {
	if !v.Equal(p.dcTitle) {
		p.dcTitle = v
		p.setAlt(p.rdf.Properties, nsDC, "title", v)
		p.dirty = true
	}
	return nil
}
