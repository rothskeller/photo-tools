package metadata

import "strings"

// A Location is the textual description of a location.
type Location struct {
	CountryCode string
	CountryName string
	State       string
	City        string
	Sublocation string
}

// ParseLocation parses a string into a Location.  It returns nil if the string
// cannot be parsed (which is different from a pointer to an empty structure,
// which means the location is unknown/empty).
func ParseLocation(s string) (l *Location) {
	l = new(Location)
	s = strings.TrimSpace(s)
	if s == "" {
		return l
	}
	parts := strings.Split(s, "/")
	switch len(parts) {
	case 4:
		l.Sublocation = strings.TrimSpace(parts[3])
		fallthrough
	case 3:
		l.City = strings.TrimSpace(parts[2])
		fallthrough
	case 2:
		l.State = strings.TrimSpace(parts[1])
	case 1:
		break
	default:
		return nil
	}
	cparts := strings.SplitN(parts[0], " ", 2)
	if l.CountryName = countryCodes[cparts[0]]; l.CountryCode == "" {
		return nil
	}
	l.CountryCode = cparts[0]
	if len(cparts) > 1 {
		if cname := strings.TrimSpace(cparts[1]); cname != "" {
			l.CountryName = cname
		}
	}
	if l.CountryCode == "USA" {
		if sname := stateCodes[l.State]; sname != "" {
			l.State = sname
		}
	}
	return l
}

func (l *Location) String() string {
	var sb strings.Builder

	if l == nil || l.CountryCode == "" {
		return ""
	}
	sb.WriteString(l.CountryCode)
	if l.CountryName != "" {
		sb.WriteByte(' ')
		sb.WriteString(l.CountryName)
	}
	if l.State != "" || l.City != "" || l.Sublocation != "" {
		sb.WriteString(" /")
	}
	if l.State != "" {
		sb.WriteByte(' ')
		sb.WriteString(l.State)
	}
	if l.State != "" && (l.City != "" || l.Sublocation != "") {
		sb.WriteByte(' ')
	}
	if l.City != "" || l.Sublocation != "" {
		sb.WriteByte('/')
	}
	if l.City != "" {
		sb.WriteByte(' ')
		sb.WriteString(l.City)
	}
	if l.City != "" && l.Sublocation != "" {
		sb.WriteByte(' ')
	}
	if l.Sublocation != "" {
		sb.WriteString("/ ")
		sb.WriteString(l.Sublocation)
	}
	return sb.String()
}

// Valid returns whether the location contains any valid data.
func (l *Location) Valid() bool {
	return l != nil && l.CountryCode != ""
}

// Equal tests two locations for equality.
func (l *Location) Equal(o *Location) bool {
	if (l == nil) != (o == nil) {
		return false
	}
	if l == nil {
		return true
	}
	return l.CountryCode == o.CountryCode && l.CountryName == o.CountryName && l.State == o.State &&
		l.City == o.City && l.Sublocation == o.Sublocation
}

