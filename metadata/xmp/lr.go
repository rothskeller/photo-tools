package xmp

import (
	"github.com/rothskeller/photo-tools/metadata/xmp/models/lr"
)

// HierarchicalSubject returns the list of HierarchicalSubjects from the XMP.
func (p *XMP) HierarchicalSubject() []string {
	if p == nil || p.doc == nil {
		return nil
	}
	if model := lr.FindModel(p.doc); model != nil {
		return model.HierarchicalSubject
	}
	return nil
}

// SetHierarchicalSubject sets the list of HierarchicalSubjects in the XMP.
func (p *XMP) SetHierarchicalSubject(v []string) {
	model, err := lr.MakeModel(p.doc)
	if err != nil {
		p.log("XMP lr.MakeModel: %s", err)
		return
	}
	model.HierarchicalSubject = v
}
