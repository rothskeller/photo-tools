package xmp

import (
	"github.com/rothskeller/photo-tools/metadata"
)

const nsPS = "http://ns.adobe.com/photoshop/1.0/"
const pfxPS = "photoshop"

// PSDateCreated returns the value of the photoshop:DateCreated tag.
func (p *XMP) PSDateCreated() metadata.DateTime { return p.psDateCreated }

func (p *XMP) getPS() {
	p.xmpDateTimeToMetadata(p.getString(p.rdf.Properties, pfxPS, nsPS, "DateCreated"), &p.psDateCreated)
	p.rdf.RegisterNamespace(pfxPS, nsPS)
}

// SetPSDateCreated sets the value of the photoshop:DateCreated tag.
func (p *XMP) SetPSDateCreated(v metadata.DateTime) (err error) {
	if v.Equivalent(p.psDateCreated) {
		return nil
	}
	p.psDateCreated = v
	p.setString(p.rdf.Properties, nsPS, "DateCreated", v.String())
	return nil
}
