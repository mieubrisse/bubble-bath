package text_input

import (
	"github.com/mieubrisse/bubble-bath"
)

type Component interface {
	bubble_bath.InteractiveComponent

	GetValue() string
	SetValue(value string)
}
