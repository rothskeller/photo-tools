package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
)

// GetKeywords returns keywords starting with the specified prefix, coming from
// the highest priority keyword tag.
func GetKeywords(h fileHandler, prefix string) (matches []metadata.Keyword) {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DigiKamTagsList) != 0 {
			if prefix == "" {
				return xmp.DigiKamTagsList
			}
			for _, kw := range xmp.DigiKamTagsList {
				if kw[0].Word == prefix {
					matches = append(matches, kw[1:])
				}
			}
			return matches
		}
		if len(xmp.LRHierarchicalSubject) != 0 {
			if prefix == "" {
				return xmp.LRHierarchicalSubject
			}
			for _, kw := range xmp.LRHierarchicalSubject {
				if kw[0].Word == prefix {
					matches = append(matches, kw[1:])
				}
			}
			return matches
		}
		if prefix == "" && len(xmp.DCSubject) != 0 {
			matches = make([]metadata.Keyword, len(xmp.DCSubject))
			for i, s := range xmp.DCSubject {
				matches[i] = []metadata.KeywordComponent{{Word: s}}
			}
			return matches
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if prefix == "" && len(iptc.Keywords) != 0 {
			matches = make([]metadata.Keyword, len(iptc.Keywords))
			for i, s := range iptc.Keywords {
				matches[i] = []metadata.KeywordComponent{{Word: s}}
			}
			return matches
		}
	}
	return nil
}

// GetKeywordsTags returns all of the keyword tags for keywords with the
// specified prefix, and their values.
func GetKeywordsTags(h fileHandler, prefix string) (tags []string, values []metadata.Keyword) {
	if xmp := h.XMP(false); xmp != nil {
		for _, kw := range xmp.DigiKamTagsList {
			if prefix == "" || kw[0].Word == prefix {
				tags = append(tags, "XMP.digiKam:TagsList")
				values = append(values, kw)
			}
		}
		for _, kw := range xmp.LRHierarchicalSubject {
			if prefix == "" || kw[0].Word == prefix {
				tags = append(tags, "XMP.lr:HierarchicalSubject")
				values = append(values, kw)
			}
		}
		for _, kw := range xmp.DCSubject {
			if prefix == "" {
				tags = append(tags, "XMP.dc:Subject")
				values = append(values, []metadata.KeywordComponent{{Word: kw}})
			}
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		for _, kw := range iptc.Keywords {
			if prefix == "" {
				tags = append(tags, "IPTC.Keyword")
				values = append(values, []metadata.KeywordComponent{{Word: kw}})
			}
		}
	}
	return tags, values
}

// SetKeywords sets the creator tags with the specified prefix.  It leaves
// keywords with other prefixes alone.
func SetKeywords(h fileHandler, prefix string, v []metadata.Keyword) error {
	var list []metadata.Keyword
	var flat []string

	if prefix != "" {
		if xmp := h.XMP(false); xmp != nil {
			var source []metadata.Keyword
			if len(xmp.DigiKamTagsList) != 0 {
				source = xmp.DigiKamTagsList
			} else if len(xmp.LRHierarchicalSubject) != 0 {
				source = xmp.LRHierarchicalSubject
			}
			for _, kw := range source {
				if kw[0].Word != prefix {
					list = append(list, kw)
				}
			}
		}
		for _, kw := range v {
			kw = append([]metadata.KeywordComponent{{Word: prefix, OmitWhenFlattened: true}}, kw...)
			list = append(list, kw)
		}
	} else {
		list = v
	}
	for _, kw := range list {
		for _, comp := range kw {
			if !comp.OmitWhenFlattened {
				var found = false
				for _, f := range flat {
					if f == comp.Word {
						found = true
						break
					}
				}
				if !found {
					flat = append(flat, comp.Word)
				}
			}
		}
	}
	if xmp := h.XMP(true); xmp != nil {
		xmp.DigiKamTagsList = list
		xmp.LRHierarchicalSubject = list
		xmp.DCSubject = flat
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.Keywords = flat
	}
	return nil
}
