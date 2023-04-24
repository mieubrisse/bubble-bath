package bubble_bath

type Component interface {
	View() string

	// Resize resizes the component
	// This is intentionally by-reference, because by-value is just too messy (see this repo's README)
	Resize(width int, height int)

	GetWidth() int
	GetHeight() int
}
