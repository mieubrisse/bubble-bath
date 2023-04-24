package filterable_list_item

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/bubble-bath/bubble_bath/text_block"
	"github.com/mieubrisse/bubble-bath/global_styles"
)

// implementation is a basic implementation of a list item
// it can be reimplemented as needed
type implementation struct {
	innerComponent text_block.Component

	contents string

	isHighlighted bool
	width         int
	height        int
}

func New(contents string) Component {
	inner := text_block.New(contents)
	return &implementation{
		innerComponent: inner,
		contents:       contents,
		isHighlighted:  false,
		width:          0,
		height:         0,
	}
}

func (impl implementation) View() string {
	lineStyle := lipgloss.NewStyle()
	if impl.isHighlighted {
		lineStyle = lineStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
	}

	// TODO do the cute little '...' cutoff
	return lineStyle.Render(impl.innerComponent.View())
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
	return impl.contents
}

func (impl implementation) IsHighlighted() bool {
	return impl.IsHighlighted()
}

func (impl *implementation) SetHighlighted(isHighlighted bool) {
	impl.isHighlighted = isHighlighted
}
