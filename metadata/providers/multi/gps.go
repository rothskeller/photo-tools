package multi

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// GPS returns the value of the GPS field.
func (p Provider) GPS() (value metadata.GPSCoords) {
	for _, sp := range p {
		if value = sp.GPS(); !value.Empty() {
			return value
		}
	}
	return metadata.GPSCoords{}
}

// GPSTags returns a list of tag names for the GPS field, and a parallel
// list of values held by those tags.
func (p Provider) GPSTags() (tags []string, values []metadata.GPSCoords) {
	for _, sp := range p {
		t, v := sp.GPSTags()
		tags = append(tags, t...)
		values = append(values, v...)
	}
	return tags, values
}

// SetGPS sets the value of the GPS field.
func (p Provider) SetGPS(value metadata.GPSCoords) error {
	var set = false

	for _, sp := range p {
		if err := sp.SetGPS(value); err != nil && err != metadata.ErrNotSupported {
			return err
		} else if err == nil {
			set = true
		}
	}
	if !set {
		return metadata.ErrNotSupported
	}
	return nil
}
