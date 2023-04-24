package filterable_list_item

import (
	"github.com/mieubrisse/bubble-bath/bubble_bath"
)

// This interface can be reimplemented for more interesting usecases
type Component interface {
	bubble_bath.InteractiveComponent

	IsHighlighted() bool
	SetHighlighted(isHighlighted bool)
	GetValue() string
}
