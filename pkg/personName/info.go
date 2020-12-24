/*
The personName package provides methods and data types for inspecting Person Name (PN)
Value Representations, as defined here:

http://dicom.nema.org/dicom/2013/output/chtml/part05/sect_6.2.html
*/
package personName

import (
	"errors"
	"regexp"
)

// Parsed Info value information from an element with a "PN" VR. See the "PN"
// entry at: http://dicom.nema.org/dicom/2013/output/chtml/part05/sect_6.2.html
type Info struct {
	// The person's family or last name.
	FamilyName string
	// The person's given or first names.
	GivenName string
	// The person's middle names.
	MiddleName string
	// The person's name prefix.
	NamePrefix string
	// The person's name suffix.
	NameSuffix string
}

// Returns dicom format PN string: '[Last]^[First]^[Middle]^[Prefix]^[Suffix]'.
func (pn Info) String() string {
	return pn.FamilyName +
		"^" + pn.GivenName +
		"^" + pn.MiddleName +
		"^" + pn.NamePrefix +
		"^" + pn.NameSuffix
}

// Regex for parsing PN values.
var pnRegex = regexp.MustCompile(
	`(?P<Fist>.*)\^(?P<Last>.*)\^(?P<Middle>.*)\^(?P<Prefix>.*)\^(?P<Suffix>.*)`,
)

// ErrParsePersonName is returned when attempting to parse Info from a string.
var ErrParsePersonName = errors.New("person name value does not match Dicom Spec" +
	" '[Last]^[First]^[Middle]^[Prefix]^[Suffix]'")

// FromDicomValueString converts a string from a dicom element with a Value
// Representation of PN to a parsed Info struct.
func FromDicomValueString(pnString string) (Info, error) {
	// Run the regex against the name.
	matches := pnRegex.FindStringSubmatch(pnString)

	// If the first match is empty or there are no matches, the value is malformed.
	if len(matches) == 0 || matches[0] == "" {
		return Info{}, ErrParsePersonName
	}

	name := Info{
		FamilyName: matches[1],
		GivenName:  matches[2],
		MiddleName: matches[3],
		NamePrefix: matches[4],
		NameSuffix: matches[5],
	}

	return name, nil
}
