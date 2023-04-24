package bubble_bath

type Component interface {
	View() string

	// This is expected to be by-ref replacing
	// All our components are by-value replacing because it gets way too messy with interface generics when being
	// done by-value
	Resize(width int, height int)

	GetWidth() int
	GetHeight() int
}
