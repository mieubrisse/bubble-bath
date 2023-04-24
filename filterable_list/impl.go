package filterable_list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/bubble-bath"
	"github.com/mieubrisse/bubble-bath/filterable_list_item"
	"strings"
)

// This is an implementation of a filterable, scrollable list
type implementation[T filterable_list_item.Component] struct {
	unfilteredItems []T

	// The indices of the filtered items within the unfiltered items list
	filteredItemsOriginalIndices []int

	// The index of the highlighted item within the *filtered list*
	highlightedItemIdx int

	isFocused bool
	width     int
	height    int
}

func New[T filterable_list_item.Component]() Component[T] {
	return &implementation[T]{
		unfilteredItems:              make([]T, 0),
		filteredItemsOriginalIndices: make([]int, 0),
		highlightedItemIdx:           0,
		width:                        0,
		height:                       0,
	}
}

func (impl implementation[T]) View() string {
	if len(impl.filteredItemsOriginalIndices) == 0 {
		return ""
	}

	// As aesthetic choices, when there are more item lines than display lines:
	// 1. We want the entire list to scroll around the cursor if it's in the center of the screen, rather than
	//    the user needing to scroll to top or bottom to get the list to move. This helps the user see more
	//    relevant information at once
	// 2. When the cursor is near the top or bottom of the list, scroll the cursor rather than the entire list
	//    so that we don't get blank space
	// The easiest way to accomplish this is to calculate the range of acceptable first-line indexes of the view,
	//   which will range from [0, num_items - num_display_lines], and when the user is in the middle of the list
	//   the view will have the cursor line in the center
	halfHeight := impl.height / 2

	// Ensure that, when near the bottom of the list, the cursor is no longer centered and scrolls to the bottom
	firstDisplayedLineIdxInclusive := bubble_bath.GetMinInt(
		impl.highlightedItemIdx-halfHeight,
		len(impl.filteredItemsOriginalIndices)-impl.height,
	)

	// Ensure that, when near the top of the list, the cursor is no longer centered and scrolls to the top
	firstDisplayedLineIdxInclusive = bubble_bath.GetMaxInt(
		firstDisplayedLineIdxInclusive,
		0,
	)

	lastDisplayedLineIdxExclusive := bubble_bath.GetMinInt(
		len(impl.filteredItemsOriginalIndices),
		firstDisplayedLineIdxInclusive+impl.height,
	)

	if firstDisplayedLineIdxInclusive == lastDisplayedLineIdxExclusive {
		return ""
	}

	displayedItems := impl.filteredItemsOriginalIndices[firstDisplayedLineIdxInclusive:lastDisplayedLineIdxExclusive]

	resultLines := []string{}
	for _, originalItemIdx := range displayedItems {
		item := impl.unfilteredItems[originalItemIdx]

		renderedLine := item.View()
		resultLines = append(resultLines, renderedLine)
	}

	result := strings.Join(resultLines, "\n")

	// TODO truncating long lines by printable char

	return lipgloss.NewStyle().
		Width(impl.width).
		Height(impl.height).
		MaxWidth(impl.width).
		MaxHeight(impl.height).
		Render(result)
}

func (impl *implementation[T]) Update(msg tea.Msg) tea.Cmd {
	// Do nothing on non-Keymsgs
	switch msg.(type) {
	case tea.KeyMsg:
		// Proceed to rest of function
	default:
		return nil
	}

	if !impl.isFocused {
		return nil
	}

	// TODO allow for KeyMap overrides here?
	castedMsg := msg.(tea.KeyMsg)
	switch castedMsg.String() {
	case "j":
		impl.Scroll(1)
	case "k":
		impl.Scroll(-1)
	case "J":
		impl.Scroll(impl.height)
	case "K":
		impl.Scroll(-impl.height)
	}
	return nil
}

