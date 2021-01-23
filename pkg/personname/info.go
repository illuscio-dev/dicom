package personname

import (
	"fmt"
	"strings"
)

const groupSep = "="

// pnGroup is an enum value for the PN group types (Alphabetic, Ideographic & Phonetic).
type pnGroup int

// String representation -- mostly for formatting error messages.
func (group pnGroup) String() string {
	switch group {
	case pnGroupAlphabetic:
		return "Alphabetic"
	case pnGroupIdeographic:
		return "Ideographic"
	case pnGroupPhonetic:
		return "Phonetic"
	default:
		panic(fmt.Errorf("bad pnGroup value: %v", int(group)))
	}
}

// Enum values for pnGroup
const (
	pnGroupAlphabetic  = 0
	pnGroupIdeographic = 1
	pnGroupPhonetic    = 2
)

// InfoTrailingNullLevel represents how many null '=' separators are present in the
// Info.DCM() return value.
type InfoTrailingNullLevel uint

// String implements fmt.Stringer, giving human-readable names to the null sep level.
//
// Returns "NONE" if no null separators were present.
//
// Returns "ALL" if the highest possible null separator was present.
//
// Otherwise, returns the name of the section that comes after the highest present null
// separator.
//
// String will panic if called on a value that exceeds InfoNullAll.
func (level InfoTrailingNullLevel) String() string {
	switch level {
	case InfoNullNone:
		return "NONE"
	case InfoNullIdeographic:
		return "Ideographic"
	case InfoNullAll:
		return "ALL"
	default:
		panic(validateInfoNullSepLevel(level))
	}
}

const (
	// InfoNullNone will render no null seps.
	InfoNullNone InfoTrailingNullLevel = iota

	// InfoNullIdeographic will render null separators up to the separator before the
	// Info.Ideographic segment
	InfoNullIdeographic

	// InfoNullAll will render null separators up to the separator before the
	// Info.Phonetic segment, or ALL possible separators.
	InfoNullAll
)

func validateInfoNullSepLevel(level InfoTrailingNullLevel) error {
	if level <= InfoNullAll {
		return nil
	}

	return newErrNullSepLevelInvalid(uint(InfoNullAll), uint(level))
}

// Info holds information from an element with a "PN" VR. See the "PN"
// entry at: http://dicom.nema.org/dicom/2013/output/chtml/part05/sect_6.2.html
//
// PN values are split into three groups which represent three different ways to
// represent a name:
//
// - Alphabetic: How a name is formally spelled using a Phonetic alphabet.
// - Ideographic: How a name is represented using ideograms / ideographs.
// - Phonetic: How a name is pronounced.
//
// Each of these groups can be inspected to access their individual segments (family
// name, Given name, etc.)
type Info struct {
	// Alphabetic group information about the Alphabetic group.
	Alphabetic GroupInfo
	// Ideographic group information about the Ideographic group.
	Ideographic GroupInfo
	// Phonetic group information about the Phonetic group.
	Phonetic GroupInfo

	// TrailingNullLevel contains the highest present null '=' separator in the DCM()
	// value. For most use cases InfoNullAll or InfoNullNone should be used when
	// creating new PN values. Use other levels only if you know what you are doing!
	TrailingNullLevel InfoTrailingNullLevel
}

// WithFormat returns a new Info object with null separator settings applied to the
// relevant Info / GroupInfo objects.
//
// infoNullSepLevel sets the highest '=' null separator to render between groups.
//
// The remaining options will apply the passed value to their groups respective
// TrailingNullLevel value.
//
// WithFormat does not mutate its receiver value, instead returning a new value
// to the caller with the passed settings.
func (info Info) WithFormat(
	infoNullSepLevel InfoTrailingNullLevel,
	alphabeticNullSepLevel,
	ideographicNullSepLevel,
	phoneticNullSepLevel GroupTrailingNullLevel,
) Info {
	info.TrailingNullLevel = infoNullSepLevel
	info.Alphabetic.TrailingNullLevel = alphabeticNullSepLevel
	info.Ideographic.TrailingNullLevel = ideographicNullSepLevel
	info.Phonetic.TrailingNullLevel = phoneticNullSepLevel
	return info
}

// WithTrailingNulls returns a new Info object that will keep trailing separators
// that surround both null groups AND group segments: (ex: 'Potter^Harry^^^==').
//
// WithTrailingNulls is equivalent to calling WithFormat() with all options set to
// GroupNullAll.
//
// WithTrailingNulls does not mutate its receiver value, instead returning a new value
// to the caller with the passed settings.
func (info *Info) WithTrailingNulls() Info {
	// We're going to take in a pointer here to avoid a double-copy when invoking
	// WithFormat, otherwise we would be passing by value twice.
	//
	// Since WithFormat is pass-by-value, this will not mutate the original info.
	return info.WithFormat(
		InfoNullAll,
		GroupNullAll,
		GroupNullAll,
		GroupNullAll,
	)
}

