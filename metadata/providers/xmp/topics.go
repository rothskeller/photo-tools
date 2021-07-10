package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// getTopics reads the value of the Topics field from the RDF.
func (p *Provider) getTopics() error {
	return nil // handled by getKeywords
}

// Topics returns the values of the Topics field.
func (p *Provider) Topics() (values []metadata.HierValue) {
	values = p.filteredKeywords(topicPredicate)
	for i := range values {
		values[i] = values[i][1:]
	}
	return values
}

// TopicsTags returns a list of tag names for the Topics field, and a
// parallel list of values held by those tags.
func (p *Provider) TopicsTags() (tags []string, values []metadata.HierValue) {
	tags, values = p.filteredKeywordsTags(topicPredicate)
	for i := range tags {
		tags[i] += ":Topics/"
		values[i] = values[i][1:]
	}
	return tags, values
}

// SetTopics sets the values of the Topics field.
func (p *Provider) SetTopics(values []metadata.HierValue) error {
	var kws = make([]metadata.HierValue, len(values))
	for i := range values {
		kws[i] = append(metadata.HierValue{"Topics"}, values[i]...)
	}
	p.setFilteredKeywords(topicPredicate, kws)
	return nil
}

// topicPredicate is the predicate satisfied by keyword tags that encode topic
// names.
func topicPredicate(kw metadata.HierValue) bool {
	return len(kw) >= 2 && kw[0] == "Topics"
}
