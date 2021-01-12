package dicom

import (
	"errors"
	"fmt"
	"github.com/suyashkumar/dicom/pkg/dcmtime"
	"github.com/suyashkumar/dicom/pkg/personname"
	"github.com/suyashkumar/dicom/pkg/tag"
	"github.com/suyashkumar/dicom/pkg/vrraw"
	"reflect"
)

// ErrDecodeValue is the base error all decoder methods wrap.
var ErrDecodeValue = errors.New("error decoding element value")

// ErrConvertTypeValue is the sentinel error for failed type conversions. Wraps
// ErrDecodeValue
var ErrConvertTypeValue = fmt.Errorf(
	"%w: underlying not expected type", ErrDecodeValue,
)

// newErrConvertValue wraps ErrConvertTypeValue with some extra contextual information.
func newErrConvertValue(expected interface{}, actual interface{}) error {
	// Using %w we can wrap the error for errors.Is()
	return fmt.Errorf(
		"%w: expected '%v', but got '%v'",
		ErrConvertTypeValue,
		reflect.TypeOf(expected),
		reflect.TypeOf(actual),
	)
}

// ErrMultipleValuesFound is the sentinel error returned when casting to a single value,
// but multiple values are found. Wraps ErrDecodeValue
var ErrMultipleValuesFound = errors.New("expected single value")

func newErrMultipleValuesFound(valueCount int) error {
	return fmt.Errorf("%w, but found %v", ErrMultipleValuesFound, valueCount)
}

// ErrDicomSpecViolation is the Sentinel error for decode failures resulting from the
// requested operation violating the DICOM spec, regardless of whether it is technically
// feasible. Wraps ErrDecodeValue.
var ErrDicomSpecViolation = fmt.Errorf(
	"%w: operation violates dicom spec",
	ErrDecodeValue,
)

// ErrDicomSpecNotSingle is the sentinel error returned when converting to a singular
// value, checkSpec is true, and the VM for the Tag of the element in the Dicom spec
// is not '1'.
//
// Wraps ErrDicomSpecViolation.
var ErrDicomSpecNotSingle = fmt.Errorf(
	"%w: value multiplicity is not '1'", ErrDicomSpecViolation,
)

// newErrSpecNotSingle wraps ErrDicomSpecNotSingle with some additional context.
func newErrSpecNotSingle(vmRaw string) error {
	return fmt.Errorf("%w: found '%v'", ErrDicomSpecNotSingle, vmRaw)
}

// ErrWrongValueRepresentation is the sentinel error returned when we are trying to
// decode a value from the wrong VR. For instance, trying to decode a non-PN value as a
// personname.Info.
//
// wraps ErrDicomSpecViolation.
var ErrWrongValueRepresentation = fmt.Errorf(
	"%w: cannot decode Value Representation", ErrDicomSpecViolation,
)

// newErrBadVR wraps ErrWrongValueRepresentation with some additional context.
func newErrBadVR(targetValue interface{}, expectedVR string, foundVR string) error {
	return fmt.Errorf(
		"%w to type '%v': expected VR of '%v', got '%v'",
		ErrWrongValueRepresentation,
		reflect.TypeOf(targetValue),
		expectedVR,
		foundVR,
	)
}

// ElementValueDecoder offers element value conversion methods for extracting the inner
// value of an element.
type ElementValueDecoder struct {
	element *Element
}

// Check a list of values to see if it is a single value.
func (converter ElementValueDecoder) checkSingleValue(
	valueCount int, checkSpec bool,
) error {
	elementTag := converter.element.Tag

	// If our value count is not 1, immediately return an error.
	if valueCount != 1 {
		return newErrMultipleValuesFound(valueCount)
	}

	// If we are ignoring the spec OR this is a private tag, we are good to go. We have
	// a single value, so we will convert it. If this is a private tag, we have no way
	// of checking the spec, so we can ignore that the user wants us to check it.
	if !checkSpec || tag.IsPrivate(elementTag.Group) {
		return nil
	}

	// If we are not ignoring the spec, look up the dicom Tag info and see if it has a
	// VM of 1.
	info, err := tag.Find(elementTag)
	// If none is found, we are good to go.
	if err != nil {
		return nil
	}

	// If the spec says this is not a single value, return an error.
	if !info.VMInfo.IsSingleValue() {
		return newErrSpecNotSingle(info.VM)
	}

	return nil
}

