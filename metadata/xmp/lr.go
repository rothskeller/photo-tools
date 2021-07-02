package xmp

import (
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
)

const nsLR = "http://ns.adobe.com/lightroom/1.0/"
const pfxLR = "lr"

// LRHierarchicalSubject returns the values of the lr:HierarchicalSubject tag.
func (p *XMP) LRHierarchicalSubject() []metadata.Keyword { return p.lrHierarchicalSubject }

func (p *XMP) getLR() {
	var kws = p.getStrings(p.rdf.Properties, pfxLR, nsLR, "HierarchicalSubject")
	for _, kw := range kws {
		p.lrHierarchicalSubject = append(p.lrHierarchicalSubject, strings.Split(kw, "|"))
	}
}

// SetLRHierarchicalSubject sets the values of the lr:HierarchicalSubject tag.
func (p *XMP) SetLRHierarchicalSubject(v []metadata.Keyword) (err error) {
	var tags, old []string

	old = p.getStrings(p.rdf.Properties, pfxLR, nsLR, "HierarchicalSubject")
	for _, mkw := range v {
		tags = append(tags, strings.Join(mkw, "|"))
	}
	if !stringSliceEqual(tags, old) {
		p.digiKamTagsList = v
		p.setBag(p.rdf.Properties, nsLR, "HierarchicalSubject", tags)
		p.dirty = true
	}
	return nil
}
