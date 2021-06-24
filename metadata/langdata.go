package metadata

// LangDatum implements the Metadatum interface for a datum with an associated
// language tag.
type LangDatum struct {
	Metadatum
	Lang string
}

// Verify interface compliance.
var _ Metadatum = (*LangDatum)(nil)

// Multilingual represents a collection of LangDatum values for different
// languages, with the first one being considered the "default" language.
type Multilingual []*LangDatum

// Duolingual represents (up to) two language variants for a datum, one of which
// must be English.
type Duolingual struct {
	English LangDatum
	Other   LangDatum
}
