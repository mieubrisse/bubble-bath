package filterable_list_item

import (
	"github.com/mieubrisse/bubble-bath"
)

// This interface can be reimplemented for more interesting usecases
type Component interface {
	bubble_bath.Component

	IsHighlighted() bool
	SetHighlighted(isHighlighted bool)

	// GetValue gets the list item's value, which is the value returned when asking "which item is selected?"
	GetValue() string
}
