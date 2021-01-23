package personname

import (
	"errors"
	"testing"
)

func TestInfo(t *testing.T) {
	testCases := []struct {
		// The Raw string to parse from
		Raw string
		// The parsed information we expect.
		Expected Info
		// Whether IsEmpty should return true after parsing Raw.
		IsEmpty bool
	}{
		// All groups
		{
			Raw: "aFamily^aGiven^aMiddle^aPrefix^aSuffix=" +
				"iFamily^iGiven^iMiddle^iPrefix^iSuffix=" +
				"pFamily^pGiven^pMiddle^pPrefix^pSuffix",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "aFamily",
					GivenName:    "aGiven",
					MiddleName:   "aMiddle",
					NamePrefix:   "aPrefix",
					NameSuffix:   "aSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "iFamily",
					GivenName:    "iGiven",
					MiddleName:   "iMiddle",
					NamePrefix:   "iPrefix",
					NameSuffix:   "iSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "pFamily",
					GivenName:    "pGiven",
					MiddleName:   "pMiddle",
					NamePrefix:   "pPrefix",
					NameSuffix:   "pSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepNone,
			},
			IsEmpty: false,
		},
		// No Phonetic
		{
			Raw: "aFamily^aGiven^aMiddle^aPrefix^aSuffix=" +
				"iFamily^iGiven^iMiddle^iPrefix^iSuffix=" +
				"",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "aFamily",
					GivenName:    "aGiven",
					MiddleName:   "aMiddle",
					NamePrefix:   "aPrefix",
					NameSuffix:   "aSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "iFamily",
					GivenName:    "iGiven",
					MiddleName:   "iMiddle",
					NamePrefix:   "iPrefix",
					NameSuffix:   "iSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepAll,
			},
			IsEmpty: false,
		},
		// No Phonetic, no seps
		{
			Raw: "aFamily^aGiven^aMiddle^aPrefix^aSuffix=" +
				"iFamily^iGiven^iMiddle^iPrefix^iSuffix" +
				"",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "aFamily",
					GivenName:    "aGiven",
					MiddleName:   "aMiddle",
					NamePrefix:   "aPrefix",
					NameSuffix:   "aSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "iFamily",
					GivenName:    "iGiven",
					MiddleName:   "iMiddle",
					NamePrefix:   "iPrefix",
					NameSuffix:   "iSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepNone,
			},
			IsEmpty: false,
		},
		// No Ideographic
		{
			Raw: "aFamily^aGiven^aMiddle^aPrefix^aSuffix=" +
				"=" +
				"",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "aFamily",
					GivenName:    "aGiven",
					MiddleName:   "aMiddle",
					NamePrefix:   "aPrefix",
					NameSuffix:   "aSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepAll,
			},
			IsEmpty: false,
		},
		// No Ideographic, no seps
		{
			Raw: "aFamily^aGiven^aMiddle^aPrefix^aSuffix",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "aFamily",
					GivenName:    "aGiven",
					MiddleName:   "aMiddle",
					NamePrefix:   "aPrefix",
					NameSuffix:   "aSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepNone,
			},
			IsEmpty: false,
		},
		// No Alphabetic, with seps
		{
			Raw: "=" +
				"iFamily^iGiven^iMiddle^iPrefix^iSuffix=" +
				"pFamily^pGiven^pMiddle^pPrefix^pSuffix",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "iFamily",
					GivenName:    "iGiven",
					MiddleName:   "iMiddle",
					NamePrefix:   "iPrefix",
					NameSuffix:   "iSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "pFamily",
					GivenName:    "pGiven",
					MiddleName:   "pMiddle",
					NamePrefix:   "pPrefix",
					NameSuffix:   "pSuffix",
					NullSepLevel: GroupNullSepNone,
				},
			},
			IsEmpty: false,
		},
		// Alphabetic, with only 1 sep
		{
			Raw: "aFamily^aGiven^aMiddle^aPrefix^aSuffix=",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "aFamily",
					GivenName:    "aGiven",
					MiddleName:   "aMiddle",
					NamePrefix:   "aPrefix",
					NameSuffix:   "aSuffix",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepIdeographic,
			},
			IsEmpty: false,
		},
		// Empty with seps
		{
			Raw: "=" +
				"=" +
				"",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepAll,
			},
			IsEmpty: true,
		},
		// Empty no seps
		{
			Raw: "",
			Expected: Info{
				Alphabetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Ideographic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				Phonetic: GroupInfo{
					FamilyName:   "",
					GivenName:    "",
					MiddleName:   "",
					NamePrefix:   "",
					NameSuffix:   "",
					NullSepLevel: GroupNullSepNone,
				},
				NullSepLevel: InfoNullSepNone,
			},
			IsEmpty: true,
		},
	}

	for _, tc := range testCases {

		// Test creating a new Info value and getting it's DCM() value.
		t.Run(tc.Raw+"_String", func(t *testing.T) {
			newInfo := Info{
				Alphabetic:   tc.Expected.Alphabetic,
				Ideographic:  tc.Expected.Ideographic,
				Phonetic:     tc.Expected.Phonetic,
				NullSepLevel: tc.Expected.NullSepLevel,
			}

			if tc.Raw != newInfo.DCM() {
				t.Errorf(
					"formatted string: expected '%v', got '%v'",
					tc.Raw,
					newInfo.DCM(),
				)
			}
		})

		// Test parsing a full PN value.
		t.Run(tc.Raw+"_parse", func(t *testing.T) {

			parsed, err := Parse(tc.Raw)
			if err != nil {
				t.Fatal("error parsing value:", err)
			}

			checkGroupInfo(
				t,
				tc.Expected.Alphabetic,
				parsed.Alphabetic,
				tc.Expected.Alphabetic.DCM(),
				"Alphabetic",
			)
			checkGroupInfo(
				t,
				tc.Expected.Ideographic,
				parsed.Ideographic,
				tc.Expected.Ideographic.DCM(),
				"Ideographic",
			)
			checkGroupInfo(
				t,
				tc.Expected.Phonetic,
				parsed.Phonetic,
				tc.Expected.Phonetic.DCM(),
				"Phonetic",
			)
		})

		// Test the .IsEmpty() method.
		t.Run(tc.Raw+"_isEmpty", func(t *testing.T) {
			newInfo, err := Parse(tc.Raw)
			if err != nil {
				t.Error("error parsing value:", err)
				t.FailNow()
			}

			if tc.IsEmpty != newInfo.IsEmpty() {
				t.Errorf(
					".IsEmpty(): got '%v', expected '%v'",
					tc.IsEmpty,
					newInfo.IsEmpty(),
				)
			}
		})
	}
}

