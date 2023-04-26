package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	bubble_bath "github.com/mieubrisse/bubble-bath"
	"github.com/mieubrisse/bubble-bath/demos/long_text/my_app"
	"os"
)

func main() {
	if _, err := bubble_bath.RunBubbleBathProgram(
		my_app.New(),
		nil,
		[]tea.ProgramOption{
			tea.WithAltScreen(),
		},
	); err != nil {
		fmt.Printf("An error occurred running the program:\n%v", err)
		os.Exit(1)
	}
}