var countryCodes = map[string]string{
	"ASC":  "Ascension Island",
	"AND":  "Andorra",
	"ARE":  "United Arab Emirates",
	"AFG":  "Afghanistan",
	"ATG":  "Antigua and Barbuda",
	"AIA":  "Anguilla",
	"ALB":  "Albania",
	"ARM":  "Armenia",
	"ANHH": "Netherlands Antilles",
	"AGO":  "Angola",
	"ATA":  "Antarctica",
	"ARG":  "Argentina",
	"ASM":  "American Samoa",
	"AUT":  "Austria",
	"AUS":  "Australia",
	"ABW":  "Aruba",
	"ALA":  "\u212Bland Islands",
	"AZE":  "Azerbaijan",
	"BIH":  "Bosnia and Herzegovina",
	"BRB":  "Barbados",
	"BGD":  "Bangladesh",
	"BEL":  "Belgium",
	"BFA":  "Burkina Faso",
	"BGR":  "Bulgaria",
	"BHR":  "Bahrain",
	"BDI":  "Burundi",
	"BEN":  "Benin",
	"BLM":  "Saint Barth\u00E9lemy",
	"BMU":  "Bermuda",
	"BRN":  "Brunei Darussalam",
	"BOL":  "Bolivia",
	"BES":  "Bonaire, Sint Eustatius and Saba",
	"BRA":  "Brazil",
	"BHS":  "Bahamas",
	"BTN":  "Bhutan",
	"BUMM": "Burma",
	"BVT":  "Bouvet Island",
	"BWA":  "Botswana",
	"BLR":  "Belarus",
	"BLZ":  "Belize",
	"CAN":  "Canada",
	"CCK":  "Cocos Islands",
	"COD":  "Congo",
	"CAF":  "Central African Republic",
	"COG":  "Congo",
	"CHE":  "Switzerland",
	"CIV":  "C\u00F4te d'Ivoire",
	"COK":  "Cook Islands",
	"CHL":  "Chile",
	"CMR":  "Cameroon",
	"CHN":  "China",
	"COL":  "Colombia",
	"CPT":  "Clipperton Island",
	"CRI":  "Costa Rica",
	"CSXX": "Serbia and Montenegro",
	"CUB":  "Cuba",
	"CPV":  "Cape Verde",
	"CUW":  "Cura\u00E7ao",
	"CXR":  "Christmas Island",
	"CYP":  "Cyprus",
	"CZE":  "Czech Republic",
	"DEU":  "Germany",
	"DGA":  "Diego Garcia",
	"DJI":  "Djibouti",
	"DNK":  "Denmark",
	"DMA":  "Dominica",
	"DOM":  "Dominican Republic",
	"DZA":  "Algeria",
	"ECU":  "Ecuador",
	"EST":  "Estonia",
	"EGY":  "Egypt",
	"ESH":  "Western Sahara",
	"ERI":  "Eritrea",
	"ESP":  "Spain",
	"ETH":  "Ethiopia",
	"FIN":  "Finland",
	"FJI":  "Fiji",
	"FLK":  "Falkland Islands",
	"FSM":  "Micronesia",
	"FRO":  "Faroe Islands",
	"FRA":  "France",
	"GAB":  "Gabon",
	"GBR":  "United Kingdom",
	"GRD":  "Grenada",
	"GEO":  "Georgia",
	"GUF":  "French Guiana",
	"GGY":  "Guernsey",
	"GHA":  "Ghana",
	"GIB":  "Gibraltar",
	"GRL":  "Greenland",
	"GMB":  "Gambia",
	"GIN":  "Guinea",
	"GLP":  "Guadeloupe",
	"GNQ":  "Equatorial Guinea",
	"GRC":  "Greece",
	"SGS":  "South Georgia and the South Sandwich Islands",
	"GTM":  "Guatemala",
	"GUM":  "Guam",
	"GNB":  "Guinea-Bissau",
	"GUY":  "Guyana",
	"HKG":  "Hong Kong",
	"HMD":  "Heard Island and McDonald Islands",
	"HND":  "Honduras",
	"HRV":  "Croatia",
	"HTI":  "Haiti",
	"HUN":  "Hungary",
	"IDN":  "Indonesia",
	"IRL":  "Ireland",
	"ISR":  "Israel",
	"IMN":  "Isle of Man",
	"IND":  "India",
	"IOT":  "British Indian Ocean Territory",
	"IRQ":  "Iraq",
	"IRN":  "Iran",
	"ISL":  "Iceland",
	"ITA":  "Italy",
	"JEY":  "Jersey",
	"JAM":  "Jamaica",
	"JOR":  "Jordan",
	"JPN":  "Japan",
	"KEN":  "Kenya",
	"KGZ":  "Kyrgyzstan",
	"KHM":  "Cambodia",
	"KIR":  "Kiribati",
	"COM":  "Comoros",
	"KNA":  "Saint Kitts and Nevis",
	"PRK":  "North Korea",
	"KOR":  "South Korea",
	"KWT":  "Kuwait",
	"CYM":  "Cayman Islands",
	"KAZ":  "Kazakhstan",
	"LAO":  "Lao",
	"LBN":  "Lebanon",
	"LCA":  "Saint Lucia",
	"LIE":  "Liechtenstein",
	"LKA":  "Sri Lanka",
	"LBR":  "Liberia",
	"LSO":  "Lesotho",
	"LTU":  "Lithuania",
	"LUX":  "Luxembourg",
	"LVA":  "Latvia",
	"LBY":  "Libya",
	"MAR":  "Morocco",
	"MCO":  "Monaco",
	"MDA":  "Moldova",
	"MNE":  "Montenegro",
	"MAF":  "Saint Martin",
	"MDG":  "Madagascar",
	"MHL":  "Marshall Islands",
	"MKD":  "Macedonia",
	"MLI":  "Mali",
	"MMR":  "Myanmar",
	"MNG":  "Mongolia",
	"MAC":  "Macao",
	"MNP":  "Northern Mariana Islands",
	"MTQ":  "Martinique",
	"MRT":  "Mauritania",
	"MSR":  "Montserrat",
	"MLT":  "Malta",
	"MUS":  "Mauritius",
	"MDV":  "Maldives",
	"MWI":  "Malawi",
	"MEX":  "Mexico",
	"MYS":  "Malaysia",
	"MOZ":  "Mozambique",
	"NAM":  "Namibia",
	"NCL":  "New Caledonia",
	"NER":  "Niger",
	"NFK":  "Norfolk Island",
	"NGA":  "Nigeria",
	"NIC":  "Nicaragua",
	"NLD":  "Netherlands",
	"NOR":  "Norway",
	"NPL":  "Nepal",
	"NRU":  "Nauru",
	"NTHH": "Neutral Zone",
	"NIU":  "Niue",
	"NZL":  "New Zealand",
	"OMN":  "Oman",
	"PAN":  "Panama",
	"PER":  "Peru",
	"PYF":  "French Polynesia",
	"PNG":  "Papua New Guinea",
	"PHL":  "Philippines",
	"PAK":  "Pakistan",
	"POL":  "Poland",
	"SPM":  "Saint Pierre and Miquelon",
	"PCN":  "Pitcairn",
	"PRI":  "Puerto Rico",
	"PSE":  "Palestine",
	"PRT":  "Portugal",
	"PLW":  "Palau",
	"PRY":  "Paraguay",
	"QAT":  "Qatar",
	"REU":  "R\u00E9union",
	"ROU":  "Romania",
	"SRB":  "Serbia",
	"RUS":  "Russian Federation",
	"RWA":  "Rwanda",
	"SAU":  "Saudi Arabia",
	"SLB":  "Solomon Islands",
	"SYC":  "Seychelles",
	"SDN":  "Sudan",
	"SWE":  "Sweden",
	"SGP":  "Singapore",
	"SHN":  "Saint Helena, Ascension and Tristan da Cunha",
	"SVN":  "Slovenia",
	"SJM":  "Svalbard and Jan Mayen",
	"SVK":  "Slovakia",
	"SLE":  "Sierra Leone",
	"SMR":  "San Marino",
	"SEN":  "Senegal",
	"SOM":  "Somalia",
	"SUR":  "Suriname",
	"SSD":  "South Sudan",
	"STP":  "Sao Tome and Principe",
	"SUN":  "USSR",
	"SLV":  "El Salvador",
	"SXM":  "Sint Maarten",
	"SYR":  "Syrian Arab Republic",
	"SWZ":  "Swaziland",
	"TAA":  "Tristan da Cunha",
	"TCA":  "Turks and Caicos Islands",
	"TCD":  "Chad",
	"ATF":  "French Southern Territories",
	"TGO":  "Togo",
	"THA":  "Thailand",
	"TJK":  "Tajikistan",
	"TKL":  "Tokelau",
	"TLS":  "Timor-Leste",
	"TKM":  "Turkmenistan",
	"TUN":  "Tunisia",
	"TON":  "Tonga",
	"TPTL": "East Timor",
	"TUR":  "Turkey",
	"TTO":  "Trinidad and Tobago",
	"TUV":  "Tuvalu",
	"TWN":  "Taiwan",
	"TZA":  "Tanzania",
	"UKR":  "Ukraine",
	"UGA":  "Uganda",
	"USA":  "United States",
	"URY":  "Uruguay",
	"UZB":  "Uzbekistan",
	"VAT":  "Vatican",
	"VCT":  "Saint Vincent and the Grenadines",
	"VEN":  "Venezuela",
	"VGB":  "British Virgin Islands",
	"VNM":  "Viet Nam",
	"VUT":  "Vanuatu",
	"WLF":  "Wallis and Futuna",
	"WSM":  "Samoa",
	"XXK":  "Kosovo",
	"YEM":  "Yemen",
	"MYT":  "Mayotte",
	"YUCS": "Yugoslavia",
	"ZAF":  "South Africa",
	"ZMB":  "Zambia",
	"ZRCD": "Zaire",
	"ZWE":  "Zimbabwe",
}

