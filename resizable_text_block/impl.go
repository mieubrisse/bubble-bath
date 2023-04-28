package resizable_text_block

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/wordwrap"
	"strings"
)

type Option func(*implementation)

func WithStyle(style lipgloss.Style) Option {
	return func(impl *implementation) {
		impl.style = style
	}
}

type implementation struct {
	style lipgloss.Style

	contents string

	// Cached upon creation
	minimumIntrinsicWidth  int
	maximumIntrinsicHeight int
	maximumIntrinsicWidth  int
	minimumIntrinsicHeight int

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

	result.maximumIntrinsicWidth = lipgloss.Width(contents)
	result.minimumIntrinsicHeight = lipgloss.Height(contents)

	minimumIntrinsicWidth := 0
	for _, word := range strings.Fields(contents) {
		printableWidth := ansi.PrintableRuneWidth(word)
		if printableWidth > minimumIntrinsicWidth {
			minimumIntrinsicWidth = printableWidth
		}
	}
	result.minimumIntrinsicWidth = minimumIntrinsicWidth
	result.maximumIntrinsicHeight = lipgloss.Height(wordwrap.String(contents, result.minimumIntrinsicWidth))

	return result
}

func (item *implementation) GetContents() string {
	return item.contents
}

func (item *implementation) View() string {
	// TODO add the nice '...' for when the item is cut off
	return item.style.Render(item.contents)
	/*
		return item.style.
			MaxWidth(item.width).
			MaxHeight(item.height).
			Render(item.contents)

	*/
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

func (item *implementation) GetMinimumIntrinsicWidth() int {
	return item.minimumIntrinsicWidth
}

func (item *implementation) GetMaximumIntrinsicWidth() int {
	return item.maximumIntrinsicWidth
}

func (item *implementation) GetHeightGivenWidth(width int) int {
	if width <= item.minimumIntrinsicWidth {
		return item.maximumIntrinsicHeight
	}
	wrappedStr := wordwrap.String(item.contents, width)
	return lipgloss.Width(wrappedStr)
}
