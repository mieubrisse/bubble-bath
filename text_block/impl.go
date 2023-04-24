package text_block

import "github.com/charmbracelet/lipgloss"

type impl struct {
	contents string

	// TODO add matched char index

	width  int
	height int
}

func New(contents string) Component {
	return &impl{
		contents: contents,
		width:    0,
		height:   0,
	}
}

func (item impl) GetContents() string {
	return item.contents
}

func (item impl) View() string {
	// TODO add the nice '...' for when the item is cut off
	return lipgloss.NewStyle().
		MaxWidth(item.width).
		MaxHeight(item.height).
		Render(item.contents)
}

func (item *impl) Resize(width int, height int) {
	item.width = width
	item.height = height
}

func (item impl) GetWidth() int {
	return item.width
}

func (item impl) GetHeight() int {
	return item.height
}
