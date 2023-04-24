package my_app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/bubble-bath/filterable_list"
	"github.com/mieubrisse/bubble-bath/filterable_list_item"
	"github.com/mieubrisse/bubble-bath/text_block"
)

type implementation struct {
	title text_block.Component

	hobbies filterable_list.Component[filterable_list_item.Component]

	width  int
	height int
}

func New() MyApp {
	titleStyle := lipgloss.NewStyle().Bold(true)
	title := text_block.New("My hobbies:", text_block.WithStyle(titleStyle))

	hobbies := []filterable_list_item.Component{
		filterable_list_item.New(text_block.New("Pourover coffee"), "coffee"),
		filterable_list_item.New(text_block.New("Coding"), "coding"),
		filterable_list_item.New(text_block.New("Jiu jitsu"), "bjj"),
	}
	hobbiesList := filterable_list.New[filterable_list_item.Component]()
	hobbiesList.SetItems(hobbies)
	hobbiesList.SetFocus(true)

	return &implementation{
		title:   title,
		hobbies: hobbiesList,
	}
}

func (i implementation) Update(msg tea.Msg) tea.Cmd {
	return i.hobbies.Update(msg)
}

func (i implementation) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		i.title.View(),
		i.hobbies.View(),
	)
}

func (i *implementation) Resize(width int, height int) {
	i.width = width
	i.height = height
	i.title.Resize(width, 1)
	i.hobbies.Resize(width, height-1)
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