// checkVR checks whether the VR we expect is what we found, adn returns an error if so.
func (converter ElementValueDecoder) checkVR(expectedVR string) bool {
	return expectedVR == converter.element.RawValueRepresentation
}

// Must returns a value decoder with decode methods panic if they hit an error.
func (converter ElementValueDecoder) Must() ElementMustValueDecoder {
	return ElementMustValueDecoder{converter: converter}
}

// ToBytes tries to coerce the value from dicom.Element.Value.GetValue() to []byte
// and returns an error on failure.
func (converter ElementValueDecoder) ToBytes() ([]byte, error) {
	value := converter.element.Value.GetValue()
	bytes, ok := value.([]byte)
	if !ok {
		return nil, newErrConvertValue(bytes, value)
	}

	return bytes, nil
}

// ToStrings tries to coerce the value from dicom.Element.Value.GetValue() to []string
// and returns an error on failure.
func (converter ElementValueDecoder) ToStrings() ([]string, error) {
	value := converter.element.Value.GetValue()
	strings, ok := value.([]string)
	if !ok {
		return nil, newErrConvertValue(strings, value)
	}

	return strings, nil
}

// ToString tries to coerce the value from dicom.Element.Value.GetValue() to a single
// string value, and returns an error on failure.
//
// This method will fail if the underlying value is the wrong type, or does not contain
// a single value.
//
// If checkSpec is true, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
func (converter ElementValueDecoder) ToString(checkSpec bool) (string, error) {
	strings, err := converter.ToStrings()
	if err != nil {
		return "", err
	}

	err = converter.checkSingleValue(len(strings), checkSpec)
	if err != nil {
		return "", err
	}

	return strings[0], nil
}

// ToInts tries to coerce the value from dicom.Element.Value.GetValue() to []int
// and returns an error on failure.
func (converter ElementValueDecoder) ToInts() ([]int, error) {
	value := converter.element.Value.GetValue()
	ints, ok := value.([]int)
	if !ok {
		return nil, newErrConvertValue(ints, value)
	}

	return ints, nil
}

// ToInt tries to coerce the value from dicom.Element.Value.GetValue() to a single
// int value, and returns an error on failure.
//
// This method will fail if the underlying value is the wrong type, or does not contain
// a single value.
//
// If checkSpec is true, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
func (converter ElementValueDecoder) ToInt(checkSpec bool) (int, error) {
	ints, err := converter.ToInts()
	if err != nil {
		return 0, err
	}

	err = converter.checkSingleValue(len(ints), checkSpec)
	if err != nil {
		return 0, err
	}

	return ints[0], nil
}

// ToPersonNames tries to coerce the value from dicom.Element.Value.GetValue() to
// []pn.PersonName and returns an error on failure.
func (converter ElementValueDecoder) ToPersonNames() ([]personname.Info, error) {
	if !converter.checkVR(vrraw.PersonName) {
		return nil, newErrBadVR(
			personname.Info{},
			vrraw.PersonName,
			converter.element.RawValueRepresentation,
		)
	}

	strings, err := converter.ToStrings()
	if err != nil {
		return nil, err
	}

	personNames := make([]personname.Info, len(strings))
	for i, thisString := range strings {
		thisPn, err := personname.Parse(thisString)
		if err != nil {
			return personNames, fmt.Errorf(
				"error converting string value %v to PersonName: %w",
				i,
				err,
			)
		}

		personNames[i] = thisPn
	}

	return personNames, nil
}

// ToPersonName tries to coerce the value from dicom.Element.Value.GetValue() to a
// single personname.Info value, and returns an error on failure.
//
// This method will fail if the underlying value is the wrong type, or does not contain
// a single value.
//
// If checkSpec is true, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
func (converter ElementValueDecoder) ToPersonName(
	checkSpec bool,
) (personname.Info, error) {
	names, err := converter.ToPersonNames()
	if err != nil {
		return personname.Info{}, err
	}

	err = converter.checkSingleValue(len(names), checkSpec)
	if err != nil {
		return personname.Info{}, err
	}

	return names[0], nil
}

