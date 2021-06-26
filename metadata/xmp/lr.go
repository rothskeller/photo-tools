package xmp

import (
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
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
		parts := strings.Split(xkw, "|")
		var comps = make(metadata.Keyword, len(parts))
		for i, c := range parts {
			comps[i] = metadata.KeywordComponent{Word: c}
		}
		p.LRHierarchicalSubject = append(p.LRHierarchicalSubject, comps)
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
		var comps = make([]string, len(mkw))
		for i, c := range mkw {
			comps[i] = c.Word
		}
		hs = append(hs, strings.Join(comps, "|"))
	}
	if !stringSliceEqual(hs, model.HierarchicalSubject) {
		model.HierarchicalSubject = hs
		p.dirty = true
	}
}
