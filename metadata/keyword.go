package metadata

type Keyword []KeywordComponent

type KeywordComponent struct {
	Word              string
	OmitWhenFlattened bool
}
