package filterable_list

import (
	"github.com/mieubrisse/bubble-bath/bubble_bath"
	"github.com/mieubrisse/bubble-bath/bubble_bath/filterable_list_item"
)

type Component[T filterable_list_item.Component] interface {
	bubble_bath.InteractiveComponent

	UpdateFilter(newFilter func(idx int, item T) bool)
	SetItems(items []T)
	Scroll(scrollOffset int)
	GetItems() []T
	GetFilteredItemIndices() []int
	GetHighlightedItemIndex() int
}
