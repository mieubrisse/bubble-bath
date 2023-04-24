package text_block

import "github.com/charmbracelet/lipgloss"

type Option func(*implementation)

func WithStyle(style lipgloss.Style) Option {
	return func(impl *implementation) {
		impl.style = style
	}
}

type implementation struct {
	style lipgloss.Style

	contents string

	// TODO add matched char index

	width  int
	height int
}

func New(contents string, options ...Option) Component {
	result := &implementation{
		style:    lipgloss.Style{},
		contents: contents,
		width:    0,
		height:   0,
	}
	for _, opt := range options {
		opt(result)
	}
	return result
}

func (item *implementation) GetContents() string {
	return item.contents
}

func (item *implementation) View() string {
	// TODO add the nice '...' for when the item is cut off
	return item.style.
		MaxWidth(item.width).
		MaxHeight(item.height).
		Render(item.contents)
}

func (item *implementation) Resize(width int, height int) {
	item.width = width
	item.height = height
}

func (item *implementation) GetWidth() int {
	return item.width
}

func (item *implementation) GetHeight() int {
	return item.height
}
