package strmeta

import (
	"errors"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

// A Topic represents a topic of a media file.
type Topic []string

// Parse parses a topic name, as a hierarchical string with levels separated by
// slashes and optional whitespace.  Pipe symbols are disallowed due to
// underlying storage formats, and empty levels are disallowed (although a
// completely empty string is allowed).
func (g *Topic) Parse(s string) error {
	kw, err := metadata.ParseKeyword(s, "")
	if err == nil {
		*g = Topic(kw)
	}
	return err
}

// String returns the formatted string form of the topic name, suitable for
// input to Parse.
func (g Topic) String() string { return metadata.Keyword(g).String() }

// Empty returns whether the topic name is empty.
func (g Topic) Empty() bool { return len(g) == 0 }

// Equal returns whether two topic names are equal.
func (g Topic) Equal(other Topic) bool {
	return metadata.Keyword(g).Equal(metadata.Keyword(other))
}

// GetTopics returns the highest priority topic values.
func GetTopics(h filefmt.FileHandler) []Topic {
	kws := getFilteredKeywords(h, topicPredicate, false)
	topics := make([]Topic, len(kws))
	for i := range kws {
		topics[i] = Topic(kws[i][1:])
	}
	return topics
}

// GetTopicTags returns all of the topic tags and their values.
func GetTopicTags(h filefmt.FileHandler) (tags []string, values []Topic) {
	tags, kws := getFilteredKeywordTags(h, topicPredicate)
	values = make([]Topic, len(kws))
	for i := range kws {
		values[i] = Topic(kws[i][1:])
	}
	return tags, values
}

// CheckTopics determines whether the topics are tagged correctly.
func CheckTopics(h filefmt.FileHandler) CheckResult {
	return checkFilteredKeywords(h, topicPredicate)
}

// SetTopics sets the topic tags.
func SetTopics(h filefmt.FileHandler, v []Topic) error {
	var kws = make([]metadata.Keyword, len(v))
	for i, g := range v {
		if g.Empty() {
			return errors.New("empty topic name not allowed")
		}
		kws[i] = append(metadata.Keyword{"Topics"}, v[i]...)
	}
	return setFilteredKeywords(h, kws, topicPredicate)
}

// topicPredicate is the predicate satisfied by keyword tags that encode topic
// names.
func topicPredicate(kw metadata.Keyword) bool {
	return len(kw) >= 2 && kw[0] == "Topics"
}
