package anonymize

import (
	"fmt"
	"math/rand"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

const (
	// ROOT_UID is the root UID allocated for this project
	ROOT_UID = "1.2.826.0.1.3680043.10.1336"
	// DEFAULT_DT is the default value for VR "DT"
	DEFAULT_DT = "00010101010101.000000+0000"
	// DEFAULT_DA is the default value for VR "DA"
	DEFAULT_DA = "00010101"
	// DEFAULT_TM is the default value for VR "TM"
	DEFAULT_TM = "000000.00"
	// DEFAUlT_ANON is the default value for anonymized fields
	DEFAUlT_ANON = "Anonymized"
)

func replaceElementUID(el *dicom.Element) error {
	uid := fmt.Sprintf("%s.%d", ROOT_UID, rand.Uint64())
	if len(uid) > 64 {
		uid = uid[:64]
	}

	// TODO: add multivalue support if necessary
	v, err := dicom.NewValue([]string{uid})
	if err != nil {
		return err
	}

	el.Value = v

	return nil
}

// ReplaceElementValue replaces the value of the given element. If the element is a
// sequence, it will replace all the sub elements accordingly.
func ReplaceElementValue(el *dicom.Element) error {
	var v dicom.Value
	var err error

	switch el.RawValueRepresentation {
	case "LO", "LT", "SH", "PN", "CS", "ST", "UT":
		v, err = dicom.NewValue([]string{DEFAUlT_ANON})
		if err != nil {
			return err
		}
	case "DS", "IS":
		v, err = dicom.NewValue([]string{"0"})
		if err != nil {
			return err
		}
	case "UI":
		return replaceElementUID(el)
	case "FD", "FL", "SS", "US", "SL", "UL":
		v, err = dicom.NewValue([]int{0})
		if err != nil {
			return err
		}
	case "DT":
		v, err = dicom.NewValue([]string{DEFAULT_DT})
		if err != nil {
			return err
		}
	case "DA":
		v, err = dicom.NewValue([]string{DEFAULT_DA})
		if err != nil {
			return err
		}
	case "TM":
		v, err = dicom.NewValue([]string{DEFAULT_TM})
		if err != nil {
			return err
		}
	case "UN":
		v, err = dicom.NewValue([]byte(DEFAUlT_ANON))
		if err != nil {
			return err
		}
	case "SQ":
		if el.Value.ValueType() == dicom.Sequences {
			seqs := el.Value.GetValue().([]*dicom.SequenceItemValue)
			for _, sq := range seqs {
				subElements := sq.GetValue().([]*dicom.Element)
				for _, subEl := range subElements {
					if err := ReplaceElementValue(subEl); err != nil {
						return err
					}
				}
			}

			v, err = dicom.NewValue(seqs)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid value type for SQ: %d", el.Value.ValueType())
		}
	default:
		panic(fmt.Sprintf("unknown VR: %s", el.RawValueRepresentation))
	}

	el.Value = v

	return nil
}

// EmptyElementValue empties the value of the given element. If the element is a
// sequence, it will empty all the sub elements accordingly.
func EmptyElementValue(el *dicom.Element) error {
	var v dicom.Value
	var err error

	switch el.RawValueRepresentation {
	case "SH", "PN", "UI", "LO", "LT", "CS", "AS", "ST", "UT":
		v, err = dicom.NewValue([]string{""})
		if err != nil {
			return err
		}
	case "DT":
		v, err = dicom.NewValue([]string{DEFAULT_DT})
		if err != nil {
			return err
		}
	case "DA":
		v, err = dicom.NewValue([]string{DEFAULT_DA})
		if err != nil {
			return err
		}
	case "TM":
		v, err = dicom.NewValue([]string{DEFAULT_TM})
		if err != nil {
			return err
		}
	case "UL", "FL", "FD", "SL", "SS", "US":
		v, err = dicom.NewValue([]int{0})
		if err != nil {
			return err
		}
	case "DS", "IS":
		v, err = dicom.NewValue([]string{"0"})
		if err != nil {
			return err
		}
	case "UN":
		v, err = dicom.NewValue([]byte(""))
		if err != nil {
			return err
		}
	case "SQ":
		if el.Value.ValueType() == dicom.Sequences {
			seqs := el.Value.GetValue().([]*dicom.SequenceItemValue)
			for _, sq := range seqs {
				subElements := sq.GetValue().([]*dicom.Element)
				for _, subEl := range subElements {
					if err := EmptyElementValue(subEl); err != nil {
						return err
					}
				}
			}

			v, err = dicom.NewValue(seqs)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid value type for SQ: %d", el.Value.ValueType())
		}
	default:
		panic(fmt.Sprintf("unknown VR: %s", el.RawValueRepresentation))
	}

	el.Value = v

	return nil
}

// DeleteElementValue deletes a tag from a given dataset. If the element is a sequence,
// it will delete all the sub elements accordingly. In case the sequence is empty
// after deleting the tag, the sequence will be deleted as well.
func DeleteElementValue(elements *[]*dicom.Element, t tag.Tag, ignoreTagCheck bool) error {
	newElements := []*dicom.Element{}
	for _, el := range *elements {
		if !ignoreTagCheck && el.Tag != t {
			// * element should not be deleted
			newElements = append(newElements, el)
			continue
		}

		if el.RawValueRepresentation == "DA" {
			// * date should be replaced
			if err := ReplaceElementValue(el); err != nil {
				return err
			}
		} else if el.RawValueRepresentation == "SQ" && el.Value.ValueType() == dicom.Sequences {
			// * iterate over all sub elements and applying the same logic
			output := [][]*dicom.Element{}
			for _, sq := range el.Value.GetValue().([]*dicom.SequenceItemValue) {
				subElements := sq.GetValue().([]*dicom.Element)
				if err := DeleteElementValue(&subElements, t, true); err != nil {
					return err
				}

				if len(subElements) > 0 {
					output = append(output, subElements)
				}
			}

			if len(output) == 0 {
				continue
			}

			v, err := dicom.NewValue(output)
			if err != nil {
				return err
			}

			el.Value = v
		} else {
			continue
		}

		newElements = append(newElements, el)
	}

	*elements = newElements

	return nil
}
