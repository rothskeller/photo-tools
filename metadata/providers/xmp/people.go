package xmp

import "github.com/rothskeller/photo-tools/metadata"

// getPeople reads the value of the People field from the RDF.
func (p *Provider) getPeople() error {
	return nil // handled by getKeywords
}

// People returns the values of the People field.
func (p *Provider) People() (values []string) {
	pmap := make(map[string]bool)
	kws := p.filteredKeywords(personPredicate)
	values = make([]string, len(kws))
	for i := range kws {
		values[i] = kws[i][1]
		pmap[values[i]] = true
	}
	for _, face := range p.mpRegPersonDisplayNames {
		if !pmap[face] {
			values = append(values, face)
			pmap[face] = true
		}
	}
	for _, face := range p.mwgrsNames {
		if !pmap[face] {
			values = append(values, face)
			pmap[face] = true
		}
	}
	return values
}

// PeopleTags returns a list of tag names for the People field, and a
// parallel list of values held by those tags.
func (p *Provider) PeopleTags() (tags []string, values [][]string) {
	var kws [][]metadata.HierValue

	tags, kws = p.filteredKeywordsTags(personPredicate)
	values = make([][]string, len(kws))
	for i := range tags {
		tags[i] += ":People/"
		values[i] = make([]string, len(kws[i]))
		for j := range kws[i] {
			values[i][j] = kws[i][j][1]
		}
	}
	return tags, values
}

// SetPeople sets the values of the People field.
func (p *Provider) SetPeople(values []string) (err error) {
	var faces = p.Faces()
	var facemap = make(map[string]bool)
	for _, face := range faces {
		facemap[face] = false
	}
	faces = faces[:0]
	var kws = make([]metadata.HierValue, len(values))
	for i := range values {
		kws[i] = metadata.HierValue{"People", values[i]}
		if seen, ok := facemap[values[i]]; ok && !seen {
			facemap[values[i]] = true
			faces = append(faces, values[i])
		}
	}
	p.setFilteredKeywords(personPredicate, kws)
	return p.SetFaces(faces)
}

// personPredicate is the predicate satisfied by keyword tags that encode person
// names.
func personPredicate(kw metadata.HierValue) bool {
	return len(kw) == 2 && kw[0] == "People"
}
