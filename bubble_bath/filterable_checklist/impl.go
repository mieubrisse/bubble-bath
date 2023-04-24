package filterable_checklist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/bubble-bath/bubble_bath/filterable_checklist_item"
	"github.com/mieubrisse/bubble-bath/bubble_bath/filterable_list"
)

type implementation[T filterable_checklist_item.Component] struct {
	innerList filterable_list.Component[T]

	items []T

	// Indices of selected items within the *unfiltered* list
	selectedItemIndices map[int]bool

	isFocused bool
	width     int
	height    int
}

func New[T filterable_checklist_item.Component]() Component[T] {
	inner := filterable_list.New[T]()
	return &implementation[T]{
		innerList:           inner,
		items:               make([]T, 0),
		selectedItemIndices: make(map[int]bool, 0),
		isFocused:           false,
		width:               0,
		height:              0,
	}
}

func (impl *implementation[T]) View() string {
	return impl.innerList.View()
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
	var returnCmd tea.Cmd
	castedMsg := msg.(tea.KeyMsg)
	switch castedMsg.String() {
	case "x", "enter":
		impl.ToggleHighlightedItemSelection()
	case "s":
		impl.SetAllViewableItemsSelection(true)
	case "d":
		impl.SetAllViewableItemsSelection(false)
	case "S":
		impl.SetAllItemsSelection(true)
	case "D":
		impl.SetAllItemsSelection(false)
	default:
		returnCmd = impl.innerList.Update(msg)
	}

	return returnCmd
}

func (impl implementation[T]) GetItems() []T {
	return impl.items
}

func (impl *implementation[T]) SetItems(items []T) {
	// TODO something about preserving the selected item indices when the list changes??
	impl.items = items
	impl.selectedItemIndices = make(map[int]bool, 0)

	impl.innerList.SetItems(items)
}

func (impl implementation[T]) GetFilterableList() filterable_list.Component[T] {
	return impl.innerList
}

func (impl implementation[T]) GetSelectedItemOriginalIndices() map[int]bool {
	return impl.selectedItemIndices
}

func (impl *implementation[T]) ToggleHighlightedItemSelection() {
	filteredItemOriginalIndicies := impl.innerList.GetFilteredItemIndices()
	if len(filteredItemOriginalIndicies) == 0 {
		return
	}
	itemOriginalIdx := filteredItemOriginalIndicies[impl.innerList.GetHighlightedItemIndex()]
	item := impl.items[itemOriginalIdx]
	isSelected := item.IsSelected()
	impl.setItemSelection(itemOriginalIdx, !isSelected)
}

func (impl *implementation[T]) SetHighlightedItemSelection(isSelected bool) {
	filteredItemIndices := impl.innerList.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return
	}

	highlightedItemIdxInFilteredList := impl.innerList.GetHighlightedItemIndex()
	highlightedItemIdxInOriginalList := filteredItemIndices[highlightedItemIdxInFilteredList]

	impl.setItemSelection(highlightedItemIdxInOriginalList, isSelected)
}

func (impl *implementation[T]) SetAllViewableItemsSelection(isSelected bool) {
	filteredItemIndices := impl.innerList.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return
	}

	for _, originalItemIdx := range filteredItemIndices {
		impl.setItemSelection(originalItemIdx, isSelected)
	}
}

func (impl *implementation[T]) SetAllItemsSelection(isSelected bool) {
	for idx := range impl.items {
		impl.setItemSelection(idx, isSelected)
	}
}

func (impl *implementation[T]) Resize(width int, height int) {
	impl.width = width
	impl.height = height

	impl.innerList.Resize(width, height)
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

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
// setItemSelection sets the selection of the given item, and does the appropriate bookkeeping
func (impl *implementation[T]) setItemSelection(itemIdx int, isSelected bool) {
	item := impl.items[itemIdx]
	item.SetSelection(isSelected)

	if isSelected {
		impl.selectedItemIndices[itemIdx] = true
	} else {
		delete(impl.selectedItemIndices, itemIdx)
	}
}
