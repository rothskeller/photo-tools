package xmpexif

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var userCommentName = rdf.Name{Namespace: nsEXIF, Name: "UserComment"}

// getCaption reads the value of the Caption field from the RDF.
func (p *Provider) getCaption() (err error) {
	if p.exifUserComment, err = getAlt(p.rdf.Properties, userCommentName); err != nil {
		return fmt.Errorf("exif:UserComment: %s", err)
	}
	return nil
}

// Caption returns the value of the Caption field.
func (p *Provider) Caption() (value string) { return p.exifUserComment.Default() }

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p *Provider) CaptionTags() (tags []string, values [][]string) {
	if len(p.exifUserComment) == 0 {
		return nil, nil
	}
	values = [][]string{nil}
	for _, as := range p.exifUserComment {
		values[0] = append(values[0], as.Value)
	}
	return []string{"XML  exif:UserComment"}, values
}

// SetCaption sets the value of the Caption field.
func (p *Provider) SetCaption(value string) error {
	p.exifUserComment = nil
	if _, ok := p.rdf.Properties[userCommentName]; ok {
		delete(p.rdf.Properties, userCommentName)
		p.dirty = true
	}
	return nil
}
