package fields

import (
	"errors"
	"strings"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/strmeta"
)

type peopleField struct {
	baseKWField
}

// A personInfo represents a single person; it's the value type of the Person
// field.  It contains the person's name and a flag indicating whether there is
// a face region for that name.  The flag is only used for display.
type personInfo struct {
	name       string
	faceRegion bool
}

// PeopleField is the field handler for people, which come both from keywords
// that start with People and from face region tags.  These represent people who
// are depicted in the media.
var PeopleField Field = &peopleField{
	baseKWField{
		baseField: baseField{
			name:        "person",
			pluralName:  "people",
			label:       "Person",
			shortLabel:  "PE",
			multivalued: true,
		},
		prefix: "People",
	},
}

// ParseValue parses a string and returns a value for the field.  It returns an
// error if the string is invalid.
func (f *peopleField) ParseValue(s string) (interface{}, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, " [F]") {
		// This shouldn't happen, but we might as well handle it.
		s = s[:len(s)-4]
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("invalid person name (must not be empty)")
	}
	if strings.IndexByte(s, '/') >= 0 || strings.IndexByte(s, '|') >= 0 {
		return nil, errors.New("invalid person name (cannot contain / or |)")
	}
	return &personInfo{name: s}, nil
}

// RenderValue takes a value for the field and renders it in string form
// for display.
func (f *peopleField) RenderValue(v interface{}) string {
	p := v.(*personInfo)
	if p.faceRegion { // Mark people who have face regions with [F]
		return p.name + " [F]"
	}
	return p.name
}

// EqualValue compares two values for equality.
func (f *peopleField) EqualValue(a interface{}, b interface{}) bool {
	return a.(*personInfo).name == b.(*personInfo).name
	// We don't care if the faceRegion flags match.
}

// GetValues returns all of the values of the field.  (For single-valued fields,
// the return slice will have at most one entry.)  Empty values should not be
// included.
func (f *peopleField) GetValues(h filefmt.FileHandler) (people []interface{}) {
	// Get the People keywords and the face regions.
	kws := f.baseKWField.GetValues(h)
	faces := strmeta.GetFaces(h)
	for _, kwi := range kws {
		kw := kwi.(metadata.Keyword)
		if len(kw) != 2 { // ignore People keywords with multiple components.
			continue
		}
		faceRegion := false // note whether the name has a face region
		for _, f := range faces {
			if f == kw[1] {
				faceRegion = true
				break
			}
		}
		people = append(people, &personInfo{kw[1], faceRegion})
	}
	return people
}

// GetTags returns the names of all of the metadata tags that correspond to the
// field in its first return slice, and a parallel slice of the values of those
// tags (which may be zero values).
func (f *peopleField) GetTags(h filefmt.FileHandler) (tags []string, values []interface{}) {
	// People tags are reported by the Keywords field.  Here I only need to
	// report the face region tags.
	var ftags, fvals = strmeta.GetFacesTags(h)
	tags = ftags
	values = make([]interface{}, len(fvals))
	for i := range fvals {
		values[i] = &personInfo{name: fvals[i]}
		// We're not setting the faceRegion flag, because that's obvious
		// from the tag name.
	}
	return tags, values
}

// SetValues sets all of the values of the field.
func (f *peopleField) SetValues(h filefmt.FileHandler, v []interface{}) error {
	// First, "set" the face region list.  We can't actually change it, but
	// we'll fail on the attempt.
	var faces []string
	for _, p := range v {
		person := p.(*personInfo)
		if person.faceRegion {
			faces = append(faces, person.name)
		}
	}
	if err := strmeta.SetFaces(h, faces); err != nil { // fails if the list changes
		return err
	}
	// Now set the People keywords.
	ptags := make([]interface{}, len(v))
	for i := range v {
		ptags[i] = metadata.Keyword{f.prefix, v[i].(*personInfo).name}
	}
	return f.baseKWField.SetValues(h, ptags)
}

// CheckValues returns whether the values of the field in the target are tagged
// correctly, and are consistent with the values of the field in the reference.
func (f *peopleField) CheckValues(ref filefmt.FileHandler, tgt filefmt.FileHandler) strmeta.CheckResult {
	// First make sure that the People keywords are valid.
	if res := f.baseKWField.CheckValues(ref, tgt); res < 0 {
		return res
	}
	// Now get the reference values and check each of the face region tags
	// for consistency with it.  Basically we're looking for face regions
	// that don't have a corresponding People keyword, or inconsistencies
	// between the two face region tag sets.
	values := f.GetValues(ref)
	if xmp := tgt.XMP(false); xmp != nil {
		if res := checkFaceTags(values, xmp.MPFaces); res < 0 {
			return res
		}
		if res := checkFaceTags(values, xmp.MWGRSFaces); res < 0 {
			return res
		}
	}
	if len(values) != 0 {
		return strmeta.ChkPresent
	}
	return strmeta.ChkOptionalAbsent
}
func checkFaceTags(ref []interface{}, faces []string) (result strmeta.CheckResult) {
	// First, check for people who have face regions, and make sure they are
	// in the faces list.
	for _, pi := range ref {
		p := pi.(*personInfo)
		if !p.faceRegion {
			continue
		}
		found := false
		for _, f := range faces {
			if f == p.name {
				found = true
				break
			}
		}
		if !found {
			if len(faces) != 0 {
				return strmeta.ChkConflictingValues
			}
			result = strmeta.ChkIncorrectlyTagged
		}
	}
	// Next, look through the faces list, and make sure there are People
	// keywords for each.
	for _, f := range faces {
		found := false
		for _, pi := range ref {
			p := pi.(*personInfo)
			if p.faceRegion && p.name == f {
				found = true
				break
			}
		}
		if !found {
			return strmeta.ChkConflictingValues
		}
	}
	return result
}