var stateCodes = map[string]string{
	"AK": "Alaska",
	"AL": "Alabama",
	"AR": "Arkansas",
	"AZ": "Arizona",
	"CA": "California",
	"CO": "Colorado",
	"CT": "Connecticut",
	"DC": "District of Columbia",
	"DE": "Delaware",
	"FL": "Florida",
	"GA": "Georgia",
	"HI": "Hawaii",
	"IA": "Iowa",
	"ID": "Idaho",
	"IL": "Illinois",
	"IN": "Indiana",
	"KS": "Kansas",
	"KY": "Kentucky",
	"LA": "Louisiana",
	"MA": "Massachusetts",
	"MD": "Maryland",
	"ME": "Maine",
	"MI": "Michigan",
	"MN": "Minnesota",
	"MO": "Missouri",
	"MS": "Mississippi",
	"MT": "Montana",
	"NC": "North Carolina",
	"ND": "North Dakota",
	"NE": "Nebraska",
	"NH": "New Hampshire",
	"NJ": "New Jersey",
	"NM": "New Mexico",
	"NV": "Nevada",
	"NY": "New York",
	"OH": "Ohio",
	"OK": "Oklahoma",
	"OR": "Oregon",
	"PA": "Pennsylvania",
	"PR": "Puerto Rico",
	"RI": "Rhode Island",
	"SC": "South Carolina",
	"SD": "South Dakota",
	"TN": "Tennessee",
	"TX": "Texas",
	"UT": "Utah",
	"VA": "Virginia",
	"VI": "Virgin Islands",
	"VT": "Vermont",
	"WA": "Washington",
	"WI": "Wisconsin",
	"WV": "West Virginia",
	"WY": "Wyoming",
}
