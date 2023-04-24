package filterable_list

import (
	"github.com/mieubrisse/bubble-bath/bubble_bath"
	"github.com/mieubrisse/bubble-bath/bubble_bath/filterable_list_item"
)

type Component[T filterable_list_item.Component] interface {
	bubble_bath.InteractiveComponent

	// UpdateFilter updates the filter by which items are currently being shown (or not)
	// If shouldPreserveHighlight is set, the highlighted item in the pre-update list will be the highlighted item
	// in the post-update list (as long as it exists)
	UpdateFilter(newFilter func(idx int, item T) bool, shouldPreserveHighlight bool)

	// Scroll scrolls the highlighted selection up or down by the specified number of items, with safeguards to
	// prevent scrolling off the ends of the list
	Scroll(scrollOffset int)

	GetItems() []T
	SetItems(items []T)

	GetFilteredItemIndices() []int
	GetHighlightedItemIndex() int
}
