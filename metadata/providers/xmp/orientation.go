package xmp

import (
	"fmt"
	"strconv"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	orientationName = rdf.Name{Namespace: nsTIFF, Name: "Orientation"}
)

// getOrientation reads the value of the Orientation field from the RDF.
func (p *Provider) getOrientation() (err error) {
	var ostring string
	var oval int

	if ostring, err = getString(p.rdf.Property(orientationName)); err != nil {
		return fmt.Errorf("tiff:Orientation: %s", err)
	}
	if oval, err = strconv.Atoi(ostring); err != nil {
		return fmt.Errorf("tiff:Orientation: %s", err)
	}
	p.tiffOrientation = metadata.Orientation(oval)
	return nil
}

// Orientation returns the value of the Orientation field.
func (p *Provider) Orientation() (value metadata.Orientation) {
	return p.tiffOrientation
}

// OrientationTags returns a list of tag names for the Orientation field, and a
// parallel list of values held by those tags.
func (p *Provider) OrientationTags() (tags []string, values [][]metadata.Orientation) {
	if p.tiffOrientation == 0 || p.tiffOrientation == metadata.Rotate0 {
		values = [][]metadata.Orientation{nil}
	} else {
		values = [][]metadata.Orientation{{p.tiffOrientation}}
	}
	return []string{"XMP  tiff:Orientation"}, values
}

// SetOrientation sets the value of the Orientation field.
func (p *Provider) SetOrientation(value metadata.Orientation) error {
	if value == 0 || value == metadata.Rotate0 {
		p.tiffOrientation = 0
		p.rdf.RemoveProperty(orientationName)
		return nil
	}
	p.tiffOrientation = value
	p.rdf.SetProperty(orientationName, makeString(strconv.Itoa(int(p.tiffOrientation))))
	return nil
}
