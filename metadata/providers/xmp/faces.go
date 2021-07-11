package xmp

import (
	"errors"
	"sort"

	"github.com/rothskeller/photo-tools/metadata/containers/rdf"
)

var (
	mpRegionInfoName           = rdf.Name{Namespace: nsMP, Name: "RegionInfo"}
	mpriRegionsName            = rdf.Name{Namespace: nsMPRI, Name: "Regions"}
	mpRegRectangleName         = rdf.Name{Namespace: nsMPReg, Name: "Rectangle"}
	mpRegPersonDisplayNameName = rdf.Name{Namespace: nsMPReg, Name: "PersonDisplayName"}
	mwgrsRegionsName           = rdf.Name{Namespace: nsMWGRS, Name: "Regions"}
	mwgrsRegionListName        = rdf.Name{Namespace: nsMWGRS, Name: "RegionList"}
	mwgrsTypeName              = rdf.Name{Namespace: nsMWGRS, Name: "Type"}
	mwgrsNameName              = rdf.Name{Namespace: nsMWGRS, Name: "Name"}
)

// getFaces reads the value of the Faces field from the RDF.
func (p *Provider) getFaces() (err error) {
	if err = p.getFacesMP(); err != nil {
		return err
	}
	return p.getFacesMWGRS()
}
func (p *Provider) getFacesMP() (err error) {
	regionInfo, ok := p.rdf.Properties[mpRegionInfoName]
	if !ok {
		return nil
	}
	regionInfoStruct, ok := regionInfo.Value.(rdf.Struct)
	if !ok {
		return errors.New("MP:RegionInfo: wrong data type")
	}
	regions, ok := regionInfoStruct[mpriRegionsName]
	if !ok {
		return nil
	}
	bag, ok := regions.Value.(rdf.Bag)
	if !ok {
		return errors.New("MPRI:Regions: wrong data type")
	}
	p.mpRegPersonDisplayNames = make([]string, 0, len(bag))
	for _, reg := range bag {
		region, ok := reg.Value.(rdf.Struct)
		if !ok {
			return errors.New("MPRI:Regions element: wrong data type")
		}
		if _, ok := region[mpRegRectangleName]; !ok {
			// digiKam tends to create region entries for anyone who has a face
			// in *any* photo, but omits the rectangle if they don't have a face
			// in *this* photo.  For our purposes, if there's no rectangle, it
			// doesn't count.
			continue
		}
		personDisplayName, ok := region[mpRegPersonDisplayNameName]
		if !ok {
			continue
		}
		name, ok := personDisplayName.Value.(string)
		if !ok {
			return errors.New("MPReg:PersonDisplayName: wrong data type")
		}
		p.mpRegPersonDisplayNames = append(p.mpRegPersonDisplayNames, name)
	}
	return nil
}
func (p *Provider) getFacesMWGRS() (err error) {
	regions, ok := p.rdf.Properties[rdf.Name{Namespace: nsMWGRS, Name: "Regions"}]
	if !ok {
		return nil
	}
	regionsStruct, ok := regions.Value.(rdf.Struct)
	if !ok {
		return errors.New("mwg-rs:Regions: wrong data type")
	}
	regionList, ok := regionsStruct[rdf.Name{Namespace: nsMWGRS, Name: "RegionList"}]
	if !ok {
		return nil
	}
	var items []rdf.Value
	bag, ok := regionList.Value.(rdf.Bag)
	if ok {
		items = bag
	} else {
		// Some images incorrectly use an rdf:Seq here instead.
		seq, ok := regionList.Value.(rdf.Seq)
		if ok {
			items = seq
		} else {
			return errors.New("mwg-rs:RegionList: wrong data type")
		}
	}
	p.mwgrsNames = make([]string, 0, len(items))
	for _, reg := range items {
		region, ok := reg.Value.(rdf.Struct)
		if !ok {
			return errors.New("mwg-rs:RegionList element: wrong data type")
		}
		typ, ok := region[rdf.Name{Namespace: nsMWGRS, Name: "Type"}]
		if !ok {
			continue
		}
		typString, ok := typ.Value.(string)
		if !ok {
			return errors.New("mwg-rs:Type: wrong data type")
		}
		if typString != "Face" {
			continue
		}
		name, ok := region[rdf.Name{Namespace: nsMWGRS, Name: "Name"}]
		if !ok {
			continue
		}
		nameString, ok := name.Value.(string)
		if !ok {
			return errors.New("mwg-rs:Name: wrong data type")
		}
		p.mwgrsNames = append(p.mwgrsNames, nameString)
	}
	return nil
}

