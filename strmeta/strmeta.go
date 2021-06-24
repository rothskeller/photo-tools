// Package strmeta contains the definition of a Metadata structure reflecting my
// desired media metadata model, and the code to read it from and write it to
// real metadata sources.
package strmeta

import (
	"fmt"

	"github.com/rothskeller/photo-tools/filefmt"
	"github.com/rothskeller/photo-tools/metadata"
)

type fileHandler = filefmt.FileHandler // copied to save typing

func tagsForStringList(tags, values []string, label string, ss []string) (newt, newv []string) {
	if len(ss) == 0 {
		tags = append(tags, label)
		values = append(values, "")
	}
	for _, s := range ss {
		tags = append(tags, label)
		values = append(values, s)
	}
	return tags, values
}

func tagsForLangStrings(tags, values []string, label string, lss []metadata.LangString) (newt, newv []string) {
	if len(lss) == 0 {
		tags = append(tags, label)
		values = append(values, "")
	}
	for i, ls := range lss {
		if i == 0 && ls.Lang == "" {
			tags = append(tags, label)
		} else {
			tags = append(tags, fmt.Sprintf("%s[%s]", label, ls.Lang))
		}
		values = append(values, ls.Value)
	}
	return tags, values
}
