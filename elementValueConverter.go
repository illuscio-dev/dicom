package dicom

import (
	"errors"
	"fmt"
	"github.com/suyashkumar/dicom/pkg/personname"
	"github.com/suyashkumar/dicom/pkg/tag"
	"reflect"
)

// Sentinel error for conversions.
var ErrConvertValue = errors.New("underlying not expected type")

// Wraps ErrConvertValue with some extra contextual information.
func newErrConvertValue(expected interface{}, actual interface{}) error {
	// Using %w we can wrap the error for errors.Is()
	return fmt.Errorf(
		"%w: expected '%v', but got '%v'",
		ErrConvertValue,
		reflect.TypeOf(expected),
		reflect.TypeOf(actual),
	)
}

// Sentinel value returned when casting to a single value, but multiple values are
// found.
var ErrMultipleValuesFound = errors.New("expected single value")
func newErrMultipleValuesFound(valueCount int) error {
	return fmt.Errorf("%w, but found %v", ErrMultipleValuesFound, valueCount)
}

// Sentinel error returned when converting to a singular value, ignoreSpec is false,
// and the VM for the Tag of the element in the Dicom spec is not '1'.
var ErrSpecNotSingle = errors.New("value multiplicity is not '1'")
func newErrSpecNotSingle(vmRaw string) error {
	return fmt.Errorf("%w: found '%v'", ErrSpecNotSingle, vmRaw)
}

// Check a list of values to see if it is a single value.
func checkSingleValue(elementTag tag.Tag, valueCount int, ignoreSpec bool) error {
	// If our value count is not 1, immediately return an error.
	if valueCount != 1 {
		return newErrMultipleValuesFound(valueCount)
	}

	// If we are ignoring the spec OR this is a private tag, we are good to go. We have
	// a single value, so we will convert it. If this is a private tag, we have no way
	// of checking the spec, so we can ignore that the user wants us to check it.
	if ignoreSpec || tag.IsPrivate(elementTag.Group) {
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

type ElementValueConverter struct {
	element *Element
}

// ToStrings tries to coerce the value from dicom.Element.Value.GetValue() to []string
// and returns an error on failure.
func (converter ElementValueConverter) ToStrings() ([]string, error) {
	value := converter.element.Value.GetValue()
	strings, ok := value.([]string)
	if !ok {
		return nil, newErrConvertValue(strings, value)
	}

	return strings, nil
}

// As ToStrings(), but panics on error.
func (converter ElementValueConverter) MustToStrings() []string {
	strings, err := converter.ToStrings()
	if err != nil {
		panic(err)
	}
	return strings
}

// ToStrings tries to coerce the value from dicom.Element.Value.GetValue() to a single
// string value, and returns an error on failure.
//
// This method will fail if the underlying value is the wrong type, or does not contain
// a single value.
//
// If ignoreSpec is false, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
func (converter ElementValueConverter) ToString(ignoreSpec bool) (string, error) {
	strings, err := converter.ToStrings()
	if err != nil {
		return "", err
	}

	err = checkSingleValue(converter.element.Tag, len(strings), ignoreSpec)
	if err != nil {
		return "", err
	}

	return strings[0], nil
}

// MustToString is as ToString, but panics on error.
func (converter ElementValueConverter) MustToString(ignoreSpec bool) string {
	stringVal, err := converter.ToString(ignoreSpec)
	if err != nil {
		panic(err)
	}
	return stringVal
}

// ToPersonNames tries to coerce the value from dicom.Element.Value.GetValue() to
// []pn.PersonName and returns an error on failure.
func (converter ElementValueConverter) ToPersonNames() ([]personname.Info, error) {
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
// If ignoreSpec is false, the VR of the element tag will be looked up in the dicom
// spec, and the operation will fail if it is not '1', regardless of how many values
// are in this specific instance of the element.
func (converter ElementValueConverter) ToPersonName(
	ignoreSpec bool,
) (personname.Info, error) {
	names, err := converter.ToPersonNames()
	if err != nil {
		return personname.Info{}, err
	}

	err = checkSingleValue(converter.element.Tag, len(names), ignoreSpec)
	if err != nil {
		return personname.Info{}, err
	}

	return names[0], nil
}
