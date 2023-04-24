package bubble_bath

import tea "github.com/charmbracelet/bubbletea"

type InteractiveComponent interface {
	Component

	// Update updates the model based on the given message
	// This is expected to be by-reference replacement, because doing by-value is just too messy
	// (you get into weird situations with generic interfaces, trying to return the model)
	Update(msg tea.Msg) tea.Cmd

	SetFocus(isFocused bool) tea.Cmd
	IsFocused() bool
}
