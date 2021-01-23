package personname

import (
	"strings"
)

const segmentSep = "^"

// GroupTrailingNullLevel represents how many null '^' separators are present in the
// GroupInfo.DCM() return value.
type GroupTrailingNullLevel uint

// String implements fmt.Stringer, giving human-readable names to the trailing null
// level.
//
// Returns "NONE" if no null separators were present.
//
// Returns "ALL" if the highest possible null separator was present.
//
// Otherwise, returns the name of the section that comes after the highest present null
// separator.
//
// String will panic if called on a value that exceeds GroupNullAll.
func (level GroupTrailingNullLevel) String() string {
	switch level {
	case GroupNullNone:
		return "NONE"
	case GroupNullGiven:
		return "GivenName"
	case GroupNullMiddle:
		return "MiddleName"
	case GroupNullPrefix:
		return "NamePrefix"
	case GroupNullAll:
		return "ALL"
	default:
		panic(validateGroupNullSepLevel(level))
	}
}

const (
	// GroupNullNone will render no null seps.
	GroupNullNone GroupTrailingNullLevel = iota

	// NullSepGiven will render null separators up to the separator before the
	// GroupInfo.GivenName segment
	GroupNullGiven

	// NullSepGiven will render null separators up to the separator before the
	// GroupInfo.MiddleName segment
	GroupNullMiddle

	// NullSepGiven will render null separators up to the separator before the
	// GroupInfo.NamePrefix segment
	GroupNullPrefix

	// NullSepGiven will render null separators up to the separator before the
	// GroupInfo.NameSuffix segment (ALL possible separators).
	GroupNullAll
)

func validateGroupNullSepLevel(level GroupTrailingNullLevel) error {
	if level <= GroupNullAll {
		return nil
	}

	return newErrNullSepLevelInvalid(uint(GroupNullAll), uint(level))
}

// GroupInfo holds the parsed information for any one of these groups the person name
// groups specified in the DICOM spec:
//
//	- Alphabetic
//	- Ideographic
//	- Phonetic
type GroupInfo struct {
	// FamilyName is the person's family or last name.
	FamilyName string
	// GivenName if the person's given or first names.
	GivenName string
	// MiddleName is the person's middle name(s).
	MiddleName string
	// NamePrefix is the person's name prefix (ex: MR., MRS.).
	NamePrefix string
	// NameSuffix is the person's name suffix (ex: Jr, III).
	NameSuffix string

	// TrailingNullLevel contains the highest present null '^' separator in the DCM()
	// value. For most use cases GroupNullAll or GroupNullNone should be used when
	// creating new PN values. Use other levels only if you know what you are doing!
	TrailingNullLevel GroupTrailingNullLevel
}

// DCM Returns original, formatted string in
// '[FamilyName]^[GivenName]^[MiddleName]^[NamePrefix]^[NameSuffix]'.
func (group GroupInfo) DCM() string {
	// validate our TrailingNullLevel and panic if it is exceeded.
	if err := validateGroupNullSepLevel(group.TrailingNullLevel); err != nil {
		panic(err)
	}

	// Put all the segments into an array.
	segments := []string{
		group.FamilyName,
		group.GivenName,
		group.MiddleName,
		group.NamePrefix,
		group.NameSuffix,
	}

	// Render our segments with the correct number of null-separators.
	return renderWithSeps(segments, segmentSep, uint(group.TrailingNullLevel))
}

// IsEmpty returns true if all group segments are empty, even if Raw value was "^^^^".
func (group GroupInfo) IsEmpty() bool {
	return group.FamilyName == "" &&
		group.GivenName == "" &&
		group.MiddleName == "" &&
		group.NamePrefix == "" &&
		group.NameSuffix == ""
}

// groupFromValueString converts a string from a dicom element with a Value
// Representation of PN to a parsed Info struct.
func groupFromValueString(groupString string, group pnGroup) (GroupInfo, error) {
	segments := strings.Split(groupString, segmentSep)
	segmentCount := len(segments)

	if segmentCount > 5 {
		return GroupInfo{}, newErrTooManyGroupSegments(group, len(segments))
	}

	groupInfo := GroupInfo{}

	// Start off with our null segment level being None
	nullSepLevel := GroupNullNone
	for i, groupValue := range segments {
		// If this segment is empty, it means there is a null sep here. Our null sep
		// level needs to reflect this.
		if groupValue == "" {
			nullSepLevel = GroupTrailingNullLevel(i)
		} else {
			// Otherwise, if there is a non-zero string value, there is no null sep
			// after it.
			nullSepLevel = GroupNullNone
		}

		switch i {
		case 0:
			groupInfo.FamilyName = groupValue
		case 1:
			groupInfo.GivenName = groupValue
		case 2:
			groupInfo.MiddleName = groupValue
		case 3:
			groupInfo.NamePrefix = groupValue
		case 4:
			groupInfo.NameSuffix = groupValue
		}
	}

	// If the string is not empty, and any of our groups ARE empty, then we are using
	// null separators.
	if strings.HasSuffix(groupString, "^") {
		groupInfo.TrailingNullLevel = nullSepLevel
	}

	return groupInfo, nil
}
