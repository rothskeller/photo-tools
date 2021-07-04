package fields

import (
	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/strmeta"
)

type topicsField struct {
	baseField
}

// TopicsField is the field handler for topic names, i.e., topics of people
// (organizations, teams, etc.) that are depicted in the media.
var TopicsField Field = &topicsField{
	baseField{
		name:        "topic",
		pluralName:  "topics",
		label:       "Topic",
		shortLabel:  "TP",
		multivalued: true,
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *topicsField) ParseValue(s string) (interface{}, error) {
	var v strmeta.Topic
	if err := v.Parse(s); err != nil {
		return nil, err
	}
	return v, nil
}

// RenderValue takes a value for the field and renders it in string form for
// display.
func (f *topicsField) RenderValue(v interface{}) string { return v.(strmeta.Topic).String() }

// EqualValue compares two values for equality.
func (f *topicsField) EqualValue(a interface{}, b interface{}) bool {
	return a.(strmeta.Topic).Equal(b.(strmeta.Topic))
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *topicsField) GetValues(h filefmt.FileHandler) []interface{} {
	var topics = strmeta.GetTopics(h)
	var ifcs = make([]interface{}, len(topics))
	for i := range topics {
		ifcs[i] = topics[i]
	}
	return ifcs
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *topicsField) GetTags(h filefmt.FileHandler) ([]string, []interface{}) {
	var tags, values = strmeta.GetTopicTags(h)
	var ifcs = make([]interface{}, len(values))
	for i := range values {
		ifcs[i] = values[i]
	}
	return tags, ifcs
}

// CheckValues returns whether the values of the field are tagged correctly.
func (f *topicsField) CheckValues(h filefmt.FileHandler) (res strmeta.CheckResult) {
	if res = strmeta.CheckTopics(h); res == strmeta.ChkPresent {
		res = strmeta.CheckResult(len(f.GetValues(h)))
	}
	return res
}

// SetValues sets all of the values of the field.
func (f *topicsField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	var topics = make([]strmeta.Topic, len(v))
	for i := range v {
		topics[i] = v[i].(strmeta.Topic)
	}
	return strmeta.SetTopics(h, topics)
}
