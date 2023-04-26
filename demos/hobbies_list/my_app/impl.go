package my_app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/bubble-bath/filterable_list"
	"github.com/mieubrisse/bubble-bath/filterable_list_item"
	"github.com/mieubrisse/bubble-bath/flexbox"
	"github.com/mieubrisse/bubble-bath/text_block"
)

type implementation struct {
	hobbiesAndTitle flexbox.Component

	width  int
	height int
}

func New() MyApp {
	hobbiesListTitleStyle := lipgloss.NewStyle().Bold(true)
	hobbiesListTitle := text_block.New("My hobbies:", text_block.WithStyle(hobbiesListTitleStyle))

	hobbies := []filterable_list_item.Component{
		filterable_list_item.New(text_block.New("Pourover coffee"), "coffee"),
		filterable_list_item.New(text_block.New("Coding"), "coding"),
		filterable_list_item.New(text_block.New("Jiu jitsu"), "bjj"),
	}
	hobbiesList := filterable_list.New[filterable_list_item.Component]()
	hobbiesList.SetItems(hobbies)
	hobbiesList.SetFocus(true)

	// Will flexibly resize as needed
	hobbiesAndTitle := flexbox.New(
		[]flexbox.FlexItem{
			{
				Component: hobbiesListTitle,
				FixedSize: 1,
			},
			{
				Component:  hobbiesList,
				FlexWeight: 1,
			},
			{
				Component:  hobbiesList,
				FlexWeight: 1,
			},
		},
		flexbox.WithDirection(flexbox.Vertical),
	)

	return &implementation{
		hobbiesAndTitle: hobbiesAndTitle,
		width:           0,
		height:          0,
	}
}

func (i implementation) Update(msg tea.Msg) tea.Cmd {
	return i.hobbiesAndTitle.Update(msg)
}

func (i implementation) View() string {
	return i.hobbiesAndTitle.View()
}

func (i *implementation) Resize(width int, height int) {
	i.width = width
	i.height = height
	i.hobbiesAndTitle.Resize(width, height)
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
