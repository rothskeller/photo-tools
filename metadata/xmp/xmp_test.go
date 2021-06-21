package xmp

import (
	"testing"
)

var start = []byte(`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"></rdf:RDF>`)

func TestAddCreators(t *testing.T) {
	xmp := Parse(start)
	xmp.SetDCCreator([]string{"Steve Roth", "Fred Flintstone"})
	out, _ := xmp.Render()
	println(string(out))
	t.Error("x")
}
