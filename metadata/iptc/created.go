package iptc

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	idDateCreated uint16 = 0x0237
	idTimeCreated uint16 = 0x023C
)

var dateCreatedRE = regexp.MustCompile(`^\d\d\d\d\d\d\d\d$`)
var timeCreatedRE = regexp.MustCompile(`^\d\d\d\d\d\d(?:[-+]\d\d\d\d)?$`)

// DateTimeCreated returns the IPTC Date Created and Time Created tags, if any.
func (p *IPTC) DateTimeCreated() (dtc string) {
	datedset := p.findDSet(idDateCreated)
	if datedset == nil {
		return ""
	} else if !dateCreatedRE.Match(datedset.data) {
		p.problems = append(p.problems, fmt.Sprintf("invalid IPTC DateCreated value %q", datedset.data))
		return ""
	} else if datedset.data[0] == '0' {
		return ""
	}
	dc := string(datedset.data)
	dtc = dc[0:4]
	if dc[4:6] == "00" {
		dtc += "-01"
	} else {
		dtc += "-" + dc[4:6]
	}
	if dc[6:8] == "00" {
		dtc += "-01"
	} else {
		dtc += "-" + dc[6:8]
	}
	timedset := p.findDSet(idTimeCreated)
	if timedset == nil {
		return dtc + "T00:00:00"
	} else if !timeCreatedRE.Match(timedset.data) {
		p.problems = append(p.problems, fmt.Sprintf("invalid IPTC TimeCreated value %q", timedset.data))
		return dtc + "T00:00:00"
	}
	tc := string(timedset.data)
	if len(tc) == 11 {
		if strings.HasSuffix(tc, "0000") {
			return fmt.Sprintf("%sT%s:%s:%sZ", dtc, tc[0:2], tc[2:4], tc[4:6])
		}
		return fmt.Sprintf("%sT%s:%s:%s:%s", dtc, tc[0:2], tc[2:4], tc[4:9], tc[9:11])
	}
	return fmt.Sprintf("%sT%s:%s:%s", dtc, tc[0:2], tc[2:4], tc[4:6])
}

// SetDateTimeCreated sets the IPTC Date Created and Time Created.
func (p *IPTC) SetDateTimeCreated(val string) {
	if p == nil {
		return
	}
	if val == "" {
		p.deleteDSet(idDateCreated)
		p.deleteDSet(idTimeCreated)
		return
	}
	newDate := strings.Replace(val, "-", "", -1)[0:8]
	newTime := strings.Replace(val[11:], ":", "", -1)
	if strings.HasSuffix(newTime, "Z") {
		newTime = newTime[:len(newTime)-1] + "+0000"
	}
	datedset := p.findDSet(idDateCreated)
	if datedset != nil {
		if string(datedset.data) != newDate {
			datedset.data = []byte(newDate)
			p.dirty = true
		}
	} else {
		p.dsets = append(p.dsets, &dsett{0, idDateCreated, []byte(newDate)})
		p.dirty = true
	}
	timedset := p.findDSet(idTimeCreated)
	if timedset != nil {
		if string(timedset.data) != newTime {
			timedset.data = []byte(newTime)
			p.dirty = true
		}
	} else {
		p.dsets = append(p.dsets, &dsett{0, idTimeCreated, []byte(newTime)})
		p.dirty = true
	}
}