// ToDates tries to coerce the value from dicom.Element.Value.GetValue() to
// []dcmtime.Date and returns an error on failure.
//
// If allowNema is true, parsing of pre-DICOM NEMA-300 style dates will be allowed.
func (converter ElementValueDecoder) ToDates(allowNema bool) ([]dcmtime.Date, error) {
	if !converter.checkVR(vrraw.Date) {
		return nil, newErrBadVR(
			personname.Info{},
			vrraw.PersonName,
			converter.element.RawValueRepresentation,
		)
	}

	dateStrings, err := converter.ToStrings()
	if err != nil {
		return nil, fmt.Errorf(
			"error getting value as strings for date conversion: %w", err,
		)
	}

	dates := make([]dcmtime.Date, len(dateStrings))
	for i, thisString := range dateStrings {
		thisDate, err := dcmtime.ParseDate(thisString, allowNema)
		if err != nil {
			return nil, fmt.Errorf("error parsing string value %v: %w", i, err)
		}
		dates[i] = thisDate
	}

	return dates, nil
}

// ToDate tries to coerce the value from dicom.Element.Value.GetValue() to a single
// dcmtime.Date value, and returns an error on failure.
//
// This method will fail if the underlying value is the wrong type, or does not contain
// a single value.
//
// If checkSpec is true, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
//
// If allowNema is true, parsing of pre-DICOM NEMA-300 style dates will be allowed.
func (converter ElementValueDecoder) ToDate(
	checkSpec bool, allowNema bool,
) (dcmtime.Date, error) {
	dates, err := converter.ToDates(allowNema)
	if err != nil {
		return dcmtime.Date{}, err
	}

	err = converter.checkSingleValue(len(dates), checkSpec)
	if err != nil {
		return dcmtime.Date{}, err
	}

	return dates[0], nil
}

// ToTimes tries to coerce the value from dicom.Element.Value.GetValue() to
// []dcmtime.Time and returns an error on failure.
func (converter ElementValueDecoder) ToTimes() ([]dcmtime.Time, error) {
	if !converter.checkVR(vrraw.Time) {
		return nil, newErrBadVR(
			personname.Info{},
			vrraw.Time,
			converter.element.RawValueRepresentation,
		)
	}

	timeStrings, err := converter.ToStrings()
	if err != nil {
		return nil, fmt.Errorf(
			"error getting value as strings for date conversion: %w", err,
		)
	}

	dates := make([]dcmtime.Time, len(timeStrings))
	for i, thisString := range timeStrings {
		thisDate, err := dcmtime.ParseTime(thisString)
		if err != nil {
			return nil, fmt.Errorf("error parsing string value %v: %w", i, err)
		}
		dates[i] = thisDate
	}

	return dates, nil
}

// ToTime tries to coerce the value from dicom.Element.Value.GetValue() to a single
// dcmtime.Time value, and returns an error on failure.
//
// This method will fail if the underlying value is the wrong type, or does not contain
// a single value.
//
// If checkSpec is true, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
func (converter ElementValueDecoder) ToTime(checkSpec bool) (dcmtime.Time, error) {
	times, err := converter.ToTimes()
	if err != nil {
		return dcmtime.Time{}, err
	}

	err = converter.checkSingleValue(len(times), checkSpec)
	if err != nil {
		return dcmtime.Time{}, err
	}

	return times[0], nil
}

// ToDatetimes tries to coerce the value from dicom.Element.Value.GetValue() to
// []dcmtime.Datetime and returns an error on failure.
func (converter ElementValueDecoder) ToDatetimes() ([]dcmtime.Datetime, error) {
	if !converter.checkVR(vrraw.DateTime) {
		return nil, newErrBadVR(
			personname.Info{},
			vrraw.DateTime,
			converter.element.RawValueRepresentation,
		)
	}

	timeStrings, err := converter.ToStrings()
	if err != nil {
		return nil, fmt.Errorf(
			"error getting value as strings for date conversion: %w", err,
		)
	}

	datetimes := make([]dcmtime.Datetime, len(timeStrings))
	for i, thisString := range timeStrings {
		thisDatetime, err := dcmtime.ParseDatetime(thisString)
		if err != nil {
			return nil, fmt.Errorf("error parsing string value %v: %w", i, err)
		}
		datetimes[i] = thisDatetime
	}

	return datetimes, nil
}

