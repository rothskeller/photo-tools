package xmp

import (
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/xmp/rdf"
)

const nsLR = "http://ns.adobe.com/lightroom/1.0/"
const pfxLR = "lr"

// LRHierarchicalSubject returns the values of the lr:HierarchicalSubject tag.
func (p *XMP) LRHierarchicalSubject() []metadata.Keyword { return p.lrHierarchicalSubject }

func (p *XMP) getLR() {
	var kws = p.getStrings(p.rdf.Properties, pfxLR, nsLR, "hierarchicalSubject")
	for _, kw := range kws {
		p.lrHierarchicalSubject = append(p.lrHierarchicalSubject, strings.Split(kw, "|"))
	}
	p.rdf.RegisterNamespace(pfxLR, nsLR)
}

// SetLRHierarchicalSubject sets the values of the lr:HierarchicalSubject tag.
func (p *XMP) SetLRHierarchicalSubject(v []metadata.Keyword) (err error) {
	var tags, old []string

	// Earlier versions of this tool incorrectly created the tag with a
	// capital H.  Remove those incorrect tags.
	delete(p.rdf.Properties, rdf.Name{Namespace: nsLR, Name: "HierarchicalSubject"})

	old = p.getStrings(p.rdf.Properties, pfxLR, nsLR, "hierarchicalSubject")
	for _, mkw := range v {
		tags = append(tags, strings.Join(mkw, "|"))
	}
	if !stringSliceEqual(tags, old) {
		p.digiKamTagsList = v
		p.setBag(p.rdf.Properties, nsLR, "hierarchicalSubject", tags)
		p.dirty = true
	}
	return nil
}
