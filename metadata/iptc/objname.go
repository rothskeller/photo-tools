package iptc

import (
	"strings"
)

// MaxObjectNameLen is the maximum length of the Object Name entry.
const MaxObjectNameLen = 64

const idObjectName uint16 = 0x0205

// ObjectName returns the value of the Object Name tag.
func (p *IPTC) ObjectName() string { return p.objectName }

func (p *IPTC) getObjectName() {
	if dset := p.findDSet(idObjectName); dset != nil {
		p.objectName = strings.TrimSpace(p.decodeString(dset.data, "ObjectName"))
	}
}

// SetObjectName sets the value of the Object Name tag.
func (p *IPTC) SetObjectName(v string) error {
	if stringEqualMax(v, p.objectName, MaxObjectNameLen) {
		return nil
	}
	p.objectName = applyMax(v, MaxObjectNameLen)
	if p.objectName == "" {
		p.deleteDSet(idObjectName)
		return nil
	}
	p.setDSet(idObjectName, []byte(p.objectName))
	return nil
}
