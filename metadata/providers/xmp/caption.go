package xmp

import (
	"fmt"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	descriptionName      = rdf.Name{Namespace: nsDC, Name: "description"}
	userCommentName      = rdf.Name{Namespace: nsEXIF, Name: "UserComment"}
	imageDescriptionName = rdf.Name{Namespace: nsTIFF, Name: "ImageDescription"}
)

// getCaption reads the value of the Caption field from the RDF.
func (p *Provider) getCaption() (err error) {
	if p.dcDescription, err = getAlt(p.rdf.Property(descriptionName)); err != nil {
		return fmt.Errorf("dc:description: %s", err)
	}
	if p.exifUserComment, err = getAlt(p.rdf.Property(userCommentName)); err != nil {
		return fmt.Errorf("exif:UserComment: %s", err)
	}
	if p.tiffImageDescription, err = getAlt(p.rdf.Property(imageDescriptionName)); err != nil {
		return fmt.Errorf("tiff:ImageDescription: %s", err)
	}
	return nil
}

// Caption returns the value of the Caption field.
func (p *Provider) Caption() (value string) {
	if value = p.dcDescription.Default(); value != "" {
		return value
	}
	if value = p.exifUserComment.Default(); value != "" {
		return value
	}
	return p.tiffImageDescription.Default()
}

// CaptionTags returns a list of tag names for the Caption field, and a
// parallel list of values held by those tags.
func (p *Provider) CaptionTags() (tags []string, values [][]string) {
	tags = append(tags, "XML  dc:description")
	if len(p.dcDescription) == 0 {
		values = append(values, nil)
	} else {
		vlist := make([]string, len(p.dcDescription))
		for i := range p.dcDescription {
			vlist[i] = p.dcDescription[i].Value
		}
		values = append(values, vlist)
	}
	tags = append(tags, "XML  tiff:ImageDescription")
	if len(p.tiffImageDescription) == 0 {
		values = append(values, nil)
	} else {
		vlist := make([]string, len(p.tiffImageDescription))
		for i := range p.tiffImageDescription {
			vlist[i] = p.tiffImageDescription[i].Value
		}
		values = append(values, vlist)
	}
	if len(p.exifUserComment) != 0 {
		vlist := make([]string, len(p.exifUserComment))
		for i, as := range p.exifUserComment {
			vlist[i] = as.Value
		}
		tags = append(tags, "XML  exif:UserComment")
		values = append(values, vlist)
	}
	return tags, values
}

// SetCaption sets the value of the Caption field.
func (p *Provider) SetCaption(value string) error {
	p.exifUserComment = nil
	p.rdf.RemoveProperty(userCommentName)
	if value == "" {
		p.dcDescription = nil
		p.rdf.RemoveProperty(descriptionName)
		p.tiffImageDescription = nil
		p.rdf.RemoveProperty(imageDescriptionName)
		return nil
	}
	if len(p.dcDescription) != 1 || value != p.dcDescription[0].Value {
		p.dcDescription = newAltString(value)
		p.rdf.SetProperty(descriptionName, makeAlt(p.dcDescription))
	}
	if len(p.tiffImageDescription) != 1 || value != p.tiffImageDescription[0].Value {
		p.tiffImageDescription = newAltString(value)
		p.rdf.SetProperty(imageDescriptionName, makeAlt(p.tiffImageDescription))
	}
	return nil
}
