package flexbox

import bubble_bath "github.com/mieubrisse/bubble-bath"

// Component is a flexbox component which will automatically handle resizing and focus-event routing for multiple children
type Component interface {
	bubble_bath.InteractiveComponent

	// SetFocusReceivingChildren indicates which children should be focused when the flexbox is focused
	// All focused children receive all events
	// Children that are not bubble_bath.InteractiveComponent will of course not receive an event
	SetFocusReceivingChildren(focusReceivingChildrenIndexes map[int]bool)
}
