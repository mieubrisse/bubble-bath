package filterable_checklist_item

import (
	"github.com/mieubrisse/bubble-bath/bubble_bath/filterable_list_item"
)

type Component interface {
	filterable_list_item.Component

	IsSelected() bool
	SetSelection(isSelected bool)
}
