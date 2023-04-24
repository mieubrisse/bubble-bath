package text_input

import "github.com/mieubrisse/bubble-bath/components"

type Component interface {
	components.InteractiveComponent

	GetValue() string
	SetValue(value string)
}
