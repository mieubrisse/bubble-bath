package bubble_bath

import tea "github.com/charmbracelet/bubbletea"

type bubbleBathModel struct {
	component InteractiveComponent
}

func NewBubbleBathProgram(app InteractiveComponent) tea.Model {
	return bubbleBathModel{component: app}
}

func (b bubbleBathModel) Init() tea.Cmd {

}

func (b bubbleBathModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, b.component.Update(msg)
}

func (b bubbleBathModel) View() string {

}
