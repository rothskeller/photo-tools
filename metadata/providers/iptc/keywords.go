package iptc

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/iim"
)

const (
	idKeyword     uint16 = 0x0219
	maxKeywordLen        = 64
)

// getKeywords reads the values of the Keywords field from the IIM.
func (p *Provider) getKeywords() (err error) {
	for _, ds := range p.iim[idKeyword] {
		if keyword, err := getString(ds); err == nil {
			p.keywords = append(p.keywords, keyword)
		} else {
			return fmt.Errorf("Keyword: %s", err)
		}
	}
	return nil
}

// Keywords returns the values of the Keywords field.
func (p *Provider) Keywords() (values []metadata.HierValue) {
	values = make([]metadata.HierValue, len(p.keywords))
	for i := range p.keywords {
		values[i] = metadata.HierValue{p.keywords[i]}
	}
	return values
}

// KeywordsTags returns a list of tag names for the Keywords field, and
// a parallel list of values held by those tags.
func (p *Provider) KeywordsTags() (tags []string, values []metadata.HierValue) {
	tags = make([]string, len(p.keywords))
	values = make([]metadata.HierValue, len(p.keywords))
	for i := range p.keywords {
		tags[i] = "IPTC Keyword"
		values[i] = metadata.HierValue{p.keywords[i]}
	}
	return tags, values
}

// SetKeywords sets the values of the Keywords field.
func (p *Provider) SetKeywords(values []metadata.HierValue) error {
	if len(values) == 0 {
		p.keywords = nil
		if _, ok := p.iim[idKeyword]; ok {
			delete(p.iim, idKeyword)
			p.dirty = true
		}
		return nil
	}
	var vmap = make(map[string]bool)
	var vlist []string
	for _, hv := range values {
		kw := hv[len(hv)-1]
		if len(kw) > maxKeywordLen {
			kw = kw[:maxKeywordLen]
		}
		if _, ok := vmap[kw]; !ok {
			vmap[kw] = false
			vlist = append(vlist, kw)
		}
	}
	var changed = false
	for _, kw := range p.keywords {
		if _, ok := vmap[kw]; ok {
			vmap[kw] = true
		} else {
			changed = true
		}
	}
	for _, seen := range vmap {
		if !seen {
			changed = true
		}
	}
	if !changed {
		return nil
	}
	p.keywords = vlist
	p.iim[idKeyword] = make([]iim.DataSet, len(vlist))
	for i := range vlist {
		p.iim[idKeyword][i] = iim.DataSet{ID: idKeyword, Data: []byte(vlist[i])}
	}
	p.setEncoding()
	p.dirty = true
	return nil
}
