package flexbox

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bubble_bath "github.com/mieubrisse/bubble-bath"
	"math"
)

const (
	defaultShouldHandleChildrenFocus = true
)

// TODO Remove this; I'm not happy about these options (should probably just be setters)
type FlexboxOption func(*implementation)

func WithDirection(direction LayoutDirection) FlexboxOption {
	return func(impl *implementation) {
		impl.direction = direction
	}
}

func WithChildFocusManaging(shouldManageChildrenFocus bool) FlexboxOption {
	return func(impl *implementation) {
		impl.shouldManageChildrenFocus = shouldManageChildrenFocus
	}
}

type LayoutDirection int

const (
	Vertical LayoutDirection = iota
	Horizontal
)

type FlexItem struct {
	// Required
	Component bubble_bath.Component

	// The fixed size that the component should take up
	// 0 indicates no fixed size
	// Overrides FlexWeight
	FixedSize int

	// The weight that the item should have, when FixedSize is not set
	// 0 indicates that the item should get no weight (will be invisible)
	FlexWeight float64
}

type implementation struct {
	items []FlexItem

	direction LayoutDirection

	// "Set" of children that should receive events when the flexbox is focused
	focusReceivingChildrenIndexes map[int]bool

	// If true, the flexbox will focus and unfocus children when the flexbox itself is focused or unfocused
	shouldManageChildrenFocus bool

	isFocused bool
	width     int
	height    int
}

// New constructs a new flexbox Component
// As a convenience, if child focus management is enabled and any of the children are focused then:
// - those children will be set to receive focus from the flexbox
// - the flexbox's focus state will be set to true
func New(items []FlexItem, options ...FlexboxOption) Component {
	impl := &implementation{
		items:                         items,
		direction:                     Horizontal,
		focusReceivingChildrenIndexes: map[int]bool{},
		shouldManageChildrenFocus:     defaultShouldHandleChildrenFocus,
		isFocused:                     false,
		width:                         0,
		height:                        0,
	}
	for _, opt := range options {
		opt(impl)
	}

	if impl.shouldManageChildrenFocus {
		newFocusReceivingChildrenIndexes := map[int]bool{}
		for idx, item := range impl.items {
			switch component := item.Component.(type) {
			case bubble_bath.InteractiveComponent:
				if component.IsFocused() {
					newFocusReceivingChildrenIndexes[idx] = true
					impl.isFocused = true
				}
			}
		}
		impl.focusReceivingChildrenIndexes = newFocusReceivingChildrenIndexes
	}

	impl.alignChildFocusesIfNecessary()

	return impl
}

func (impl implementation) Update(msg tea.Msg) tea.Cmd {
	if !impl.isFocused {
		return nil
	}

	cmds := make([]tea.Cmd, 0)
	for idx, item := range impl.items {
		if _, found := impl.focusReceivingChildrenIndexes[idx]; !found {
			continue
		}

		switch component := item.Component.(type) {
		case bubble_bath.InteractiveComponent:
			cmds = append(cmds, component.Update(msg))
		}
	}
	return tea.Batch(cmds...)
}

func (impl implementation) View() string {
	// For coercing down the size of any unruly children who try to grow too big
	var bully func(string, int) string
	if impl.direction == Horizontal {
		bully = func(view string, itemWidth int) string {
			return lipgloss.NewStyle().
				MaxHeight(impl.height).
				MaxWidth(itemWidth).
				Render(view)
		}
	} else {
		bully = func(view string, itemHeight int) string {
			return lipgloss.NewStyle().
				MaxHeight(itemHeight).
				MaxWidth(impl.width).
				Render(view)
		}
	}

	itemViews := make([]string, len(impl.items))
	childSizes := impl.calculateChildSizes()
	for idx, size := range childSizes {
		itemComponent := impl.items[idx].Component
		itemViews[idx] = bully(itemComponent.View(), size)
	}

	var result string
	if impl.direction == Horizontal {
		result = lipgloss.JoinHorizontal(lipgloss.Center, itemViews...)
	} else {
		result = lipgloss.JoinVertical(lipgloss.Left, itemViews...)
	}

	// Add an extra sanity check to ensure we don't exceed our own bounds
	return lipgloss.NewStyle().
		MaxWidth(impl.width).
		MaxHeight(impl.height).
		Render(result)
}

func (impl *implementation) SetFocusReceivingChildren(focusedChildrenIndexSet map[int]bool) {
	impl.focusReceivingChildrenIndexes = focusedChildrenIndexSet
	impl.alignChildFocusesIfNecessary()
}

func (impl *implementation) Resize(width int, height int) {
	impl.width = width
	impl.height = height

	var resizingFunction func(component bubble_bath.Component, space int)
	if impl.direction == Horizontal {
		resizingFunction = func(component bubble_bath.Component, space int) {
			component.Resize(space, height)
		}
	} else {
		resizingFunction = func(component bubble_bath.Component, space int) {
			component.Resize(width, space)
		}
	}

	childSizes := impl.calculateChildSizes()
	for idx, size := range childSizes {
		childComponent := impl.items[idx].Component
		resizingFunction(childComponent, size)
	}
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl *implementation) SetFocus(isFocused bool) tea.Cmd {
	impl.isFocused = isFocused
	return impl.alignChildFocusesIfNecessary()
}

func (impl *implementation) IsFocused() bool {
	return impl.isFocused
}

// ====================================================================================================
//                                   Private Helper Functions
// ====================================================================================================

// Calculates per-child sizes along the major axis of the flexbox
func (impl *implementation) calculateChildSizes() []int {
	availableSpaceGetter := bubble_bath.Component.GetWidth
	if impl.direction == Vertical {
		availableSpaceGetter = bubble_bath.Component.GetHeight
	}

	availableSpace := availableSpaceGetter(impl)

	// First, add up the total sizes and fixed sizes
	totalFixedSizeConsumed := 0
	totalWeight := 0.0
	for _, item := range impl.items {
		if item.FixedSize != 0 {
			totalFixedSizeConsumed += item.FixedSize
		} else {
			totalWeight += item.FlexWeight
		}
	}

	spaceForFlexingElements := bubble_bath.GetMaxInt(0, availableSpace-totalFixedSizeConsumed)
	spacePerWeight := float64(spaceForFlexingElements) / totalWeight

	// Now, allocate
	results := make([]int, len(impl.items))
	for idx, item := range impl.items {
		var desiredItemSpace int
		if item.FixedSize != 0 {
			desiredItemSpace = item.FixedSize
		} else {
			desiredItemSpace = int(math.Round(item.FlexWeight * spacePerWeight))
		}
		actualItemSpace := bubble_bath.GetMinInt(availableSpace, desiredItemSpace)
		results[idx] = actualItemSpace

		availableSpace -= actualItemSpace
	}

	return results
}

// Idempotently aligns children to the right focus state
func (impl *implementation) alignChildFocusesIfNecessary() tea.Cmd {
	if !impl.shouldManageChildrenFocus {
		return nil
	}

	cmds := make([]tea.Cmd, 0)
	for idx, item := range impl.items {
		switch component := item.Component.(type) {
		case bubble_bath.InteractiveComponent:
			_, canChildReceiveFocus := impl.focusReceivingChildrenIndexes[idx]

			shouldChildBeFocused := canChildReceiveFocus && impl.isFocused

			// Skip sending the focus event for children that are already in the desired state
			if component.IsFocused() == shouldChildBeFocused {
				continue
			}

			cmds = append(cmds, component.SetFocus(shouldChildBeFocused))
		}
	}
	return tea.Batch(cmds...)
}
