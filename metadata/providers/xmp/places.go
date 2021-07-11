package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// getPlaces reads the value of the Places field from the RDF.
func (p *Provider) getPlaces() error {
	return nil // handled by getKeywords
}

// Places returns the values of the Places field.
func (p *Provider) Places() (values []metadata.HierValue) {
	values = p.filteredKeywords(placePredicate)
	for i := range values {
		values[i] = values[i][1:]
	}
	return values
}

// PlacesTags returns a list of tag names for the Places field, and a
// parallel list of values held by those tags.
func (p *Provider) PlacesTags() (tags []string, values [][]metadata.HierValue) {
	tags, values = p.filteredKeywordsTags(placePredicate)
	for i := range tags {
		tags[i] += ":Places/"
		for j := range values[i] {
			values[i][j] = append(metadata.HierValue{}, values[i][j][1:]...)
		}
	}
	return tags, values
}

// SetPlaces sets the values of the Places field.
func (p *Provider) SetPlaces(values []metadata.HierValue) error {
	var kws = make([]metadata.HierValue, len(values))
	for i := range values {
		kws[i] = append(metadata.HierValue{"Places"}, values[i]...)
	}
	p.setFilteredKeywords(placePredicate, kws)
	return nil
}

// placePredicate is the predicate satisfied by keyword tags that encode place
// names.
func placePredicate(kw metadata.HierValue) bool {
	return len(kw) >= 2 && kw[0] == "Places"
}