// WithoutTrailingNulls returns a new Info object that will remove trailing
// separators that surround both null groups AND group segments:
// (ex: 'Potter^Harry').
//
// WithoutTrailingNulls is equivalent to calling WithFormat() with all options set to
// GroupNullNone.
//
// WithoutTrailingNulls does not mutate its receiver value, instead returning a new
// value to the caller with the passed settings.
func (info *Info) WithoutTrailingNulls() Info {
	// We're going to take in a pointer here to avoid a double-copy when invoking
	// WithFormat, otherwise we would be passing by value twice.
	//
	// Since WithFormat is pass-by-value, this will not mutate the original info.
	return info.WithFormat(
		InfoNullNone,
		GroupNullNone,
		GroupNullNone,
		GroupNullNone,
	)
}

// WithoutEmptyGroups sets Info.TrailingNullLevel to false, then checks eac
// group, and if it contains no actual information, sets that group's TrailingNullLevel
// to false.
//
// Groups with Partial information will retain their null separators.
func (info Info) WithoutEmptyGroups() Info {
	info.TrailingNullLevel = InfoNullNone

	// Iterate over references to our group values (we aren't mutating our receiver
	// here since it's passed by value and already a deep copy).
	for _, group := range []*GroupInfo{&info.Alphabetic, &info.Ideographic, &info.Phonetic} {
		if group.IsEmpty() {
			group.TrailingNullLevel = GroupNullNone
		}
	}

	return info
}

// DCM returns the serialized DICOM representation of the PN value, in
// '[Alphabetic]=[Ideographic]=[Phonetic]' format.
func (info Info) DCM() string {
	// validate our TrailingNullLevel and panic if it is exceeded.
	if err := validateInfoNullSepLevel(info.TrailingNullLevel); err != nil {
		panic(err)
	}

	// Convert the groups into their formatted string representations.=
	groupStrings := []string{
		info.Alphabetic.DCM(), info.Ideographic.DCM(), info.Phonetic.DCM(),
	}

	return renderWithSeps(groupStrings, groupSep, uint(info.TrailingNullLevel))
}

// IsEmpty returns whether the PN value contains any actual information. This method
// ignores separator, so both '' and '^^^^=^^^^=^^^^' would return true.
func (info Info) IsEmpty() bool {
	return info.Alphabetic.IsEmpty() &&
		info.Ideographic.IsEmpty() &&
		info.Phonetic.IsEmpty()
}

// Parse PN dicom value into a personname.Info value.
//
// NOTE ON PARSING:
//
// The personname.Info and personname.GroupInfo values only track whether any null
// separators were used, not how many. This means if a PN value has some null
// separators, but not the full amount, round-tripping the value will result in adding
// the missing separators. See examples below. If you wish to make sure that NO
// alterations are made to the original value after inspecting re-serializing, the
// original value should be used directly instead.
func Parse(valueString string) (Info, error) {
	groups := strings.Split(valueString, groupSep)

	// If there are more than three groups, then the value does not conform to the DICOM
	// spec.
	if len(groups) > 3 {
		return Info{}, newErrTooManyGroups(len(groups))
	}

	// Set up a new info value.
	info := Info{}

	// Range over the groups and assign them based on index.
	// Start off with our null segment level being None
	nullSepLevel := InfoNullNone
	for i, groupString := range groups {
		// If this group is empty, it means there is a null sep here. Our null sep
		// level needs to reflect this.
		if groupString == "" {
			nullSepLevel = InfoTrailingNullLevel(i)
		} else {
			// Otherwise, if there is a non-zero string value, there is no null sep
			// after it.
			nullSepLevel = InfoNullNone
		}

		// Convert the index into a pnGroup enum.
		group := pnGroup(i)

		// Parse the group info.
		groupInfo, err := groupFromValueString(groupString, group)
		if err != nil {
			return Info{}, err
		}

		// Apply the group info.
		switch group {
		case pnGroupAlphabetic:
			info.Alphabetic = groupInfo
		case pnGroupIdeographic:
			info.Ideographic = groupInfo
		case pnGroupPhonetic:
			info.Phonetic = groupInfo
		}
	}

	info.TrailingNullLevel = nullSepLevel

	return info, nil
}
