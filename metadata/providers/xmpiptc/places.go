package xmpiptc

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// Places returns no value.  The mirrored IPTC data does not contain places.
func (p *Provider) Places() (values []metadata.HierValue) { return nil }

// PlacesTags returns no tags.  The mirrored IPTC data does not contain places.
func (p *Provider) PlacesTags() (tags []string, values [][]metadata.HierValue) { return nil, nil }

// SetPlaces clears the value of the Location field if the places that are being
// set do not include the location.
func (p *Provider) SetPlaces(values []metadata.HierValue) (err error) {
	var location = p.Location()
	var parts []string
	if location.CountryName != "" {
		parts = append(parts, location.CountryName)
	}
	if location.State != "" {
		parts = append(parts, location.State)
	}
	if location.City != "" {
		parts = append(parts, location.City)
	}
	if location.Sublocation != "" {
		parts = append(parts, location.Sublocation)
	}
	if len(parts) == 0 {
		return metadata.ErrNotSupported
	}
	for _, p := range values {
		if locationCongruentToPlace(parts, p) {
			return metadata.ErrNotSupported
		}
	}
	if err = p.SetLocation(metadata.Location{}); err != nil {
		return err
	}
	return metadata.ErrNotSupported
}
func locationCongruentToPlace(loc []string, place metadata.HierValue) bool {
	for len(loc) != 0 && len(place) != 0 {
		if loc[0] == place[0] {
			loc = loc[1:]
		}
		place = place[1:]
	}
	return len(loc) == 0
}
