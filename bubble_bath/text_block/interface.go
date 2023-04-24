package text_block

import (
	"github.com/mieubrisse/bubble-bath/bubble_bath"
)

type Component interface {
	bubble_bath.InteractiveComponent

	// TODO remove this??
	GetContents() string
}
