package filterable_list_item

import (
	"github.com/mieubrisse/bubble-bath/components"
)

// This interface can be reimplemented for more interesting usecases
type Component interface {
	components.Component

	IsHighlighted() bool
	SetHighlighted(isHighlighted bool)
	GetValue() string
}
