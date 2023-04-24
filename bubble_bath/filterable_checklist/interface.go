package filterable_checklist

import (
	"github.com/mieubrisse/bubble-bath/bubble_bath"
	"github.com/mieubrisse/bubble-bath/bubble_bath/filterable_checklist_item"
	"github.com/mieubrisse/bubble-bath/bubble_bath/filterable_list"
)

type Component[T filterable_checklist_item.Component] interface {
	bubble_bath.InteractiveComponent

	// Used for manipulations of the inner list (no need to reimplement all the functions)
	// The items in the original list will match the items from GetItems
	GetFilterableList() filterable_list.Component[T]

	SetItems(items []T)
	GetItems() []T

	// GetSelectedItemOriginalIndices gets the indices within the current items list that are selected
	GetSelectedItemOriginalIndices() map[int]bool

	ToggleHighlightedItemSelection()

	// SetHighlightedItemSelection sets the selection for the currently-highlighted item
	SetHighlightedItemSelection(isSelected bool)

	// SetAllViewableItemsSelection sets the selection for all items that are currently shown (i.e. matching the filter)
	SetAllViewableItemsSelection(isSelected bool)

	// SetAllItemsSelection sets the selection on ALL items in the list (whether they're viewable or not)
	SetAllItemsSelection(isSelected bool)
}