// ToDatetime tries to coerce the value from dicom.Element.Value.GetValue() to a single
// dcmtime.Datetime value, and returns an error on failure.
//
// This method will fail if the underlying value is the wrong type, or does not contain
// a single value.
//
// If checkSpec is true, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
func (converter ElementValueDecoder) ToDatetime(checkSpec bool) (dcmtime.Datetime, error) {
	datetimes, err := converter.ToDatetimes()
	if err != nil {
		return dcmtime.Datetime{}, err
	}

	err = converter.checkSingleValue(len(datetimes), checkSpec)
	if err != nil {
		return dcmtime.Datetime{}, err
	}

	return datetimes[0], nil
}

// ElementMustValueDecoder acts as a namespace for ElementValueDecoder methods
// that panic instead of returning an error.
type ElementMustValueDecoder struct {
	converter ElementValueDecoder
}

// ToBytes is as ElementValueDecoder.ToBytes(), but panics on error.
func (must ElementMustValueDecoder) ToBytes() []byte {
	bin, err := must.converter.ToBytes()
	if err != nil {
		panic(err)
	}
	return bin
}

// ToStrings is as ElementValueDecoder.ToStrings(), but panics on error.
func (must ElementMustValueDecoder) ToStrings() []string {
	strings, err := must.converter.ToStrings()
	if err != nil {
		panic(err)
	}
	return strings
}

// ToString is as ElementValueDecoder.ToString(), but panics on error.
func (must ElementMustValueDecoder) ToString(checkSpec bool) string {
	stringVal, err := must.converter.ToString(checkSpec)
	if err != nil {
		panic(err)
	}
	return stringVal
}

// ToInts is as ElementValueDecoder.ToInts(), but panics on error.
func (must ElementMustValueDecoder) ToInts() []int {
	intsVal, err := must.converter.ToInts()
	if err != nil {
		panic(err)
	}
	return intsVal
}

// ToInt is as ElementValueDecoder.ToInt(), but panics on error.
func (must ElementMustValueDecoder) ToInt(checkSpec bool) int {
	intVal, err := must.converter.ToInt(checkSpec)
	if err != nil {
		panic(err)
	}
	return intVal
}

// ToPersonNames is as ElementValueDecoder.ToPersonNames(), but panics on error.
func (must ElementMustValueDecoder) ToPersonNames() []personname.Info {
	names, err := must.converter.ToPersonNames()
	if err != nil {
		panic(err)
	}
	return names
}

// ToPersonName is as ElementValueDecoder.ToPersonName(), but panics on error.
func (must ElementMustValueDecoder) ToPersonName(checkSpec bool) personname.Info {
	name, err := must.converter.ToPersonName(checkSpec)
	if err != nil {
		panic(err)
	}
	return name
}

// ToDates is as ElementValueDecoder.ToDates(), but panics on error.
func (must ElementMustValueDecoder) ToDates(allowNema bool) []dcmtime.Date {
	names, err := must.converter.ToDates(allowNema)
	if err != nil {
		panic(err)
	}
	return names
}

// ToDate is as ElementValueDecoder.ToDate(), but panics on error.
func (must ElementMustValueDecoder) ToDate(
	checkSpec bool, allowNema bool,
) dcmtime.Date {
	date, err := must.converter.ToDate(checkSpec, allowNema)
	if err != nil {
		panic(err)
	}
	return date
}

// ToTimes is as ElementValueDecoder.ToTimes(), but panics on error.
func (must ElementMustValueDecoder) ToTimes() []dcmtime.Time {
	times, err := must.converter.ToTimes()
	if err != nil {
		panic(err)
	}
	return times
}

// ToTime is as ElementValueDecoder.ToTime(), but panics on error.
func (must ElementMustValueDecoder) ToTime(checkSpec bool) dcmtime.Time {
	timeVal, err := must.converter.ToTime(checkSpec)
	if err != nil {
		panic(err)
	}
	return timeVal
}

// ToDatetimes is as ElementValueDecoder.ToDatetimes(), but panics on error.
func (must ElementMustValueDecoder) ToDatetimes() []dcmtime.Datetime {
	datetimes, err := must.converter.ToDatetimes()
	if err != nil {
		panic(err)
	}
	return datetimes
}

// ToDatetime is as ElementValueDecoder.ToDatetime(), but panics on error.
func (must ElementMustValueDecoder) ToDatetime(checkSpec bool) dcmtime.Datetime {
	datetimeVal, err := must.converter.ToDatetime(checkSpec)
	if err != nil {
		panic(err)
	}
	return datetimeVal
}
