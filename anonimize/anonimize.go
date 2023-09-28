package anonimize

import (
	"errors"
	"fmt"
)

const ErrNotImplemented = errors.New("operation not implemented")

// Delete the tag form dataset.
// X - remove
func remove(dataset Dataset, tag Tag) {
	element := dataset.Get(tag)
	deleteElement(dataset, element)
}

// Keep the tag in dataset, no effect.
// K - keep (unchanged for non-sequence attributes, cleaned for sequences)
func keep(dataset Dataset, tag Tag) {
	return
}

// Clean the tag in dataset.
// C - clean, that is replace with values of similar meaning known not to contain identifying information and consistent with the VR
func clean(dataset Dataset, tag Tag) {
	fmt.Errorf(ErrNotImplemented)
	return
}
