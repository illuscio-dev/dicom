package dcmtime

import (
	"errors"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	testCases := []struct {
		DAValue           string
		ExpectedString    string
		Expected          time.Time
		ExpectedPrecision PrecisionLevel
		AllowNema         bool
	}{
		// DICOM full date
		{
			DAValue:        "20200304",
			ExpectedString: "2020-03-04",
			Expected: time.Date(
				2020,
				3,
				4,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			ExpectedPrecision: PrecisionFull,
			AllowNema:         false,
		},
		// DICOM no day
		{
			DAValue:        "202003",
			ExpectedString: "2020-03",
			Expected: time.Date(
				2020,
				3,
				1,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			ExpectedPrecision: PrecisionMonth,
			AllowNema:         false,
		},
		// DICOM no month
		{
			DAValue:        "2020",
			ExpectedString: "2020",
			Expected: time.Date(
				2020,
				1,
				1,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			ExpectedPrecision: PrecisionYear,
			AllowNema:         false,
		},
		// NEMA full date
		{
			DAValue:        "2020.03.04",
			ExpectedString: "2020-03-04",
			Expected: time.Date(
				2020,
				3,
				4,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			ExpectedPrecision: PrecisionFull,
			AllowNema:         true,
		},
		// NEMA no day
		{
			DAValue:        "2020.03",
			ExpectedString: "2020-03",
			Expected: time.Date(
				2020,
				3,
				1,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			ExpectedPrecision: PrecisionMonth,
			AllowNema:         true,
		},
		// NEMA no month
		{
			DAValue:        "2020",
			ExpectedString: "2020",
			Expected: time.Date(
				2020,
				1,
				1,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			ExpectedPrecision: PrecisionYear,
			AllowNema:         true,
		},
	}

	for _, tc := range testCases {
		nameNema := ""
		if tc.AllowNema {
			nameNema = "_AllowNema"
		}

		var parsed Date
		t.Run(tc.DAValue+nameNema, func(t *testing.T) {
			var err error
			parsed, err = ParseDate(tc.DAValue, tc.AllowNema)
			if err != nil {
				t.Fatal("parse err:", err)
			}

			if !tc.Expected.Equal(parsed.Time) {
				t.Errorf(
					"parsed time (%v) != expected (%v)",
					parsed.Time,
					tc.Expected,
				)

			}

			if parsed.Precision != tc.ExpectedPrecision {
				t.Errorf(
					"precision: expected %v, got %v",
					tc.ExpectedPrecision.String(),
					parsed.Precision.String(),
				)
			}
		})

		t.Run(tc.DAValue+nameNema+"_String", func(t *testing.T) {
			stringVal := parsed.String()
			if stringVal != tc.ExpectedString {
				t.Fatalf(
					"got String() value '%v', expected '%v'",
					stringVal,
					tc.ExpectedString,
				)
			}
		})

		t.Run(tc.DAValue+nameNema+"_DCM", func(t *testing.T) {
			dcmVal := parsed.DCM()
			if dcmVal != tc.DAValue {
				t.Fatalf(
					"got DCM() value '%v', expected '%v'", dcmVal, tc.DAValue,
				)
			}
		})
	}
}

func TestParseDateErr(t *testing.T) {
	testCases := []struct {
		Name      string
		BadValue  string
		AllowNema bool
	}{
		{
			Name:      "TooManyDigits",
			BadValue:  "101002034",
			AllowNema: false,
		},
		{
			Name:      "TooManyDigits_AllowNema",
			BadValue:  "101002034",
			AllowNema: true,
		},
		{
			Name:      "MissingDigit_Days",
			BadValue:  "1010023",
			AllowNema: false,
		},
		{
			Name:      "MissingDigit_Days_AllowNema",
			BadValue:  "1010023",
			AllowNema: true,
		},
		{
			Name:      "MissingDigit_Months",
			BadValue:  "10102",
			AllowNema: false,
		},
		{
			Name:      "MissingDigit_Months_AllowNema",
			BadValue:  "10102",
			AllowNema: true,
		},
		{
			Name:      "MissingDigit_Year",
			BadValue:  "101",
			AllowNema: false,
		},
		{
			Name:      "MissingDigit_Year_AllowNema",
			BadValue:  "101",
			AllowNema: true,
		},
		{
			Name:      "NemaDate_AllowNemaFalse",
			BadValue:  "1010.02.03",
			AllowNema: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := ParseDate(tc.BadValue, tc.AllowNema)
			if !errors.Is(err, ErrParseDA) {
				t.Errorf("expected ErrParseDA, got %v", err)
			}
		})
	}
}

func TestDate_DCM(t *testing.T) {
	testCases := []struct {
		Time      time.Time
		Precision PrecisionLevel
		Expected  string
	}{
		{
			Time: time.Date(
				1010,
				2,
				3,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			Precision: PrecisionFull,
			Expected:  "10100203",
		},
		{
			Time: time.Date(
				1010,
				2,
				3,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			Precision: PrecisionDay,
			Expected:  "10100203",
		},
		{
			Time: time.Date(
				1010,
				2,
				3,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			Precision: PrecisionMonth,
			Expected:  "101002",
		},
		{
			Time: time.Date(
				1010,
				2,
				3,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			Precision: PrecisionYear,
			Expected:  "1010",
		},
	}

	for _, tc := range testCases {
		name := tc.Expected + "_" + tc.Precision.String()
		t.Run(name, func(t *testing.T) {
			da := Date{
				Time:      tc.Time,
				Precision: tc.Precision,
			}

			if da.DCM() != tc.Expected {
				t.Errorf("DCM(): expected '%v', got '%v'", tc.Expected, da.DCM())
			}
		})
	}
}