// Faces returns the values of the Faces field.
func (p *Provider) Faces() (values []string) {
	if len(p.mpRegPersonDisplayNames) == 0 && len(p.mwgrsNames) == 0 {
		return nil
	}
	var facemap = make(map[string]bool)
	for _, face := range p.mpRegPersonDisplayNames {
		facemap[face] = true
	}
	for _, face := range p.mwgrsNames {
		facemap[face] = true
	}
	values = make([]string, 0, len(facemap))
	for face := range facemap {
		values = append(values, face)
	}
	sort.Strings(values)
	return values
}

// FacesTags returns a list of tag names for the Faces field, and a
// parallel list of values held by those tags.
func (p *Provider) FacesTags() (tags []string, values [][]string) {
	if len(p.mpRegPersonDisplayNames) != 0 {
		tags = append(tags, "XMP  MP:Regions")
		values = append(values, p.mpRegPersonDisplayNames)
	}
	if len(p.mwgrsNames) != 0 {
		tags = append(tags, "XMP  mwg-rs:RegionInfo")
		values = append(values, p.mwgrsNames)
	}
	return tags, values
}

// SetFaces sets the values of the Faces field.
func (p *Provider) SetFaces(values []string) error {
	var (
		regions   rdf.Struct
		bag       rdf.Bag
		nextInBag int
		vmap      = make(map[string]bool)
	)
	for _, value := range values {
		vmap[value] = false
	}
	if val, ok := p.rdf.Properties[mpRegionInfoName]; ok {
		regions = val.Value.(rdf.Struct)
		if val, ok = regions[mpriRegionsName]; ok {
			bag = val.Value.(rdf.Bag)
		}
	}
	p.mpRegPersonDisplayNames = nil
	for _, oldv := range bag {
		if _, ok := oldv.Value.(rdf.Struct)[mpRegRectangleName]; !ok {
			bag[nextInBag] = oldv
			nextInBag++
		} else if n, ok := oldv.Value.(rdf.Struct)[mpRegPersonDisplayNameName]; !ok {
			bag[nextInBag] = oldv
			nextInBag++
		} else if _, ok := vmap[n.Value.(string)]; ok {
			bag[nextInBag] = oldv
			vmap[n.Value.(string)] = true
			p.mpRegPersonDisplayNames = append(p.mpRegPersonDisplayNames, n.Value.(string))
		}
	}
	if nextInBag != len(bag) {
		regions[mpriRegionsName] = rdf.Value{Value: bag[:nextInBag]}
		p.dirty = true
	}
	bag, nextInBag = nil, 0
	if val, ok := p.rdf.Properties[mwgrsRegionsName]; ok {
		regions = val.Value.(rdf.Struct)
		if val, ok = regions[mwgrsRegionListName]; ok {
			bag = val.Value.(rdf.Bag)
		}
	}
	p.mwgrsNames = nil
	for _, oldv := range bag {
		if t, ok := oldv.Value.(rdf.Struct)[mwgrsTypeName]; !ok || t.Value.(string) != "Face" {
			bag[nextInBag] = oldv
			nextInBag++
		} else if n, ok := oldv.Value.(rdf.Struct)[mwgrsNameName]; !ok {
			bag[nextInBag] = oldv
			nextInBag++
		} else if _, ok := vmap[n.Value.(string)]; ok {
			bag[nextInBag] = oldv
			vmap[n.Value.(string)] = true
			p.mwgrsNames = append(p.mwgrsNames, n.Value.(string))
		}
	}
	if nextInBag != len(bag) {
		regions[mwgrsRegionListName] = rdf.Value{Value: bag[:nextInBag]}
		p.dirty = true
	}
	for _, seen := range vmap {
		if !seen {
			return errors.New("cannot add face regions")
		}
	}
	return nil
}
