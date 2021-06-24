package metadata

import (
	"errors"
	"strings"
)

// A Location is the textual description of a location.
type Location struct {
	Lang        string
	CountryCode string
	CountryName string
	State       string
	City        string
	Sublocation string
}

// ErrParseLocation is the error returned when a string cannot be parsed into a
// Location value.
var ErrParseLocation = errors.New("invalid Location value")

// Parse sets the value from the input string.  It returns ErrParseLocation if
// the input string is invalid.
func (loc *Location) Parse(val string) error {
	*loc = Location{}
	if val == "" {
		return nil
	}
	parts := strings.Split(val, "/")
	switch len(parts) {
	case 1:
		break
	case 2:
		loc.State = strings.TrimSpace(parts[1])
	case 3:
		loc.State = strings.TrimSpace(parts[1])
		loc.City = strings.TrimSpace(parts[2])
	case 4:
		loc.State = strings.TrimSpace(parts[1])
		loc.City = strings.TrimSpace(parts[2])
		loc.Sublocation = strings.TrimSpace(parts[3])
	default:
		return ErrParseLocation
	}
	cparts := strings.SplitN(parts[0], " ", 2)
	if loc.CountryName = countryCodes[cparts[0]]; loc.CountryName == "" {
		return ErrParseLocation
	}
	loc.CountryCode = cparts[0]
	if len(cparts) > 1 {
		if cname := strings.TrimSpace(cparts[1]); cname != "" {
			loc.CountryName = cname
		}
	}
	if loc.CountryCode == "USA" {
		if sname := stateCodes[loc.State]; sname != "" {
			loc.State = sname
		}
	}
	return nil
}

// String returns the value in string form, suitable for input to Parse.
func (loc *Location) String() string {
	var sb strings.Builder

	if loc.Empty() {
		return ""
	}
	sb.WriteString(loc.CountryCode)
	if loc.CountryName != "" {
		sb.WriteByte(' ')
		sb.WriteString(loc.CountryName)
	}
	if loc.State != "" || loc.City != "" || loc.Sublocation != "" {
		sb.WriteString(" /")
	}
	if loc.State != "" {
		sb.WriteByte(' ')
		sb.WriteString(loc.State)
	}
	if loc.State != "" && (loc.City != "" || loc.Sublocation != "") {
		sb.WriteByte(' ')
	}
	if loc.City != "" || loc.Sublocation != "" {
		sb.WriteByte('/')
	}
	if loc.City != "" {
		sb.WriteByte(' ')
		sb.WriteString(loc.City)
	}
	if loc.City != "" && loc.Sublocation != "" {
		sb.WriteByte(' ')
	}
	if loc.Sublocation != "" {
		sb.WriteString("/ ")
		sb.WriteString(loc.Sublocation)
	}
	return sb.String()
}

// ParseComponents fills in the supplied components of the Location.  It returns
// ErrParseLocation if they are invalid.
func (loc *Location) ParseComponents(countryCode, countryName, state, city, sublocation string) error {
	*loc = Location{}
	if countryCode == "" {
		if countryName != "" || state != "" || city != "" || sublocation != "" {
			return ErrParseLocation
		}
		return nil
	}
	if loc.CountryName = countryCodes[countryCode]; loc.CountryName == "" {
		return ErrParseLocation
	}
	loc.CountryCode = countryCode
	if countryName != "" {
		loc.CountryName = countryName
	}
	loc.State, loc.City, loc.Sublocation = state, city, sublocation
	return nil
}

// Empty returns true if the value contains no data.
func (loc *Location) Empty() bool {
	return loc == nil || loc.CountryCode == ""
}

// Equal returns true if the receiver is equal to the argument.
func (loc *Location) Equal(other *Location) bool {
	if (loc == nil) != (other == nil) {
		return false
	}
	if loc == nil {
		return true
	}
	return loc.CountryCode == other.CountryCode && loc.CountryName == other.CountryName && loc.State == other.State &&
		loc.City == other.City && loc.Sublocation == other.Sublocation
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
