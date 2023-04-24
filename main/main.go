package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/bubble-bath"
	"github.com/mieubrisse/bubble-bath/main/my_app"
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
