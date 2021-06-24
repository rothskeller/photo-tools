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
			if objname := strings.TrimSpace(string(dset.data)); objname != "" {
				p.ObjectName.SetMaxLength(MaxObjectNameLen)
				p.ObjectName.Parse(objname)
			}
		} else {
			p.log("ignoring non-UTF8 Object Name")
		}
	}
}

func (p *IPTC) setObjectName() {
	if p.ObjectName.Empty() {
		p.deleteDSet(idObjectName)
		return
	}
	p.ObjectName.SetMaxLength(MaxObjectNameLen)
	encoded := []byte(p.ObjectName.String())
	p.setDSet(idObjectName, encoded)
}
