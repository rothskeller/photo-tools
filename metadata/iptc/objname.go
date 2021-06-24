package iptc

import (
	"strings"
	"unicode/utf8"
)

// MaxObjectNameLen is the maximum length of the Object Name entry.
const MaxObjectNameLen = 64

const idObjectName uint16 = 0x0205

func (p *IPTC) getObjectName() {
	if dset := p.findDSet(idObjectName); dset != nil {
		if utf8.Valid(dset.data) {
			p.ObjectName = strings.TrimSpace(string(dset.data))
		} else {
			p.log("ignoring non-UTF8 Object Name")
		}
	}
}

func (p *IPTC) setObjectName() {
	if p.ObjectName == "" {
		p.deleteDSet(idObjectName)
		return
	}
	encoded := []byte(applyMax(p.ObjectName, MaxObjectNameLen))
	p.setDSet(idObjectName, encoded)
}