func TestParse_Err(t *testing.T) {
	testCases := []struct {
		Raw       string
		ErrString string
	}{
		{
			Raw: "===",
			ErrString: "error parsing PN value: no more than 3 groups with " +
				"'[Alphabetic]=[Ideographic]=[Phonetic]' format are allowed: value " +
				"contains 4 groups. see 'PN' entry in official dicom spec: " +
				"http://dicom.nema.org/medical/dicom/current/output/html/part05.html#sect_6.2",
		},
		{
			Raw: "^^^^^",
			ErrString: "error parsing PN value: no more than 5 segments with " +
				"'[Last]^[First]^[Middle]^[Prefix]^[Suffix]' format are allowed: " +
				"value group Alphabetic contains 6 segments. see 'PN' entry in " +
				"official dicom spec: " +
				"http://dicom.nema.org/medical/dicom/current/output/html/part05.html#sect_6.2",
		},
		{
			Raw: "=^^^^^",
			ErrString: "error parsing PN value: no more than 5 segments with " +
				"'[Last]^[First]^[Middle]^[Prefix]^[Suffix]' format are allowed: " +
				"value group Ideographic contains 6 segments. see 'PN' entry in " +
				"official dicom spec: " +
				"http://dicom.nema.org/medical/dicom/current/output/html/part05.html#sect_6.2",
		},
		{
			Raw: "==^^^^^",
			ErrString: "error parsing PN value: no more than 5 segments with " +
				"'[Last]^[First]^[Middle]^[Prefix]^[Suffix]' format are allowed: " +
				"value group Phonetic contains 6 segments. see 'PN' entry in " +
				"official dicom spec: " +
				"http://dicom.nema.org/medical/dicom/current/output/html/part05.html#sect_6.2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Raw, func(t *testing.T) {
			_, err := Parse(tc.Raw)

			// Test errors.Is() with ErrParsePersonName
			if !errors.Is(err, ErrParsePersonName) {
				t.Errorf("err is not ErrParsePersonName")
			}

			if err.Error() != tc.ErrString {
				t.Error("error string unexpected:", err.Error())
			}
		})
	}
}

func TestInfo_WithNullSeparators(t *testing.T) {
	pnVal := "Potter^Harry"

	parsed, err := Parse(pnVal)
	if err != nil {
		t.Fatalf("error parsing pnVal: %v", err)
	}

	altered := parsed.WithTrailingNulls()

	expected := "Potter^Harry^^^=^^^^=^^^^"
	if altered.DCM() != expected {
		t.Errorf("expected '%v', got '%v'", expected, altered.DCM())
	}
}

func TestInfo_WithoutNullSeparators(t *testing.T) {
	pnVal := "Potter^Harry^^^=^^^^=^^^^"

	parsed, err := Parse(pnVal)
	if err != nil {
		t.Fatalf("error parsing pnVal: %v", err)
	}

	altered := parsed.WithoutTrailingNulls()

	expected := "Potter^Harry"
	if altered.DCM() != expected {
		t.Errorf("expected '%v', got '%v'", expected, altered.DCM())
	}
}
