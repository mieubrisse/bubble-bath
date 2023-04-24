package text_input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/bubble-bath/global_styles"
	"github.com/mieubrisse/bubble-bath/helpers"
	"github.com/muesli/ansi"
)

type Model struct {
	input           textinput.Model
	foregroundColor lipgloss.Color

	isFocused bool
	width     int
	height    int
}

func New(promptText string) Model {
	input := textinput.New()

	input.Prompt = promptText
	return Model{
		input:           input,
		foregroundColor: "",
		isFocused:       false,
		width:           0,
		height:          0,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model *Model) Update(msg tea.Msg) tea.Cmd {
	if !model.isFocused {
		return nil
	}

	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)
	return cmd
}

func (model Model) View() string {
	baseStyle := lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		Foreground(model.foregroundColor)
	if model.input.Focused() {
		baseStyle = baseStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
	}

	// It'd be really nice if textinput.Model allowed us to control the style based on focus, but alas it does not
	return baseStyle.Render(model.input.View())
}

func (model *Model) SetValue(newValue string) {
	model.input.SetValue(newValue)
}

func (model Model) GetValue() string {
	return model.input.Value()
}

func (model *Model) Focus() tea.Cmd {
	model.isFocused = true
	return model.input.Focus()
}

func (model *Model) Blur() tea.Cmd {
	model.isFocused = false
	model.input.Blur()
	return nil
}

func (model Model) Focused() bool {
	return model.isFocused
}

func (model *Model) SetForegroundColor(color lipgloss.Color) {
	model.foregroundColor = color
}

func (model *Model) Resize(width int, height int) {
	model.width = width
	model.height = height

	promptPrintableLength := ansi.PrintableRuneWidth(model.input.Prompt)

	// I'm not actually sure why we need the extra - 1 here (something to do with how Charm renders the max width); if we
	// don't have it though, things get weird
	maxNumDesiredDisplayedChars := width - promptPrintableLength - 1
	maxNumActualDisplayedChars := helpers.GetMaxInt(0, maxNumDesiredDisplayedChars)

	// The width on the Charm input is actually the max number of characters displayed at once NOT including the prompt!
	// This is why we do all the calculations prior
	model.input.Width = maxNumActualDisplayedChars
}

func (model Model) GetHeight() int {
	return model.height
}

func (model Model) GetWidth() int {
	return model.width
}
