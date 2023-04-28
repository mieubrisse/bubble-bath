package my_app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bubble_bath "github.com/mieubrisse/bubble-bath"
	"github.com/mieubrisse/bubble-bath/resizable_text_block"
	"github.com/muesli/reflow/wordwrap"
	"math"
)

type implementation struct {
	text resizable_text_block.Component

	width  int
	height int
}

func New() MyApp {
	text := resizable_text_block.New("Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal.")

	return &implementation{
		text:   text,
		width:  0,
		height: 0,
	}
}

func (i implementation) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (i implementation) View() string {
	// return i.text.View()
	texts := []resizable_text_block.Component{
		resizable_text_block.New("Mother Mary lordy jesus"),
		resizable_text_block.New("Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal."),
		resizable_text_block.New("Foo bar bang this is a thing"),
	}

	sizes := make([]int, len(texts))
	sumMaximumIntrinsicWidths := 0
	for idx, text := range texts {
		maxIntrinsicWidth := text.GetMaximumIntrinsicWidth()
		sizes[idx] = maxIntrinsicWidth
		sumMaximumIntrinsicWidths += maxIntrinsicWidth
	}

	// TODO apply flex-grow
	// TODO make flex-shrink be optional
	// Apply flex-shrink
	availableSpace := i.width - sumMaximumIntrinsicWidths
	if availableSpace < 0 {
		// This is basically saying that flex-basis == maxContent
		totalShrinkApplied := 0
		for idx, size := range sizes {
			if idx == len(sizes)-1 {
				// Use all leftover shrink, no matter how the rounding has gone
				sizes[idx] = size + (availableSpace - totalShrinkApplied)
				break
			}

			// We use flex-basis as a multiplier for the shrink factor
			weight := float64(size) / float64(sumMaximumIntrinsicWidths)
			shrinkFloat := weight * float64(availableSpace)
			shrinkInt := int(math.Round(shrinkFloat))
			totalShrinkApplied += shrinkInt
			proposedReducedSize := size + shrinkInt

			minimumSize := texts[idx].GetMinimumIntrinsicWidth()
			sizes[idx] = bubble_bath.GetMaxInt(proposedReducedSize, minimumSize)
		}
	}

	results := make([]string, len(texts))
	for idx, size := range sizes {
		results[idx] = wordwrap.String(texts[idx].View(), size)
	}

	result := lipgloss.JoinHorizontal(lipgloss.Top, results...)

	// We're technically adding margin/padding
	return lipgloss.NewStyle().Height(i.height).MaxHeight(i.height).MaxWidth(i.width).Render(result)
}

func (i *implementation) Resize(width int, height int) {
	i.width = width
	i.height = height

	// i.hobbiesAndTitle.Resize(width, height)
}

func (i *implementation) GetWidth() int {
	return i.width
}

func (i *implementation) GetHeight() int {
	return i.height
}

func (i implementation) SetFocus(isFocused bool) tea.Cmd {
	// App is always focused
	return nil
}

func (i implementation) IsFocused() bool {
	return true
}
