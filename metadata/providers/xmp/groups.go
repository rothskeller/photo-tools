package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// getGroups reads the value of the Groups field from the RDF.
func (p *Provider) getGroups() error {
	return nil // handled by getKeywords
}

// Groups returns the values of the Groups field.
func (p *Provider) Groups() (values []metadata.HierValue) {
	values = p.filteredKeywords(groupPredicate)
	for i := range values {
		values[i] = values[i][1:]
	}
	return values
}

// GroupsTags returns a list of tag names for the Groups field, and a
// parallel list of values held by those tags.
func (p *Provider) GroupsTags() (tags []string, values [][]metadata.HierValue) {
	tags, values = p.filteredKeywordsTags(groupPredicate)
	for i := range tags {
		tags[i] += ":Groups/"
		for j := range values[i] {
			values[i][j] = append(metadata.HierValue{}, values[i][j][1:]...)
		}
	}
	return tags, values
}

// SetGroups sets the values of the Groups field.
func (p *Provider) SetGroups(values []metadata.HierValue) error {
	var kws = make([]metadata.HierValue, len(values))
	for i := range values {
		kws[i] = append(metadata.HierValue{"Groups"}, values[i]...)
	}
	p.setFilteredKeywords(groupPredicate, kws)
	return nil
}

// groupPredicate is the predicate satisfied by keyword tags that encode group
// names.
func groupPredicate(kw metadata.HierValue) bool {
	return len(kw) >= 2 && kw[0] == "Groups"
}
