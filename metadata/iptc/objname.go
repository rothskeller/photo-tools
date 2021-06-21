package iptc

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

// MaxObjectNameLen is the maximum length of the Object Name entry.
const MaxObjectNameLen = 64

const idObjectName uint16 = 0x0205

// ObjectName returns the IPTC ObjectName tag, if any.
func (p *IPTC) ObjectName() string {
	if dset := p.findDSet(idObjectName); dset != nil {
		if utf8.Valid(dset.data) {
			return strings.TrimSpace(string(dset.data))
		}
		p.problems = append(p.problems, "ignoring non-UTF8 IPTC Object Name")
	}
	return ""
}

// SetObjectName sets the IPTC ObjectName.
func (p *IPTC) SetObjectName(name string) {
	if p == nil {
		return
	}
	if name == "" {
		p.deleteDSet(idObjectName)
		return
	}
	dset := p.findDSet(idObjectName)
	encoded := []byte(applyMax(name, MaxObjectNameLen))
	if dset != nil {
		if !bytes.Equal(encoded, dset.data) {
			dset.data = encoded
			p.dirty = true
		}
	} else {
		p.dsets = append(p.dsets, &dsett{0, idObjectName, encoded})
		p.dirty = true
	}
}
