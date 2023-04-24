package filterable_list_item

import (
	"github.com/charmbracelet/lipgloss"
)

var defaultHighlightedItemStyle = lipgloss.NewStyle().Background(lipgloss.Color("#282828")).Bold(true)

// implementation is a basic implementation of a list item
// More complex implementations can be created as needed
type implementation struct {
	// lipgloss.Style to apply to items that are highlighted
	HighlightedItemStyle lipgloss.Style

	innerComponent Component

	value string

	isHighlighted bool
	width         int
	height        int
}

func New(innerComponent Component, value string) Component {
	return &implementation{
		HighlightedItemStyle: defaultHighlightedItemStyle,
		innerComponent:       innerComponent,
		value:                value,
		isHighlighted:        false,
		width:                0,
		height:               0,
	}
}

func (impl implementation) View() string {
	result := impl.innerComponent.View()
	if impl.isHighlighted {
		result = impl.HighlightedItemStyle.Render(result)
	}
	return result
}

func (impl *implementation) Resize(width int, height int) {
	impl.innerComponent.Resize(width, height)
	impl.width = width
	impl.height = height
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetValue() string {
	return impl.value
}

func (impl implementation) IsHighlighted() bool {
	return impl.IsHighlighted()
}

func (impl *implementation) SetHighlighted(isHighlighted bool) {
	impl.isHighlighted = isHighlighted
}
