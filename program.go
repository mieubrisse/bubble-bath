package bubble_bath

import tea "github.com/charmbracelet/bubbletea"

type bubbleBathModel struct {
	initCmd   tea.Cmd
	component InteractiveComponent
}

// NewBubbleBathProgram creates a new tea.Model for tea.NewProgram based off InteractiveComponent
func NewBubbleBathProgram(app InteractiveComponent, initCmd tea.Cmd) tea.Model {
	return bubbleBathModel{
		initCmd:   initCmd,
		component: app,
	}
}

func (b bubbleBathModel) Init() tea.Cmd {
	return b.initCmd
}

func (b bubbleBathModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, b.component.Update(msg)
}

func (b bubbleBathModel) View() string {
	return b.component.View()
}
