package xmp

import (
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/models/lr"
	"trimmer.io/go-xmp/xmp"
)

// LRHierarchicalSubject returns the values of the lr:HierarchicalSubject tag.
func (p *XMP) LRHierarchicalSubject() []metadata.Keyword { return p.lrHierarchicalSubject }

func (p *XMP) getLR() {
	var model *lr.Lightroom

	if p != nil && p.doc != nil {
		model = lr.FindModel(p.doc)
	}
	if model == nil || len(model.HierarchicalSubject) == 0 {
		return
	}
	for _, xkw := range model.HierarchicalSubject {
		p.lrHierarchicalSubject = append(p.lrHierarchicalSubject, strings.Split(xkw, "|"))
	}
}

// SetLRHierarchicalSubject sets the values of the lr:HierarchicalSubject tag.
func (p *XMP) SetLRHierarchicalSubject(v []metadata.Keyword) (err error) {
	var model *lr.Lightroom

	if model, err = lr.MakeModel(p.doc); err != nil {
		return fmt.Errorf("can't add lr model to XMP: %s", err)
	}
	if len(v) == len(p.lrHierarchicalSubject) {
		mismatch := false
		for i := range v {
			if !v[i].Equal(p.lrHierarchicalSubject[i]) {
				mismatch = true
				break
			}
		}
		if !mismatch {
			return nil
		}
	}
	p.lrHierarchicalSubject = v
	model.HierarchicalSubject = make(xmp.StringArray, len(v))
	for i := range v {
		model.HierarchicalSubject[i] = strings.Join(v[i], "|")
	}
	p.dirty = true
	return nil
}
