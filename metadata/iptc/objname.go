package iptc

import (
	"strings"
)

// MaxObjectNameLen is the maximum length of the Object Name entry.
const MaxObjectNameLen = 64

const idObjectName uint16 = 0x0205

func (p *IPTC) getObjectName() {
	if dset := p.findDSet(idObjectName); dset != nil {
		p.ObjectName = strings.TrimSpace(p.decodeString(dset.data, "ObjectName"))
		p.saveObjectName = p.ObjectName
	}
}

func (p *IPTC) setObjectName() {
	if stringEqualMax(p.ObjectName, p.saveObjectName, MaxObjectNameLen) {
		return
	}
	if p.ObjectName == "" {
		p.deleteDSet(idObjectName)
		return
	}
	p.setDSet(idObjectName, []byte(applyMax(p.ObjectName, MaxObjectNameLen)))
}
