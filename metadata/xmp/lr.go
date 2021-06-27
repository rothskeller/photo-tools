package xmp

import (
	"strings"

	"github.com/rothskeller/photo-tools/metadata/xmp/models/lr"
	"trimmer.io/go-xmp/xmp"
)

func (p *XMP) getLR() {
	var model *lr.Lightroom

	if p != nil && p.doc != nil {
		model = lr.FindModel(p.doc)
	}
	if model == nil || len(model.HierarchicalSubject) == 0 {
		return
	}
	for _, xkw := range model.HierarchicalSubject {
		p.LRHierarchicalSubject = append(p.LRHierarchicalSubject, strings.Split(xkw, "|"))
	}
}

func (p *XMP) setLR() {
	var (
		model *lr.Lightroom
		hs    xmp.StringArray
		err   error
	)
	if model, err = lr.MakeModel(p.doc); err != nil {
		panic(err)
	}
	for _, mkw := range p.LRHierarchicalSubject {
		hs = append(hs, strings.Join(mkw, "|"))
	}
	if !stringSliceEqual(hs, model.HierarchicalSubject) {
		model.HierarchicalSubject = hs
		p.dirty = true
	}
}
