package strmeta

import (
	"github.com/rothskeller/photo-tools/filefmt"
)

func GetKeywords(h filefmt.FileHandler, prefix string, obscured bool) (values []string, change bool) {
	return nil, false
}

func AddKeyword(h filefmt.FileHandler, kw string) error {
	return nil
}

func RemoveKeyword(h filefmt.FileHandler, kw string) error {
	return nil
}

func RemoveAllKeywords(h filefmt.FileHandler, prefix string) error {
	return nil
}

func RecordOldPlaceKeywords(h filefmt.FileHandler) error {
	return nil
}

func RemoveOldPlaceKeywords(h filefmt.FileHandler) error {
	return nil
}
