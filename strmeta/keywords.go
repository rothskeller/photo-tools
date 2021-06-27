package strmeta

import (
	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/iptc"
)

// GetKeywords returns keywords from the highest priority keyword tag.
func GetKeywords(h fileHandler) (kws []metadata.Keyword) {
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DigiKamTagsList) != 0 {
			kws = xmp.DigiKamTagsList
		} else if len(xmp.LRHierarchicalSubject) != 0 {
			kws = xmp.LRHierarchicalSubject
		}
	}
	var flat = getFlatKeywords(h)
	if flat == nil {
		return kws
	}
	for _, kw := range kws {
		for i, word := range kw {
			if _, ok := flat[word.Word]; ok {
				flat[word.Word] = false
			} else {
				kw[i].OmitWhenFlattened = true
			}
		}
	}
	for word, unseen := range flat {
		if unseen {
			kws = append(kws, metadata.Keyword{{Word: word}})
		}
	}
	return kws
}

// getFlatKeywords returns the flat keywords from the highest priority flat
// keyword tag, in the form of a map with the keyword as key and true as the
// value.
func getFlatKeywords(h fileHandler) (kws map[string]bool) {
	kws = make(map[string]bool)
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DCSubject) != 0 {
			for _, kw := range xmp.DCSubject {
				kws[kw] = true
			}
			return kws
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if len(iptc.Keywords) != 0 {
			for _, kw := range iptc.Keywords {
				kws[kw] = true
			}
			return kws
		}
	}
	return nil
}

// GetKeywordsTags returns all of the keyword tags and their values.
func GetKeywordsTags(h fileHandler) (tags []string, values []metadata.Keyword) {
	if xmp := h.XMP(false); xmp != nil {
		for _, kw := range xmp.DigiKamTagsList {
			tags = append(tags, "XMP.digiKam:TagsList")
			values = append(values, kw)
		}
		for _, kw := range xmp.LRHierarchicalSubject {
			tags = append(tags, "XMP.lr:HierarchicalSubject")
			values = append(values, kw)
		}
		for _, kw := range xmp.DCSubject {
			tags = append(tags, "XMP.dc:Subject")
			values = append(values, []metadata.KeywordComponent{{Word: kw}})
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		for _, kw := range iptc.Keywords {
			tags = append(tags, "IPTC.Keyword")
			values = append(values, []metadata.KeywordComponent{{Word: kw}})
		}
	}
	return tags, values
}

// CheckKeywords checks whether the keywords are correctly tagged, and are
// consistent with the reference, and whether expected keywords with the
// specified prefix are present.
func CheckKeywords(ref, h fileHandler) (res CheckResult) {
	var (
		value   = GetKeywords(ref)
		flat    = map[string]bool{}
		flatmax = map[string]bool{}
	)
	for _, kw := range value {
		for _, c := range kw {
			if !c.OmitWhenFlattened {
				flat[c.Word] = true
				if len(c.Word) > iptc.MaxKeywordLen {
					flatmax[c.Word[:iptc.MaxKeywordLen]] = true
				} else {
					flatmax[c.Word] = true
				}
			}
		}
	}
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DigiKamTagsList) != 0 {
			if !keywordsEqual(value, xmp.DigiKamTagsList) {
				return ChkConflictingValues
			}
		} else if len(value) != 0 {
			res = ChkIncorrectlyTagged
		}
		if len(xmp.LRHierarchicalSubject) != 0 {
			if !keywordsEqual(value, xmp.LRHierarchicalSubject) {
				return ChkConflictingValues
			}
		} else if len(value) != 0 {
			res = ChkIncorrectlyTagged
		}
		if len(xmp.DCSubject) != 0 {
			if len(xmp.DCSubject) != len(flat) {
				return ChkConflictingValues
			}
			smap := make(map[string]bool)
			for _, s := range xmp.DCSubject {
				if !flat[s] {
					return ChkConflictingValues
				}
				smap[s] = true
			}
			for s := range flat {
				if !smap[s] {
					return ChkConflictingValues
				}
			}
		} else if len(flat) != 0 {
			res = ChkIncorrectlyTagged
		}
	}
	if i := h.IPTC(); i != nil {
		if len(i.Keywords) != 0 {
			if len(i.Keywords) != len(flat) {
				return ChkConflictingValues
			}
			smap := make(map[string]bool)
			for _, s := range i.Keywords {
				if len(s) > iptc.MaxKeywordLen {
					res = ChkIncorrectlyTagged
					s = s[:iptc.MaxKeywordLen]
				}
				if !flatmax[s] {
					return ChkConflictingValues
				}
				smap[s] = true
			}
			for s := range flatmax {
				if !smap[s] {
					return ChkConflictingValues
				}
			}
		} else if len(flatmax) != 0 {
			res = ChkIncorrectlyTagged
		}
	}
	if len(value) != 0 && res == 0 {
		return ChkPresent
	}
	return res
}
func keywordsEqual(a, b []metadata.Keyword) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j].Word != b[i][j].Word || a[i][j].OmitWhenFlattened != b[i][j].OmitWhenFlattened {
				return false
			}
		}
	}
	return true
}

// SetKeywords sets the creator tags with the specified prefix.  It leaves
// keywords with other prefixes alone.
func SetKeywords(h fileHandler, v []metadata.Keyword) error {
	var flat []string

	for _, kw := range v {
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
		xmp.DigiKamTagsList = v
		xmp.LRHierarchicalSubject = v
		xmp.DCSubject = flat
	}
	if iptc := h.IPTC(); iptc != nil {
		iptc.Keywords = flat
	}
	return nil
}
