package strmeta

import (
	"sort"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

// A keywordFilter is a predicate function used to filter keywords.
type keywordFilter func(metadata.Keyword) bool

// allKeywordsFilter is a keywordFilter that matches everything.
func allKeywordsFilter(_ metadata.Keyword) bool { return true }

// getFilteredKeywords returns the highest priority keywords that satisfy the
// specified filter.  If includeFlat is true, it includes flat keywords that
// satisfy the filter, if they are not implied by a hierarchical keyword.
func getFilteredKeywords(h filefmt.FileHandler, pred keywordFilter, includeFlat bool) (kws []metadata.Keyword) {
	var (
		all     []metadata.Keyword
		flat    []string
		flatmap = make(map[string]bool)
	)
	// First, get all keywords.
	if xmp := h.XMP(false); xmp != nil {
		if len(xmp.DigiKamTagsList()) != 0 {
			all = xmp.DigiKamTagsList()
		} else if len(xmp.LRHierarchicalSubject()) != 0 {
			all = xmp.LRHierarchicalSubject()
		}
		if len(xmp.DCSubject()) != 0 {
			flat = xmp.DCSubject()
		}
	}
	if iptc := h.IPTC(); iptc != nil && flat == nil {
		flat = iptc.Keywords()
	}
	for _, f := range flat {
		flatmap[f] = true
	}
	for _, kw := range all {
		delete(flatmap, kw[len(kw)-1])
	}
	// Now filter the hierarchical list with the predicate.
	for _, kw := range all {
		if pred(kw) {
			kws = append(kws, kw)
		}
	}
	// Finally, if asked, add flat keywords that are missing from the
	// hierarchical keywords, but only if they satisfy the predicate.
	if includeFlat {
		for word := range flatmap {
			var kw = metadata.Keyword{word}
			if pred(kw) {
				kws = append(kws, kw)
			}
		}
	}
	return kws
}

// getFilteredKeywordTags returns all keyword tags and their values, for those
// keywords that satisfy the specified filter.
func getFilteredKeywordTags(h filefmt.FileHandler, pred keywordFilter) (tags []string, values []metadata.Keyword) {
	if xmp := h.XMP(false); xmp != nil {
		for _, kw := range xmp.DigiKamTagsList() {
			if pred(kw) {
				tags = append(tags, "XMP  digiKam:TagsList")
				values = append(values, kw)
			}
		}
		for _, kw := range xmp.LRHierarchicalSubject() {
			if pred(kw) {
				tags = append(tags, "XMP  lr:HirerchicalSubject")
				values = append(values, kw)
			}
		}
		for _, s := range xmp.DCSubject() {
			var kw = metadata.Keyword{s}
			if pred(kw) {
				tags = append(tags, "XMP  dc:subject")
				values = append(values, kw)
			}
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		for _, s := range iptc.Keywords() {
			var kw = metadata.Keyword{s}
			if pred(kw) {
				tags = append(tags, "IPTC Keywords")
				values = append(values, kw)
			}
		}
	}
	return tags, values
}

// checkFilteredKeywords determines whether the keywords satisfying the
// specified filter are tagged correctly, and are consistent with those in the
// reference media file.
func checkFilteredKeywords(ref, tgt filefmt.FileHandler, pred keywordFilter) (res CheckResult) {
	var (
		values   = getFilteredKeywords(ref, pred, false)
		valuemap = make(map[string]bool)
	)
	for _, kw := range values {
		valuemap[kw.String()] = true
	}
	if xmp := tgt.XMP(false); xmp != nil {
		res = checkFilteredHierarchicalKeywords(valuemap, xmp.DigiKamTagsList(), pred)
		if r := checkFilteredHierarchicalKeywords(valuemap, xmp.LRHierarchicalSubject(), pred); r < res {
			res = r
		}
	}
	if res == 0 && len(values) != 0 {
		res = ChkPresent
	}
	return res
}
func checkFilteredHierarchicalKeywords(refmap map[string]bool, tgt []metadata.Keyword, pred keywordFilter) (res CheckResult) {
	var tgtmap = make(map[string]bool)

	for _, kw := range tgt {
		if pred(kw) {
			tgtmap[kw.String()] = true
		}
	}
	return checkMaps(refmap, tgtmap)
}
func checkMaps(refmap, tgtmap map[string]bool) (res CheckResult) {
	if len(refmap) == 0 && len(tgtmap) == 0 {
		return ChkOptionalAbsent
	}
	if len(refmap) != 0 && len(tgtmap) == 0 {
		return ChkIncorrectlyTagged
	}
	if len(refmap) == 0 && len(tgtmap) != 0 {
		return ChkConflictingValues
	}
	for kw := range refmap {
		if !tgtmap[kw] {
			return ChkConflictingValues
		}
	}
	for kw := range tgtmap {
		if !refmap[kw] {
			return ChkConflictingValues
		}
	}
	return ChkPresent
}

// setFilteredKeeywords replaces all keywords that match the specified filter
// with the specified list.  It leaves all other keywords alone.
func setFilteredKeywords(h filefmt.FileHandler, values []metadata.Keyword, pred keywordFilter) error {
	var (
		all     []metadata.Keyword
		set     []metadata.Keyword
		flatmap = make(map[string]bool)
		flat    []string
	)
	all = getFilteredKeywords(h, allKeywordsFilter, true)
	for _, kw := range all {
		if !pred(kw) && len(kw) != 0 {
			set = append(set, kw)
			flatmap[kw[len(kw)-1]] = true
		}
	}
	for _, kw := range values {
		if !pred(kw) {
			panic("setFilteredKeywords: values don't satisfy predicate")
		}
		set = append(set, kw)
		flatmap[kw[len(kw)-1]] = true
	}
	for kw := range flatmap {
		flat = append(flat, kw)
	}
	sort.Strings(flat)
	if xmp := h.XMP(true); xmp != nil {
		if err := xmp.SetDigiKamTagsList(set); err != nil {
			return err
		}
		if err := xmp.SetLRHierarchicalSubject(set); err != nil {
			return err
		}
		if err := xmp.SetDCSubject(flat); err != nil {
			return err
		}
	}
	if iptc := h.IPTC(); iptc != nil {
		if err := iptc.SetKeywords(flat); err != nil {
			return err
		}
	}
	return nil
}
