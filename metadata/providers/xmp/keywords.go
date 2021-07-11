package xmp

import (
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	tagsListName            = rdf.Name{Namespace: nsDigiKam, Name: "TagsList"}
	hierarchicalSubjectName = rdf.Name{Namespace: nsLR, Name: "hierarchicalSubject"}
	subjectName             = rdf.Name{Namespace: nsDC, Name: "subject"}
)

// getKeywords reads the values of the Groups, Keywords, People, Places, and
// Topics fields from the RDF.
func (p *Provider) getKeywords() (err error) {
	var list []string

	if list, err = getStrings(p.rdf.Properties, tagsListName); err != nil {
		return fmt.Errorf("digiKam:TagsList: %s", err)
	}
	for _, xkw := range list {
		p.digiKamTagsList = append(p.digiKamTagsList, strings.Split(xkw, "/"))
	}
	if list, err = getStrings(p.rdf.Properties, hierarchicalSubjectName); err != nil {
		return fmt.Errorf("lr:hierarchicalSubject: %s", err)
	}
	for _, xkw := range list {
		p.lrHierarchicalSubject = append(p.lrHierarchicalSubject, strings.Split(xkw, "|"))
	}
	if p.dcSubject, err = getStrings(p.rdf.Properties, subjectName); err != nil {
		return fmt.Errorf("dc:subject: %s", err)
	}
	return nil
}

// Keywords returns the values of the Keywords field.
func (p *Provider) Keywords() (values []metadata.HierValue) {
	return p.filteredKeywords(otherKeywordPredicate)
}

// KeywordsTags returns a list of tag names for the Keywords field, and
// a parallel list of values held by those tags.
func (p *Provider) KeywordsTags() (tags []string, values [][]metadata.HierValue) {
	return p.filteredKeywordsTags(otherKeywordPredicate)
}

// SetKeywords sets the values of the Keywords field.
func (p *Provider) SetKeywords(values []metadata.HierValue) error {
	p.setFilteredKeywords(otherKeywordPredicate, values)
	return nil
}

// otherKeywordPredicate is the predicate satisfied by keyword tags that encode
// keyword names.
func otherKeywordPredicate(kw metadata.HierValue) bool {
	return !groupPredicate(kw) && !personPredicate(kw) && !placePredicate(kw) && !topicPredicate(kw)
}

// A keywordFilter is a predicate function used to filter keywords.
type keywordFilter func(metadata.HierValue) bool

// allKeywordsFilter is a keywordFilter that matches everything.
func allKeywordsFilter(_ metadata.HierValue) bool { return true }

func (p *Provider) filteredKeywords(pred keywordFilter) (values []metadata.HierValue) {
	var leaves = make(map[string]bool)

	if len(p.digiKamTagsList) != 0 {
		for _, kw := range p.digiKamTagsList {
			if pred(kw) {
				values = append(values, kw)
			}
			if len(kw) != 0 {
				leaves[kw[len(kw)-1]] = true
			}
		}
		return values
	}
	if len(p.lrHierarchicalSubject) != 0 {
		for _, kw := range p.lrHierarchicalSubject {
			if pred(kw) {
				values = append(values, kw)
			}
			if len(kw) != 0 {
				leaves[kw[len(kw)-1]] = true
			}
		}
		return values
	}
	for _, s := range p.dcSubject {
		if leaves[s] {
			continue
		}
		kw := metadata.HierValue{s}
		if pred(kw) {
			values = append(values, kw)
		}
	}
	return values
}

func (p *Provider) filteredKeywordsTags(pred keywordFilter) (tags []string, values [][]metadata.HierValue) {
	var (
		tl, hs, s []metadata.HierValue
		leaves    = make(map[string]bool)
	)
	for _, hv := range p.digiKamTagsList {
		if pred(hv) {
			tl = append(tl, hv)
		}
		if len(hv) != 0 {
			leaves[hv[len(hv)-1]] = true
		}
	}
	for _, hv := range p.lrHierarchicalSubject {
		if pred(hv) {
			hs = append(hs, hv)
		}
		if len(hv) != 0 {
			leaves[hv[len(hv)-1]] = true
		}
	}
	tags = []string{"XMP  digiKam:TagsList", "XMP  lr:hierarchicalSubject"}
	values = [][]metadata.HierValue{tl, hs}
	for i := range p.dcSubject {
		if leaves[p.dcSubject[i]] {
			continue
		}
		var kw = metadata.HierValue{p.dcSubject[i]}
		if pred(kw) {
			s = append(s, kw)
		}
	}
	if len(s) != 0 {
		tags = append(tags, "XMP  dc:subject")
		values = append(values, s)
	}
	return tags, values
}

func (p *Provider) setFilteredKeywords(pred keywordFilter, values []metadata.HierValue) {
	p.setFilteredKeywordsHier(pred, values, &p.digiKamTagsList, tagsListName, "/")
	p.setFilteredKeywordsHier(pred, values, &p.lrHierarchicalSubject, hierarchicalSubjectName, "|")
	p.setFilteredKeywordsFlat(pred, values, &p.dcSubject, subjectName)
}
func (p *Provider) setFilteredKeywordsHier(
	pred keywordFilter, values []metadata.HierValue, plist *[]metadata.HierValue, name rdf.Name, sep string,
) {
	var kws []metadata.HierValue

	kws = append(kws, values...)
	for _, kw := range *plist {
		if !pred(kw) {
			kws = append(kws, kw)
		}
	}
	if len(kws) == 0 {
		*plist = nil
		if _, ok := p.rdf.Properties[name]; ok {
			delete(p.rdf.Properties, name)
			p.dirty = true
		}
		return
	}
	var kwstrs = make([]string, len(kws))
	var kwmap = make(map[string]bool)
	for _, kw := range kws {
		kwstr := strings.Join(kw, sep)
		kwstrs = append(kwstrs, kwstr)
		kwmap[kwstr] = false
	}
	var changed = false
	for _, kw := range *plist {
		kwstr := strings.Join(kw, sep)
		if _, ok := kwmap[kwstr]; !ok {
			changed = true
		} else {
			kwmap[kwstr] = true
		}
	}
	for _, seen := range kwmap {
		if !seen {
			changed = true
		}
	}
	if !changed {
		return
	}
	setBag(p.rdf.Properties, name, kwstrs)
	*plist = kws
	p.dirty = true
}
func (p *Provider) setFilteredKeywordsFlat(
	pred keywordFilter, values []metadata.HierValue, plist *[]string, name rdf.Name,
) {
	var vmap = make(map[string]bool)
	for _, kw := range values {
		vmap[kw[len(kw)-1]] = false
	}
	for _, kw := range *plist {
		if !pred(metadata.HierValue{kw}) {
			vmap[kw] = false
		}
	}
	if len(vmap) == 0 {
		*plist = nil
		if _, ok := p.rdf.Properties[name]; ok {
			delete(p.rdf.Properties, name)
			p.dirty = true
		}
		return
	}
	var changed = false
	for _, kw := range *plist {
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
		return
	}
	var vlist = make([]string, 0, len(vmap))
	for v := range vmap {
		vlist = append(vlist, v)
	}
	setBag(p.rdf.Properties, name, vlist)
	p.dcSubject = vlist
	p.dirty = true
}
