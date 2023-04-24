package bubble_bath

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BubbleBathOption func(*bubbleBathModel)

func WithInitCmd(cmd tea.Cmd) BubbleBathOption {
	return func(model *bubbleBathModel) {
		model.initCmd = cmd
	}
}

func WithQuitSequences(quitSequenceSet map[string]bool) BubbleBathOption {
	return func(model *bubbleBathModel) {
		model.quitSequenceSet = quitSequenceSet
	}
}

var defaultQuitSequenceSet = map[string]bool{
	"ctrl+c": true,
	"ctrl+d": true,
}

type bubbleBathModel struct {
	// The tea.Cmd that will be fired upon initialization
	initCmd tea.Cmd

	// Sequences matching String() of tea.KeyMsg that will quit the program
	quitSequenceSet map[string]bool

	appComponent InteractiveComponent
}

// NewBubbleBathModel creates a new tea.Model for tea.NewProgram based off the given InteractiveComponent
func NewBubbleBathModel(app InteractiveComponent, options ...BubbleBathOption) tea.Model {
	result := &bubbleBathModel{
		initCmd:         nil,
		quitSequenceSet: defaultQuitSequenceSet,
		appComponent:    app,
	}
	for _, opt := range options {
		opt(result)
	}
	return result
}

func (b bubbleBathModel) Init() tea.Cmd {
	return b.initCmd
}

func (b bubbleBathModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if _, found := b.quitSequenceSet[msg.String()]; found {
			return b, tea.Quit

		}
	case tea.WindowSizeMsg:
		b.appComponent.Resize(msg.Width, msg.Height)
		return b, nil
	}

	return b, b.appComponent.Update(msg)
}

func (b bubbleBathModel) View() string {
	return b.appComponent.View()
}

func (b bubbleBathModel) GetAppComponent() InteractiveComponent {
	return b.appComponent
}

func RunBubbleBathProgram[T InteractiveComponent](
	appComponent T,
	bubbleBathOptions []BubbleBathOption,
	teaOptions []tea.ProgramOption,
) (T, error) {
	model := NewBubbleBathModel(appComponent, bubbleBathOptions...)

	finalModel, err := tea.NewProgram(model, teaOptions...).Run()
	castedModel := finalModel.(bubbleBathModel)
	castedAppComponent := castedModel.appComponent.(T)
	return castedAppComponent, err
}
