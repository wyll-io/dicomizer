package anonymize

import (
	"fmt"

	"github.com/suyashkumar/dicom"
)

func Anonymize(dataset *dicom.Dataset) error {
	for _, t := range tags {
		el, err := dataset.FindElementByTag(t.tag)
		if err != nil {
			if err == dicom.ErrorElementNotFound {
				continue
			}

			fmt.Printf("error searching tag %s: %v", t.tag, err)
			continue
		}

		switch t.value {
		case EmptyAction, DeleteOrEmptyAction:
			fmt.Printf("emptying tag %s\n", t.tag.String())
			if err := EmptyElementValue(el); err != nil {
				return err
			}
		case ReplaceAction, EmptyOrReplaceAction, DeleteOrReplaceAction, DeleteOrEmptyOrReplaceAction:
			fmt.Printf("replacing tag %s\n", t.tag.String())
			if err := ReplaceElementValue(el); err != nil {
				return err
			}
		case DeleteAction:
			fmt.Printf("deleting tag %s\n", t.tag.String())
			if err := DeleteElementValue(&dataset.Elements, t.tag, false); err != nil {
				return err
			}
		case ReplaceUIDAction:
			fmt.Printf("replacing tag UID %s\n", t.tag.String())
			if err := replaceElementUID(el); err != nil {
				return err
			}
		case DeleteOrEmptyOrReplaceUIDAction:
			if el.RawValueRepresentation == "UI" {
				replaceElementUID(el)
				continue
			}

			EmptyElementValue(el)
		default:
			panic(fmt.Sprintf("unknown value: %d", t.value))
		}
	}

	return nil
}