func (impl *implementation[T]) UpdateFilter(newFilter func(idx int, item T) bool, shouldPreserveHighlight bool) {
	// This is a hack to indicate "the filtered list was empty, so there's no highlighted item original idx"
	oldHighlightedItemOriginalIdx := -1

	// If there are items being displayed, de-highlight the current item (if there are any)
	if len(impl.filteredItemsOriginalIndices) > 0 {
		oldHighlightedItemOriginalIdx = impl.filteredItemsOriginalIndices[impl.highlightedItemIdx]
		oldHighlightedItem := impl.unfilteredItems[oldHighlightedItemOriginalIdx]
		oldHighlightedItem.SetHighlighted(false)
	}

	// By default, assume that the highlighted item in the pre-filter list doesn't exist in the
	// post-filter list (but we'll fix this below if the assumption is false)
	newHighlightedItemIdx := 0

	newFilteredItemOriginalIndices := []int{}
	for idx, item := range impl.unfilteredItems {
		if newFilter(idx, item) {
			newFilteredItemOriginalIndices = append(newFilteredItemOriginalIndices, idx)

			if shouldPreserveHighlight && idx == oldHighlightedItemOriginalIdx {
				newHighlightedItemIdx = len(newFilteredItemOriginalIndices) - 1
			}
		}
	}

	impl.filteredItemsOriginalIndices = newFilteredItemOriginalIndices
	impl.highlightedItemIdx = newHighlightedItemIdx

	// Highlight the new item (if possible)
	if len(impl.filteredItemsOriginalIndices) > 0 {
		originalIdx := impl.filteredItemsOriginalIndices[impl.highlightedItemIdx]
		item := impl.unfilteredItems[originalIdx]
		item.SetHighlighted(true)
	}
}

func (impl *implementation[T]) SetItems(items []T) {
	filteredIndices := []int{}
	for idx := range items {
		filteredIndices = append(filteredIndices, idx)
	}

	impl.unfilteredItems = items
	impl.filteredItemsOriginalIndices = filteredIndices
	impl.highlightedItemIdx = 0

	if len(impl.filteredItemsOriginalIndices) > 0 {
		highlightedItemOriginalIdx := impl.filteredItemsOriginalIndices[impl.highlightedItemIdx]
		item := impl.unfilteredItems[highlightedItemOriginalIdx]
		item.SetHighlighted(true)
	}
}

func (impl *implementation[T]) Scroll(scrollOffset int) {
	newHighlightedItemIdx := impl.highlightedItemIdx + scrollOffset
	if newHighlightedItemIdx < 0 {
		newHighlightedItemIdx = 0
	}
	if newHighlightedItemIdx >= len(impl.filteredItemsOriginalIndices) {
		newHighlightedItemIdx = len(impl.filteredItemsOriginalIndices) - 1
	}

	if newHighlightedItemIdx == impl.highlightedItemIdx {
		return
	}

	// De-highlight the previous item
	oldHighlightOriginalIdx := impl.filteredItemsOriginalIndices[impl.highlightedItemIdx]
	oldItem := impl.unfilteredItems[oldHighlightOriginalIdx]
	oldItem.SetHighlighted(false)

	// Highlight the new item
	newHighlightOriginalIdx := impl.filteredItemsOriginalIndices[newHighlightedItemIdx]
	newItem := impl.unfilteredItems[newHighlightOriginalIdx]
	newItem.SetHighlighted(true)

	impl.highlightedItemIdx = newHighlightedItemIdx
}

func (impl implementation[T]) GetItems() []T {
	return impl.unfilteredItems
}

// Gets the indices (within the original list) of the items currently being displayed
func (impl implementation[T]) GetFilteredItemIndices() []int {
	return impl.filteredItemsOriginalIndices
}

// GetHighlightedItemIndex returns the index *within the filtered list* of the highlighted item
func (impl implementation[T]) GetHighlightedItemIndex() int {
	return impl.highlightedItemIdx
}

func (impl *implementation[T]) Resize(width int, height int) {
	impl.width = width
	impl.height = height

	for _, item := range impl.unfilteredItems {
		// TODO Allow items to wrap (but requires a whole viewing framework)
		item.Resize(width, 1)
	}
}

func (impl implementation[T]) GetHeight() int {
	return impl.height
}

func (impl implementation[T]) GetWidth() int {
	return impl.width
}

func (impl *implementation[T]) SetFocus(isFocused bool) tea.Cmd {
	impl.isFocused = isFocused
	return nil
}

func (impl implementation[T]) IsFocused() bool {
	return impl.isFocused
}
