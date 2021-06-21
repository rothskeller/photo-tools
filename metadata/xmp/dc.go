package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/dc"
	"trimmer.io/go-xmp/xmp"
)

// DCCreator returns the list of Creators from the XMP.
func (p *XMP) DCCreator() []string {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := dc.FindModel(p.doc); model != nil {
		return model.Creator
	}
	return nil
}

// SetDCCreator sets the list of Creators in the XMP.
func (p *XMP) SetDCCreator(v []string) {
	model, err := dc.MakeModel(p.doc)
	if err != nil {
		p.log("XMP dc.MakeModel: %s", err)
		return
	}
	model.Creator = v
}

// DCDescription returns the Descriptions from the XMP, as an ordered list of
// alternatives, each of which is a [language, value] pair.  The first one is
// the default.
func (p *XMP) DCDescription() (descs [][]string) {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := dc.FindModel(p.doc); model != nil {
		for _, alt := range model.Description {
			descs = append(descs, []string{alt.Lang, alt.Value})
		}
		return descs
	}
	return nil
}

// SetDCDescription sets the Descriptions in the XMP.  Note that it sets the
// default language and removes any other language alternatives.
func (p *XMP) SetDCDescription(v string) {
	model, err := dc.MakeModel(p.doc)
	if err != nil {
		p.log("XMP dc.MakeModel: %s", err)
		return
	}
	if v != "" {
		model.Description = xmp.NewAltString(v)
	} else {
		model.Description = nil
	}
}

// DCSubject returns the list of Subjects from the XMP.
func (p *XMP) DCSubject() []string {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := dc.FindModel(p.doc); model != nil {
		return model.Subject
	}
	return nil
}

// SetDCSubject sets the list of Subjects in the XMP.
func (p *XMP) SetDCSubject(v []string) {
	model, err := dc.MakeModel(p.doc)
	if err != nil {
		p.log("XMP dc.MakeModel: %s", err)
		return
	}
	model.Subject = v
}

// DCTitle returns the Title from the XMP.  It returns a list of alternatives,
// each of which is a [language, value] pair.  The first one is the default.
func (p *XMP) DCTitle() (titles [][]string) {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := dc.FindModel(p.doc); model != nil {
		for _, alt := range model.Title {
			titles = append(titles, []string{alt.Lang, alt.Value})
		}
		return titles
	}
	return nil
}

// SetDCTitle sets the Title in the XMP.  Note that it sets the default language
// and removes any other language alternatives.
func (p *XMP) SetDCTitle(v string) {
	model, err := dc.MakeModel(p.doc)
	if err != nil {
		p.log("XMP dc.MakeModel: %s", err)
		return
	}
	model.Title = xmp.NewAltString(v)
}
