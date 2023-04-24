package flexbox

import bubble_bath "github.com/mieubrisse/bubble-bath"

type Component interface {
	bubble_bath.InteractiveComponent

	// SetFocusedChildren indicates which children will receive events
	// All focused children receive all events
	SetFocusedChildren(focusedChildrenIndexSet map[int]bool)
}
